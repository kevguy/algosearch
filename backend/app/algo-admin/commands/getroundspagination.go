package commands

import (
	"context"
	"fmt"
	"github.com/kevguy/algosearch/backend/business/couchdata/block"
	"github.com/kevguy/algosearch/backend/foundation/couchdb"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

func GetRoundsPaginationCmd(traceID string, log *zap.SugaredLogger, couchCfg couchdb.Config, latestBlockNum, noOfItems, pageNo int64, order string) error {
	db, err := couchdb.Open(couchCfg)
	if err != nil {
		return errors.Wrap(err, ": connect to couchdb database")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer db.Close(ctx)
	defer cancel()

	blockStore := block.NewStore(log, db)

	blocks, err := blockStore.GetBlocksPagination(ctx, traceID, log, latestBlockNum, order, pageNo, noOfItems)
	for _, item := range blocks {
		fmt.Println(item.Round)
	}
	return nil
}
