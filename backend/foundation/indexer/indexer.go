package indexer

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"github.com/pkg/errors"
	"time"
)

type Config struct {
	IndexerAddr		string
	IndexerToken	string
}

func Open(cfg Config) (*indexer.Client, error) {

	// Create an indexer client
	indexerClient, err := indexer.MakeClient(cfg.IndexerAddr, cfg.IndexerToken)
	if err != nil {
		return nil, errors.Wrap(err, "failed to construct an algod client")
	}

	return indexerClient, nil
}

// StatusCheck returns nil if it can successfully talk to the indexer node. It
// returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, indexerClient *indexer.Client) error {
	// First check we can ping the database.
	for attempts := 1; attempts < 20; attempts++ {
		healthCheck, err := indexerClient.HealthCheck().Do(ctx) // nodeStatus
		if err == nil {
			break
		}
		if healthCheck.Errors != nil && len(healthCheck.Errors) == 0 {
			break
		}
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
	return nil
}
