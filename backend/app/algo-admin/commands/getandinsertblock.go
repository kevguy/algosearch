package commands

import (
	"context"
	"github.com/kevguy/algosearch/backend/app/algosearch/blocksynchronizer"
	"github.com/kevguy/algosearch/backend/business/core/account"
	algod2 "github.com/kevguy/algosearch/backend/business/core/algod"
	"github.com/kevguy/algosearch/backend/business/core/application"
	"github.com/kevguy/algosearch/backend/business/core/asset"
	"github.com/kevguy/algosearch/backend/business/core/block"
	"github.com/kevguy/algosearch/backend/business/core/transaction"
	"github.com/kevguy/algosearch/backend/foundation/algod"
	"github.com/kevguy/algosearch/backend/foundation/couchdb"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

func GetAndInsertBlockCmd(log *zap.SugaredLogger, cfg algod.Config, couchCfg couchdb.Config, dbName string, fromBlock uint64, toBlock uint64) error {

	client, err := algod.Open(cfg)
	if err != nil {
		return errors.Wrap(err, "connect to Algorand Node")
	}

	algodCore := algod2.NewCore(log, client)

	db, err := couchdb.Open(couchCfg)
	if err != nil {
		return errors.Wrap(err, "connect to couchdb database")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 600000*time.Second)
	defer db.Close(ctx)
	defer cancel()

	blockCore := block.NewCore(log, db, dbName)
	transactionCore := transaction.NewCore(log, db, dbName)
	accountCore := account.NewCore(log, db, dbName)
	assetCore := asset.NewCore(log, db, dbName)
	appCore := application.NewCore(log, db, dbName)

	for i := fromBlock; i <= toBlock; i++ {
		if err := blocksynchronizer.GetAndInsertBlockData(
			log,
			client,
			&blockCore,
			&transactionCore,
			&accountCore,
			&assetCore,
			&appCore,
			&algodCore,
			i); err != nil {
			//return err
			log.Errorf("Failed to add Block Number %d\n", i)
		}
		log.Infof("Added Block Number %d\n", i)
	}

	return nil
}
