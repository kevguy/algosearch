package commands

import (
	"context"
	algod2 "github.com/kevguy/algosearch/backend/business/core/algod"
	"github.com/kevguy/algosearch/backend/foundation/algod"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

// GetCurrentRoundCmd retrieves information about the block for the latest round and prints it out.
func GetCurrentRoundCmd(log *zap.SugaredLogger, cfg algod.Config) error {
	client, err := algod.Open(cfg)
	if err != nil {
		return errors.Wrap(err, "connect to Algorand Node")
	}

	algodCore := algod2.NewCore(log, client)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	num, err := algodCore.GetCurrentRoundNum(ctx)
	if err != nil {
		return errors.Wrap(err, "getting current round num from Algorand Node")
	}
	rawBlock, err := algodCore.GetRoundInRawBytes(ctx, num)
	if err != nil {
		return errors.Wrap(err, "getting current round from Algorand Node")
	}
	if err := algod2.PrintBlockInfoFromRawBytes(rawBlock); err != nil {
		return errors.Wrap(err, "process current round raw block")
	}

	return nil
}
