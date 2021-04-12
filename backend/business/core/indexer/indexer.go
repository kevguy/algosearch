// Package indexer provides the core business API of handling
// everything account related.
package indexer

import (
	indexerv2 "github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"go.uber.org/zap"
)

// Core manages the set of API's for block access.
type Core struct {
	log *zap.SugaredLogger
	client *indexerv2.Client
}

// NewCore constructs a core for product api access.
func NewCore(log *zap.SugaredLogger, indexerClient *indexerv2.Client) Core {
	return Core{
		log: log,
		client: indexerClient,
	}
}


//func Fuck(ctx context.Context, indexerClient *indexerv2.Client, roundNum uint64) (models.Block, error) {
//}

//func ConvertBlockJSON(ctx context.Context, jsonBlock models.Block) (block.NewBlock, error) {
//
//	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "indexer.ConvertBlockJSON")
//	defer span.End()
//
//	var genesisHashStr = base64.StdEncoding.EncodeToString(jsonBlock.GenesisHash[:])
//
//	var newBLock = block.NewBlock{
//		GenesisHash:        genesisHashStr,
//		GenesisID:          jsonBlock.GenesisId,
//		PrevBlockHash:      "",
//		Rewards:            block.Rewards{},
//		Round:              0,
//		Seed:               "",
//		Timestamp:          0,
//		Transactions:       nil,
//		TransactionsRoot:   "",
//		TransactionCounter: 0,
//		UpgradeState:       block.UpgradeState{},
//		UpgradeVote:        block.UpgradeVote{},
//		Proposer:           "",
//		BlockHash:          "",
//	}
//}
