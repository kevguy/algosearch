package blocksynchronizer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/kevguy/algosearch/backend/business/core/account"
	algod2 "github.com/kevguy/algosearch/backend/business/core/algod"
	"github.com/kevguy/algosearch/backend/business/core/application"
	"github.com/kevguy/algosearch/backend/business/core/asset"
	"github.com/kevguy/algosearch/backend/business/core/block"
	"github.com/kevguy/algosearch/backend/business/core/transaction"
	"github.com/kevguy/algosearch/backend/foundation/couchdb"
	"github.com/kevguy/algosearch/backend/foundation/websocket"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"log"
	"sync"
	"time"
)

// BlockSynchronizer provides the ability to retrieve block data
// on an interval.
type BlockSynchronizer struct {
	log             *zap.SugaredLogger
	wg              sync.WaitGroup
	timer           *time.Timer
	shutdown        chan struct{}
	algodClient     *algod.Client
	blockCore       *block.Core
	transactionCore *transaction.Core
	accountCore     *account.Core
	assetCore       *asset.Core
	appCore         *application.Core
	algodCore 		*algod2.Core
	hub				*websocket.Hub
	dbName			string
}

// New creates a BlockSynchronizer for retrieving block data and saving it to CouchDB.
func New(log *zap.SugaredLogger, interval time.Duration, algodClient *algod.Client, cfg couchdb.Config, hub *websocket.Hub, dbName string) (*BlockSynchronizer, error) {
	p := BlockSynchronizer{
		log:       log,
		timer:     time.NewTimer(interval),
		shutdown:  make(chan struct{}),
		algodClient: algodClient,
		hub: 		hub,
		dbName: dbName,
	}

	db, err := couchdb.Open(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "connect to couchdb database")
	}

	algodCore := algod2.NewCore(log, algodClient)
	p.algodCore = &algodCore

	blockStore := block.NewCore(log, db, dbName)
	p.blockCore = &blockStore

	transactionStore := transaction.NewCore(log, db, dbName)
	p.transactionCore = &transactionStore

	accountStore := account.NewCore(log, db, dbName)
	p.accountCore = &accountStore

	assetStore := asset.NewCore(log, db, dbName)
	p.assetCore = &assetStore

	appStore := application.NewCore(log, db, dbName)
	p.appCore = &appStore

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

func (p *BlockSynchronizer) broadcastUpdate(newInfo WsMessage) error {
	fmt.Println(newInfo)
	res, err := json.Marshal(newInfo)
	if err != nil {
		return fmt.Errorf("marshaling new block info for broadcasting: %w", err)
	}
	p.hub.ExternalBroadcast <- res
	return nil
}

// TODO: add retry
// update pulls the block data and saves it to CouchDB.
func (p *BlockSynchronizer) update() {

	lastSyncedBlockNum, err := p.blockCore.GetLastSyncedRoundNumber(context.Background())
	if err != nil {
		p.log.Errorw("blocksynchronizer", "status", "get last synced round number", "ERROR", err)
	}
	// 111983 is amazing, it has a a ton of transactions
	//lastSyncedBlockNum = 111980

	currentRoundNum, err := p.algodCore.GetCurrentRoundNum(context.Background())
	if err != nil {
		p.log.Errorw("blocksynchronizer", "status", "get current round num", "ERROR", err)
	}
	p.log.Infof("Current Round Number is %d", currentRoundNum)

	p.log.Infow("Updating latest round here", "last synced round", lastSyncedBlockNum)

	if (currentRoundNum - lastSyncedBlockNum) > 1 {
		p.log.Infof("Trying to get round number: %d\n", lastSyncedBlockNum + 1)

		var newBlockPayload = WsMessage{
			Block:           models.Block{},
			TransactionList: nil,
			AccountList:     nil,
			AssetList:       nil,
			AppList:         nil,
		}

		getRoundSuccessful := false
		var rawBlock []byte
		for !getRoundSuccessful {
			rawBlock, err = p.algodCore.GetRoundInRawBytes(context.Background(), lastSyncedBlockNum + 1)
			if err != nil {
				p.log.Errorw("blocksynchronizer", "status", "get round in raw bytes", "ERROR", err)
				// Assuming it's not just block data not available, jump to the next round
				lastSyncedBlockNum += 1
			} else {
				getRoundSuccessful = true
				p.log.Infof("Block data for round #%d retrieved.", lastSyncedBlockNum + 1)
				//app.PrintBlockInfoFromRawBytes(rawBlock)
			}
		}

		//fmt.Printf("raw block: %v\n", rawBlock)
		//fmt.Printf("last synced num: %d\n", lastSyncedBlockNum + 1)
		p.log.Infof("Adding Round #%d\n", lastSyncedBlockNum + 1)

		newBlock, err := algod2.ConvertBlockRawBytes(context.Background(), rawBlock)
		if err != nil {
			p.log.Errorw("blocksynchronizer", "status", "convert raw bytes to block data", "ERROR", err)
		}
		//p.log.Infof("Block data for round %d: %v\n", newBlock.Round, newBlock)
		p.log.Infof("Block data for round %d after processing:", newBlock.Round)
		payloadStr, err := json.Marshal(newBlock.Block)
		if err != nil {
			p.log.Errorw("blocksynchronizer", "status", "marshal new block bytes to block data", "ERROR", err)
		}

		newBlockPayload.Block = newBlock.Block

		var prettyJSON bytes.Buffer
		err = json.Indent(&prettyJSON, payloadStr, "", "\t")
		if err != nil {
			log.Println("JSON parse error: ", err)
		}
		fmt.Println("Pretty pretty print json!!:", prettyJSON.String())
		//indexer.PrintBlockInfoFromJsonBlock(newBlock.Block)

		//docID, rev, err := field.BlockStore.AddBlock(ctx, newBlock)
		blockDocID, blockDocRev, err := p.blockCore.AddBlock(context.Background(), newBlock)
		if err != nil {
			p.log.Errorw("blocksynchronizer", "status", "can't add new block", "ERROR", err)
		}
		p.log.Infof("Added block %s with rev %s to CouchDB Block table", blockDocID, blockDocRev)

		var accountList []models.Account
		var assetList []models.Asset
		var appList []models.Application
		var txnIDList []string

		if len(newBlock.Transactions) > 0 {

			for _, txn := range newBlock.Transactions {
				txnIDList = append(txnIDList, txn.Id)
			}
			newBlockPayload.TransactionList = txnIDList

			_, err = p.transactionCore.AddTransactions(context.Background(), newBlock.Transactions)
			if err != nil {
				p.log.Errorw("blocksynchronizer", "status", "can't add new transaction(s)", "ERROR", err)
			}
			p.log.Infof("Added %d transactions with block %s to CouchDB Transaction table", len(newBlock.Transactions), newBlock.BlockHash)

			for _, txn := range newBlock.Transactions {

				accountIDs := algod2.ExtractAccountAddrsFromTxn(txn)
				applicationIDs := algod2.ExtractApplicationIdsFromTxn(txn)
				assetIDs := algod2.ExtractAssetIdsFromTxn(txn)

				for _, acctID := range accountIDs {
					accountInfo, err := p.algodCore.GetAccount(context.Background(),"", acctID)
					if err != nil {
						p.log.Errorw("blocksynchronizer", "status", "can't get account", "ERROR", err)
					}
					p.log.Infof("Retrieved Account info: %v", *accountInfo)
					accountList = append(accountList, *accountInfo)
				}

				for _, appID := range applicationIDs {
					appInfo, err := p.algodCore.GetApplication(context.Background(),"", appID)
					if err != nil {
						p.log.Errorw("blocksynchronizer", "status", "can't get app", "ERROR", err)
					}
					p.log.Infof("Retrieved Application info: %v", *appInfo)
					appList = append(appList, *appInfo)
				}

				for _, assetID := range assetIDs {
					assetInfo, err := p.algodCore.GetAsset(context.Background(),"", assetID)
					if err != nil {
						p.log.Errorw("blocksynchronizer", "status", "can't get asset", "ERROR", err)
					}
					p.log.Infof("Retrieved Asset info: %v", *assetInfo)
					assetList = append(assetList, *assetInfo)
				}
			}
		}

		if len(accountList) > 0 {
			_, err = p.accountCore.AddAccounts(context.Background(), accountList)
			if err != nil {
				p.log.Errorw("blocksynchronizer", "status", "can't add/update account(s)", "ERROR", err)
			}

			var acctIDList []string
			for _, acct := range accountList {
				acctIDList = append(acctIDList, acct.Address)
			}
			newBlockPayload.AccountList = acctIDList
		}

		if len(assetList) > 0 {
			_, err = p.assetCore.AddAssets(context.Background(), assetList)
			if err != nil {
				p.log.Errorw("blocksynchronizer", "status", "can't add/update asset(s)", "ERROR", err)
			}

			var assetIDList []uint64
			for _, asset := range assetList {
				assetIDList = append(assetIDList, asset.Index)
			}
			newBlockPayload.AssetList = assetIDList
		}

		if len(appList) > 0 {
			_, err = p.appCore.AddApplications(context.Background(), appList)
			if err != nil {
				p.log.Errorw("blocksynchronizer", "status", "can't add/update application(s)", "ERROR", err)
			}

			var appIDList []uint64
			for _, app := range appList {
				appIDList = append(appIDList, app.Id)
			}
			newBlockPayload.AppList = appIDList
		}

		speed, err := p.blockCore.GetBlockTxnSpeed(context.Background())
		if err != nil {
			p.log.Errorw("blocksynchronizer", "status", "can't find block txn speed", "ERROR", err)
		}
		newBlockPayload.AvgBlockTxnSpeed = speed

		if err = p.broadcastUpdate(newBlockPayload); err != nil {
			p.log.Errorw("blocksynchronizer", "status", "can't broadcast block update through websocket", "ERROR", err)
		}

		//for _, transaction := range newBlock.Transactions {
		//	fmt.Println("Got transaction")
		//	fmt.Printf("%v\n", transaction)
		//	transactionDocId, transactionDocRev, err := p.transactionCore.AddTransaction(context.Background(), transaction)
		//	if err != nil {
		//		p.log.Errorw("blocksynchronizer", "status", "can't add new transaction", "ERROR", err)
		//	}
		//	p.log.Infof("\t\t- Added transaction %s with rev %s to CouchDB Transaction table\n", transactionDocId, transactionDocRev)
		//}
	}
}


func GetAndInsertBlockData(
	log					*zap.SugaredLogger,
	algodClient			*algod.Client,
	blockCore *block.Core,
	transactionCore *transaction.Core,
	accountCore *account.Core,
	assetCore *asset.Core,
	appCore *application.Core,
	algodCore *algod2.Core,
	blockNum			uint64) error {
	log.Infof("Trying to get round number: %d\n", blockNum)

	getRoundSuccessful := false
	var rawBlock []byte
	var err error
	for !getRoundSuccessful {
		rawBlock, err = algodCore.GetRoundInRawBytes(context.Background(), blockNum)
		if err != nil {
			log.Errorw("blocksynchronizer", "status", "get round in raw bytes", "ERROR", err)
			// Assuming it's not just block data not available, jump to the next round
			blockNum += 1
		} else {
			getRoundSuccessful = true
			log.Infof("Block data for round #%d retrieved.", blockNum)
			//app.PrintBlockInfoFromRawBytes(rawBlock)
		}
	}

	//fmt.Printf("raw block: %v\n", rawBlock)
	//fmt.Printf("last synced num: %d\n", lastSyncedBlockNum + 1)
	log.Infof("Adding Round #%d\n", blockNum)

	newBlock, err := algod2.ConvertBlockRawBytes(context.Background(), rawBlock)
	if err != nil {
		log.Errorw("blocksynchronizer", "status", "convert raw bytes to block data", "ERROR", err)
		return err
	}
	//p.log.Infof("Block data for round %d: %v\n", newBlock.Round, newBlock)
	log.Infof("Block data for round %d after processing:", newBlock.Round)
	payloadStr, err := json.Marshal(newBlock.Block)
	if err != nil {
		log.Errorw("blocksynchronizer", "status", "marshal new block bytes to block data", "ERROR", err)
		return err
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, payloadStr, "", "\t")
	if err != nil {
		fmt.Println("JSON parse error: ", err)
	}
	fmt.Println("Pretty pretty print json!!:", prettyJSON.String())
	//indexer.PrintBlockInfoFromJsonBlock(newBlock.Block)

	//docID, rev, err := field.BlockStore.AddBlock(ctx, newBlock)
	blockDocID, blockDocRev, err := blockCore.AddBlock(context.Background(), newBlock)
	if err != nil {
		log.Errorw("blocksynchronizer", "status", "can't add new block", "ERROR", err)
		return err
	}
	log.Infof("Added block %s with rev %s to CouchDB Block table", blockDocID, blockDocRev)

	var accountList []models.Account
	var assetList []models.Asset
	var appList []models.Application

	if len(newBlock.Transactions) > 0 {
		_, err = transactionCore.AddTransactions(context.Background(), newBlock.Transactions)
		if err != nil {
			log.Errorw("blocksynchronizer", "status", "can't add new transaction(s)", "ERROR", err)
			return err
		}
		log.Infof("Added %d transactions with block %s to CouchDB Transaction table", len(newBlock.Transactions), newBlock.BlockHash)

		for _, txn := range newBlock.Transactions {

			accountIDs := algod2.ExtractAccountAddrsFromTxn(txn)
			applicationIDs := algod2.ExtractApplicationIdsFromTxn(txn)
			assetIDs := algod2.ExtractAssetIdsFromTxn(txn)

			for _, acctID := range accountIDs {
				accountInfo, err := algodCore.GetAccount(context.Background(),"", acctID)
				if err != nil {
					log.Errorw("blocksynchronizer", "status", "can't get account", "ERROR", err)
					//return err
				}
				log.Infof("Retrieved Account info: %v", *accountInfo)
				accountList = append(accountList, *accountInfo)
			}

			for _, appID := range applicationIDs {
				appInfo, err := algodCore.GetApplication(context.Background(),"", appID)
				if err != nil {
					log.Errorw("blocksynchronizer", "status", "can't get app", "ERROR", err)
					//return err
				}
				log.Infof("Retrieved Application info: %v", *appInfo)
				appList = append(appList, *appInfo)
			}

			for _, assetID := range assetIDs {
				assetInfo, err := algodCore.GetAsset(context.Background(),"", assetID)
				if err != nil {
					log.Errorw("blocksynchronizer", "status", "can't get asset", "ERROR", err)
					//return err
				}
				log.Infof("Retrieved Asset info: %v", *assetInfo)
				assetList = append(assetList, *assetInfo)
			}
		}
	}

	if len(accountList) > 0 {
		_, err = accountCore.AddAccounts(context.Background(), accountList)
		if err != nil {
			log.Errorw("blocksynchronizer", "status", "can't add/update account(s)", "ERROR", err)
			//return err
		}
	}

	if len(assetList) > 0 {
		_, err = assetCore.AddAssets(context.Background(), assetList)
		if err != nil {
			log.Errorw("blocksynchronizer", "status", "can't add/update asset(s)", "ERROR", err)
			//return err
		}
	}

	if len(appList) > 0 {
		_, err = appCore.AddApplications(context.Background(), appList)
		if err != nil {
			log.Errorw("blocksynchronizer", "status", "can't add/update application(s)", "ERROR", err)
			//return err
		}
	}

	return nil
}
