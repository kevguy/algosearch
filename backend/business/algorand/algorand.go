package algorand

import (
	"context"
	algodv2 "github.com/algorand/go-algorand-sdk/client/v2/algod"
	indexerv2 "github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"github.com/kevguy/algosearch/backend/business/algod"
	"github.com/kevguy/algosearch/backend/business/couchdata/block"
	"github.com/kevguy/algosearch/backend/business/indexer"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type Agent struct {
	log *zap.SugaredLogger
	indexerClient *indexerv2.Client
	algodClient *algodv2.Client
	blockStore *block.Store
}

// NewAgent constructs an Algorand for api access.
func NewAgent(log *zap.SugaredLogger, indexerClient *indexerv2.Client, algodClient *algodv2.Client, blockStore *block.Store) Agent {
	return Agent{
		log: log,
		indexerClient: indexerClient,
		algodClient: algodClient,
		blockStore: blockStore,
	}
}

// GetRound retrieves a block based on the round number given
// It works by first trying the indexer, if there's a connection
// it fetches the block data from Indexer and the additional data from Couch (this sounds
// redundant, but this is written with mind of getting rid of Couch in future), if not then it tries
// Algod, and finally only Couch.
func (a Agent) GetRound(ctx context.Context, traceID string, log *zap.SugaredLogger, roundNum uint64) (*block.NewBlock, error) {

	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "algorand.GetRound")
	span.SetAttributes(attribute.Int64("round", int64(roundNum)))
	defer span.End()

	log.Infow("algorand.GetRound", "traceid", traceID)

	var blockData block.NewBlock
	var err error

	// Whatever we do, we still have to get data from Couch (for proposer and block hash, at least for the time being)
	couchBlock, err := a.blockStore.GetBlockByNum(ctx, traceID, log, roundNum)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to query couch for round %d", roundNum)
	}

	// Try Indexer
	if a.indexerClient != nil {
		idxBlock, err := indexer.GetRound(ctx, traceID, log, a.indexerClient, roundNum)
		if err != nil {
			log.Errorf("unable to get block data from indexer for round %d\n", roundNum)
		} else {
			blockData = block.NewBlock{
				Block:     idxBlock,
				Proposer:  couchBlock.Proposer,
				BlockHash: couchBlock.BlockHash,
			}
			return &blockData, nil
		}
	}

	// Try Algod
	if a.algodClient != nil {
		algodBlock, err := algod.GetRound(ctx, traceID, log, a.algodClient, roundNum)
		if err != nil {
			log.Errorf("unable to get block data from algod for round %d\n", roundNum)

			// Use the data from Couch since Indexer and Algod are not working
			return &couchBlock.NewBlock, nil
		} else {
			blockData = *algodBlock
			return &blockData, nil
		}
	}
	return &blockData, nil
}