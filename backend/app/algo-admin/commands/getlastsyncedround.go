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

// GetLastSyncedRoundCmd retrieves information about the block for the latest round and prints it out.
func GetLastSyncedRoundCmd(traceID string, log *zap.SugaredLogger, couchCfg couchdb.Config) error {

	db, err := couchdb.Open(couchCfg)
	if err != nil {
		return errors.Wrap(err, "connect to couchdb database")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer db.Close(ctx)
	defer cancel()

	blockCore := block.NewCore(log, db)

	num, err := blockCore.GetLastSyncedRoundNumber(ctx)
	if err != nil {
		return errors.Wrap(err, "can't get last synced block num")
	}

	fmt.Printf("The last synced block number is %d\n", num)

	return nil
}
