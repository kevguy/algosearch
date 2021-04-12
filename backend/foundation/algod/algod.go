package algod

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"time"
)

// Config is the required properties to use the algod client.
type Config struct {
	AlgodAddr	string
	AlgodToken	string
	KmdAddr		string
	KmdToken	string
}

// Open knows how to open a algorand connection based on the configuration.
func Open(cfg Config) (*algod.Client, error) {

	// Create an algod client
	algodClient, err := algod.MakeClient(cfg.AlgodAddr, cfg.AlgodToken)
	if err != nil {
		return nil, fmt.Errorf("failed to construct an algod client %w", err)
	}

	return algodClient, nil
}

// StatusCheck returns nil if it can successfully talk to the algod node. It
// returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, algodClient *algod.Client) error {
	// First check we can ping the database.
	for attempts := 1; attempts < 20; attempts++ {
		status, pingError := algodClient.Status().Do(ctx) // nodeStatus
		if status.LastRound >= 0 && pingError == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
	return nil
}
