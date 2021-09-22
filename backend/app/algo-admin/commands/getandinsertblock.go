package commands

import (
	"context"
	"github.com/kevguy/algosearch/backend/app/algosearch/blocksynchronizer"
	"github.com/kevguy/algosearch/backend/business/data/store/account"
	"github.com/kevguy/algosearch/backend/business/data/store/application"
	"github.com/kevguy/algosearch/backend/business/data/store/asset"
	"github.com/kevguy/algosearch/backend/business/data/store/block"
	"github.com/kevguy/algosearch/backend/business/data/store/transaction"
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

	blockStore := block.NewStore(log, db)
	transactionStore := transaction.NewStore(log, db)
	accountStore := account.NewStore(log, db)
	assetStore := asset.NewStore(log, db)
	appStore := application.NewStore(log, db)

	for i := fromBlock; i <= toBlock; i++ {
		if err := blocksynchronizer.GetAndInsertBlockData(
			log,
			client,
			&blockStore,
			&transactionStore,
			&accountStore,
			&assetStore,
			&appStore,
			i); err != nil {
			return err
		}
		log.Infof("Added Block Number %d\n", i)
	}

	return nil
}
