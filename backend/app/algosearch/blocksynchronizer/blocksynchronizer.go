package blocksynchronizer

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	app "github.com/kevguy/algosearch/backend/business/algod"
	"github.com/kevguy/algosearch/backend/business/couchdata/account"
	"github.com/kevguy/algosearch/backend/business/couchdata/application"
	"github.com/kevguy/algosearch/backend/business/couchdata/asset"
	"github.com/kevguy/algosearch/backend/business/couchdata/block"
	"github.com/kevguy/algosearch/backend/business/couchdata/transaction"
	"github.com/kevguy/algosearch/backend/foundation/couchdb"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"sync"
	"time"
)

// BlockSynchronizer provides the ability to retrieve block data
// on an interval.
type BlockSynchronizer struct {
	log       			*zap.SugaredLogger
	wg        			sync.WaitGroup
	timer     			*time.Timer
	shutdown  			chan struct{}
	algodClient 		*algod.Client
	blockStore 			*block.Store
	transactionStore	*transaction.Store
	accountStore		*account.Store
	assetStore			*asset.Store
	appStore			*application.Store
}

// New creates a BlockSynchronizer for retrieving block data and saving it to CouchDB.
func New(log *zap.SugaredLogger, interval time.Duration, algodClient *algod.Client, cfg couchdb.Config) (*BlockSynchronizer, error) {
	p := BlockSynchronizer{
		log:       log,
		timer:     time.NewTimer(interval),
		shutdown:  make(chan struct{}),
		algodClient: algodClient,
	}

	db, err := couchdb.Open(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "connect to couchdb database")
	}

	blockStore := block.NewStore(log, db)
	p.blockStore = &blockStore

	transactionStore := transaction.NewStore(log, db)
	p.transactionStore = &transactionStore

	accountStore := account.NewStore(log, db)
	p.accountStore = &accountStore

	assetStore := asset.NewStore(log, db)
	p.assetStore = &assetStore

	appStore := application.NewStore(log, db)
	p.appStore = &appStore

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for {
			p.timer.Reset(interval)
			select {
			case <-p.timer.C:
				p.update()
			case <-p.shutdown:
				return
			}
		}
	}()

	return &p, nil
}

// Stop is used to shutdown the goroutine for syncing block data.
func (p *BlockSynchronizer) Stop() {
	close(p.shutdown)
	p.wg.Wait()
}

// TODO: add retry
// update pulls the block data and saves it to CouchDB.
func (p *BlockSynchronizer) update() {

	lastSyncedBlockNum, err := p.blockStore.GetLastSyncedRoundNumber(context.Background())
	if err != nil {
		p.log.Errorw("blocksynchronizer", "status", "get last synced round number", "ERROR", err)
	}

	currentRoundNum, err := app.GetCurrentRoundNum(context.Background(), p.algodClient)
	if err != nil {
		p.log.Errorw("blocksynchronizer", "status", "get current round num", "ERROR", err)
	}

	p.log.Infow("updating latest round here", "last synced round", lastSyncedBlockNum)

	if (currentRoundNum - lastSyncedBlockNum) > 1 {
		fmt.Printf("Trying to get round number: %d\n", lastSyncedBlockNum + 1)

		getRoundSuccessful := false
		var rawBlock []byte
		for !getRoundSuccessful {
			rawBlock, err = app.GetRoundInRawBytes(context.Background(), p.algodClient, lastSyncedBlockNum + 1)
			if err != nil {
				p.log.Errorw("blocksynchronizer", "status", "get round in raw bytes", "ERROR", err)
				// Assuming it's not just block data not available, jump to the next round
				lastSyncedBlockNum += 1
			} else {
				getRoundSuccessful = true
			}
		}

		//fmt.Printf("raw block: %v\n", rawBlock)
		//fmt.Printf("last synced num: %d\n", lastSyncedBlockNum + 1)
		p.log.Infof("Adding Round #%d\n", lastSyncedBlockNum + 1)

		newBlock, err := app.ConvertBlockRawBytes(context.Background(), rawBlock)
		if err != nil {
			p.log.Errorw("blocksynchronizer", "status", "convert raw bytes to block data", "ERROR", err)
		}

		//docID, rev, err := field.BlockStore.AddBlock(ctx, newBlock)
		blockDocId, blockDocRev, err := p.blockStore.AddBlock(context.Background(), newBlock)
		if err != nil {
			p.log.Errorw("blocksynchronizer", "status", "can't add new block", "ERROR", err)
		}
		p.log.Infof("\t- Added block %s with rev %s to CouchDB Block table\n", blockDocId, blockDocRev)

		var accountList []models.Account
		var assetList []models.Asset
		var appList []models.Application

		if len(newBlock.Transactions) > 0 {
			_, err = p.transactionStore.AddTransactions(context.Background(), newBlock.Transactions)
			if err != nil {
				p.log.Errorw("blocksynchronizer", "status", "can't add new transaction(s)", "ERROR", err)
			}
			p.log.Infof("\t\t- Added %d transactions with block %s to CouchDB Transaction table\n", len(newBlock.Transactions), newBlock.BlockHash)

			for _, txn := range newBlock.Transactions {

				accountIDs := app.ExtractAccountAddrsFromTxn(txn)
				applicationIDs := app.ExtractApplicationIdsFromTxn(txn)
				assetIDs := app.ExtractAssetIdsFromTxn(txn)

				for _, acctID := range accountIDs {
					accountInfo, err := app.GetAccount(context.Background(),"", p.log, p.algodClient, acctID)
					if err != nil {
						p.log.Errorw("blocksynchronizer", "status", "can't get account", "ERROR", err)
					}
					accountList = append(accountList, *accountInfo)
					//p.accountStore.AddAccount(context.Background(), *accountInfo)
				}

				for _, appID := range applicationIDs {
					appInfo, err := app.GetApplication(context.Background(),"", p.log, p.algodClient, appID)
					if err != nil {
						p.log.Errorw("blocksynchronizer", "status", "can't get app", "ERROR", err)
					}
					appList = append(appList, *appInfo)
					//p.appStore.AddApplication(context.Background(), *appInfo)
				}

				for _, assetID := range assetIDs {
					assetInfo, err := app.GetAsset(context.Background(),"", p.log, p.algodClient, assetID)
					if err != nil {
						p.log.Errorw("blocksynchronizer", "status", "can't get asset", "ERROR", err)
					}
					assetList = append(assetList, *assetInfo)
					//p.assetStore.AddAsset(context.Background(), *assetInfo)
				}
			}

		}

		_, err = p.accountStore.AddAccounts(context.Background(), accountList)
		if err != nil {
			p.log.Errorw("blocksynchronizer", "status", "can't add/update account(s)", "ERROR", err)
		}

		_, err = p.assetStore.AddAssets(context.Background(), assetList)
		if err != nil {
			p.log.Errorw("blocksynchronizer", "status", "can't add/update asset(s)", "ERROR", err)
		}

		_, err = p.appStore.AddApplications(context.Background(), appList)
		if err != nil {
			p.log.Errorw("blocksynchronizer", "status", "can't add/update application(s)", "ERROR", err)
		}

		//for _, transaction := range newBlock.Transactions {
		//	fmt.Println("Got transaction")
		//	fmt.Printf("%v\n", transaction)
		//	transactionDocId, transactionDocRev, err := p.transactionStore.AddTransaction(context.Background(), transaction)
		//	if err != nil {
		//		p.log.Errorw("blocksynchronizer", "status", "can't add new transaction", "ERROR", err)
		//	}
		//	p.log.Infof("\t\t- Added transaction %s with rev %s to CouchDB Transaction table\n", transactionDocId, transactionDocRev)
		//}
	}
}
