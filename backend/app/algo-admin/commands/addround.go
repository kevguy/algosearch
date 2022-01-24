package commands

import (
	"context"
	"fmt"
	algod2 "github.com/kevguy/algosearch/backend/business/core/algod"
	"github.com/kevguy/algosearch/backend/business/core/block"
	"github.com/kevguy/algosearch/backend/foundation/algod"
	"github.com/kevguy/algosearch/backend/foundation/couchdb"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

// AddRoundCmd retrieves a block based on block number and inserts it into CouchDB.
func AddRoundCmd(traceID string, log *zap.SugaredLogger, cfg algod.Config, couchCfg couchdb.Config, dbName string, blockNum uint64) error {
	client, err := algod.Open(cfg)
	if err != nil {
		return errors.Wrap(err, "connect to Algorand Node")
	}

	algodCore := algod2.NewCore(log, client)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rawBlock, err := algodCore.GetRoundInRawBytes(ctx, blockNum)
	if err != nil {
		return errors.Wrap(err, "getting current round from Algorand Node")
	}
	if err := algod2.PrintBlockInfoFromRawBytes(rawBlock); err != nil {
		return errors.Wrap(err, "process current round raw block")
	}

	db, err := couchdb.Open(couchCfg)
	if err != nil {
		return errors.Wrap(err, "connect to couchdb database")
	}

	ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second)
	defer db.Close(ctx)
	defer cancel()

	blockCore := block.NewCore(log, db, dbName)

	newBlock, err := algod2.ConvertBlockRawBytes(ctx, rawBlock)
	if err != nil {
		return errors.Wrap(err, "convert raw bytes to block data")
	}

	blockCore.AddBlock(ctx, newBlock)
	if err != nil {
		return errors.Wrap(err, "can't add new block")
	}

	//if err := schema.Migrate(ctx, db); err != nil {
	//	return errors.Wrap(err, "migrate couchdb database")
	//}

	fmt.Println("block added")

	return nil
}
