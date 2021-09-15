package commands

import (
	"context"
	"fmt"
	app "github.com/kevguy/algosearch/backend/business/algod"
	"github.com/kevguy/algosearch/backend/business/data/store/block"
	"github.com/kevguy/algosearch/backend/foundation/algod"
	"github.com/kevguy/algosearch/backend/foundation/couchdb"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

// GetCurrentRoundCmd retrieves information about the block for the latest round and prints it out.
func AddCurrentRoundCmd(traceID string, log *zap.SugaredLogger, cfg algod.Config, couchCfg couchdb.Config) error {
	client, err := algod.Open(cfg)
	if err != nil {
		return errors.Wrap(err, "connect to Algorand Node")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	num, err := app.GetCurrentRoundNum(ctx, client)
	if err != nil {
		return errors.Wrap(err, "getting current round num from Algorand Node")
	}
	rawBlock, err := app.GetRoundInRawBytes(ctx, client, num)
	if err != nil {
		return errors.Wrap(err, "getting current round from Algorand Node")
	}
	if err := app.PrintBlockInfoFromRawBytes(rawBlock); err != nil {
		return errors.Wrap(err, "process current round raw block")
	}

	db, err := couchdb.Open(couchCfg)
	if err != nil {
		return errors.Wrap(err, "connect to couchdb database")
	}

	ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second)
	defer db.Close(ctx)
	defer cancel()

	blockStore := block.NewStore(log, db)

	newBlock, err := app.ConvertBlockRawBytes(ctx, rawBlock)
	if err != nil {
		return errors.Wrap(err, "convert raw bytes to block data")
	}

	blockStore.AddBlock(ctx, newBlock)
	if err != nil {
		return errors.Wrap(err, "can't add new block")
	}

	//if err := schema.Migrate(ctx, db); err != nil {
	//	return errors.Wrap(err, "migrate couchdb database")
	//}

	fmt.Println("block added")

	return nil
}
