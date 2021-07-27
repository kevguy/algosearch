package commands

import (
	"context"
	app "github.com/kevguy/algosearch/backend/business/algod"
	"github.com/kevguy/algosearch/backend/foundation/algod"
	"github.com/pkg/errors"
	"time"
)

// GetCurrentRoundCmd retrieves information about the block for the latest round and prints it out.
func GetCurrentRoundCmd(cfg algod.Config) error {
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

	return nil
}
