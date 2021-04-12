package commands

import (
	"context"
	"fmt"
	"github.com/kevguy/algosearch/backend/business/core/block"
	"github.com/kevguy/algosearch/backend/foundation/couchdb"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

// GetRoundNumInDBCmd retrieves information about the block for the latest round and prints it out.
func GetRoundNumInDBCmd(traceID string, log *zap.SugaredLogger, couchCfg couchdb.Config, dbName string, blockNum uint64) error {

	db, err := couchdb.Open(couchCfg)
	if err != nil {
		return errors.Wrap(err, "connect to couchdb database")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer db.Close(ctx)
	defer cancel()

	blockCore := block.NewCore(log, db, dbName)

	block, err := blockCore.GetBlockByNum(ctx, blockNum)
	if err != nil {
		return errors.Wrap(err, "get block")
	}

	fmt.Printf("Block Data: %v\n", block)
	fmt.Println(block.BlockHash)


	return nil
}
