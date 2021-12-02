package commands

import (
	"context"
	"fmt"
	app "github.com/kevguy/algosearch/backend/business/algod"
	algod2 "github.com/kevguy/algosearch/backend/business/core/algod"
	"github.com/kevguy/algosearch/backend/foundation/algod"
	"github.com/pkg/errors"
	"time"
)

// GetRoundCmd retrieves information about the block for the specified round and prints it out.
func GetRoundCmd(cfg algod.Config, blockNum uint64) error {
	client, err := algod.Open(cfg)
	if err != nil {
		return errors.Wrap(err, "connect to Algorand Node")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rawBlock, err := app.GetRoundInRawBytes(ctx, client, blockNum)
	if err != nil {
		return fmt.Errorf("getting rount %d from Algorand Node %w", blockNum, err)
	}
	if err := algod2.PrintBlockInfoFromRawBytes(rawBlock); err != nil {
		return fmt.Errorf("process round %d raw block from Algorand Node %w", blockNum, err)
	}

	return nil
}
