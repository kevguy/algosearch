package commands

import (
	"context"
	"fmt"
	"github.com/kevguy/algosearch/backend/business/data/store/block"
	"github.com/kevguy/algosearch/backend/foundation/couchdb"
	"go.uber.org/zap"
	"time"
)

func GetBlocksCountFromDbCmd(log *zap.SugaredLogger, couchCfg couchdb.Config) error {

	db, err := couchdb.Open(couchCfg)
	if err != nil {
		return fmt.Errorf("connect to couchdb database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	blockStore := block.NewStore(log, db)

	count, err := blockStore.GetNumOfBlocks(ctx)
	if err != nil {
		return fmt.Errorf("getting number of blocks in db: %w", err)
	}

	fmt.Println("=====================================================")
	fmt.Printf("Number of blocks found in DB: %d\n", count)
	fmt.Println("=====================================================")

	return nil
}
