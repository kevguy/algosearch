package commands

import (
	"context"
	"github.com/kevguy/algosearch/backend/app/algosearch/blocksynchronizer"
	"github.com/kevguy/algosearch/backend/business/core/account"
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

func GetAndInsertBlockCmd(log *zap.SugaredLogger, cfg algod.Config, couchCfg couchdb.Config, fromBlock uint64, toBlock uint64) error {

	client, err := algod.Open(cfg)
	if err != nil {
		return errors.Wrap(err, "connect to Algorand Node")
	}

	db, err := couchdb.Open(couchCfg)
	if err != nil {
		return errors.Wrap(err, "connect to couchdb database")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 600000*time.Second)
	defer db.Close(ctx)
	defer cancel()

	blockCore := block.NewCore(log, db)
	transactionCore := transaction.NewCore(log, db)
	accountCore := account.NewCore(log, db)
	assetCore := asset.NewCore(log, db)
	appCore := application.NewCore(log, db)

	for i := fromBlock; i <= toBlock; i++ {
		if err := blocksynchronizer.GetAndInsertBlockData(
			log,
			client,
			&blockCore,
			&transactionCore,
			&accountCore,
			&assetCore,
			&appCore,
			i); err != nil {
			//return err
			log.Errorf("Failed to add Block Number %d\n", i)
		}
		log.Infof("Added Block Number %d\n", i)
	}

	return nil
}
