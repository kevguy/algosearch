package commands

import (
	"context"
	app "github.com/kevguy/algosearch/backend/business/algorand"
	"github.com/kevguy/algosearch/backend/foundation/algorand"
	"github.com/pkg/errors"
	"time"
)

// GetRoundCmd retrieves information about the block for the specified round and prints it out.
func GetRoundCmd(cfg algorand.Config, blockNum uint64) error {
	client, err := algorand.Open(cfg)
	if err != nil {
		return errors.Wrap(err, "connect to Algorand Node")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rawBlock, err := app.GetRoundInRawBytes(ctx, client, blockNum)
	if err != nil {
		return errors.Wrap(err, "getting current round from Algorand Node")
	}
	if err := app.PrintBlockInfoFromRawBytes(rawBlock); err != nil {
		return errors.Wrap(err, "process current round raw block")
	}

	return nil
}