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

// GetRoundInDbCmd retrieves information about the block for the latest round and prints it out.
func GetRoundInDbCmd(traceID string, log *zap.SugaredLogger, couchCfg couchdb.Config, blockNum uint64) error {

	db, err := couchdb.Open(couchCfg)
	if err != nil {
		return errors.Wrap(err, "connect to couchdb database")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer db.Close(ctx)
	defer cancel()

	blockStore := block.NewStore(log, db)

	block, err := blockStore.GetBlock(ctx, uint64(blockNum))
	if err != nil {
		return errors.Wrap(err, "get block")
	}

	fmt.Printf("Block Data: %v\n", block)


	return nil
}
