package algorand

import (
	"context"
	"fmt"
	algodv2 "github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	indexerv2 "github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"github.com/go-kivik/kivik/v4"
	"github.com/kevguy/algosearch/backend/business/core/algod"
	"github.com/kevguy/algosearch/backend/business/core/block"
	blockDb "github.com/kevguy/algosearch/backend/business/core/block/db"
	"github.com/kevguy/algosearch/backend/business/core/indexer"
	"github.com/kevguy/algosearch/backend/business/core/transaction"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type Agent struct {
	log *zap.SugaredLogger
	//indexerClient *indexerv2.Client
	algodClient *algodv2.Client

	algodCore 		*algod.Core
	indexerCore		*indexer.Core
	blockCore 		*block.Core
	transactionCore *transaction.Core
}

// NewAgent constructs an Algorand for api access.
func NewAgent(log *zap.SugaredLogger,
	indexerClient *indexerv2.Client,
	algodClient *algodv2.Client,
	couchClient *kivik.Client,
	dbName string) Agent {

	algodCore := algod.NewCore(log, algodClient)
	blockCore := block.NewCore(log, couchClient, dbName)
	transactionCore := transaction.NewCore(log, couchClient, dbName)

	var indexerCore *indexer.Core
	if indexerClient != nil {
		core := indexer.NewCore(log, indexerClient)
		indexerCore = &core
	}

	return Agent{
		log: log,
		algodClient: algodClient,

		algodCore: &algodCore,
		indexerCore: indexerCore,
		blockCore: &blockCore,
		transactionCore: &transactionCore,
	}
}

// GetRound retrieves a block based on the round number given
// It works by first trying the indexer, if there's a connection
// it fetches the block data from Indexer and the additional data from Couch (this sounds
// redundant, but this is written with mind of getting rid of Couch in future), if not then it tries
// Algod, and finally only Couch.
func (a Agent) GetRound(ctx context.Context, traceID string, log *zap.SugaredLogger, roundNum uint64) (*blockDb.NewBlock, error) {

	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "algorand.GetRound")
	span.SetAttributes(attribute.Int64("round", int64(roundNum)))
	defer span.End()

	log.Infow("algorand.GetRound", "traceid", traceID)

	var blockData blockDb.NewBlock
	var err error

	// Whatever we do, we still have to get data from Couch (for proposer and block hash, at least for the time being)
	couchBlock, err := a.blockCore.GetBlockByNum(ctx, roundNum)
	if err != nil {
		return nil, fmt.Errorf("unable to query couch for round %d: %w", roundNum, err)
	}

	// Try Indexer
	if a.indexerCore != nil {
		idxBlock, err := a.indexerCore.GetRound(ctx, traceID, roundNum)
		if err != nil {
			log.Errorf("unable to get block data from indexer for round %d\n", roundNum)
		} else {
			blockData = blockDb.NewBlock{
				Block:     idxBlock,
				Proposer:  couchBlock.Proposer,
				BlockHash: couchBlock.BlockHash,
			}
			return &blockData, nil
		}
	}

	// Try Algod
	if a.algodClient != nil {
		algodBlock, err := a.algodCore.GetRound(ctx, traceID, roundNum)
		if err != nil {
			log.Errorf("unable to get block data from algod for round %d\n", roundNum)

			// Use the data from Couch since Indexer and Algod are not working
			return &couchBlock.NewBlock, nil
		} else {
			blockData = *algodBlock
			return &blockData, nil
		}
	}
	return &couchBlock.NewBlock, nil
}

// GetTransaction retrieves a transaction based on the transaction ID given.
// It works by first trying the indexer, if there's a connection
// it fetches the transaction data from Indexer and the additional data from Couch (this sounds
// redundant, but this is written with mind of getting rid of Couch in future), if not
// then return the data from Couch.
func (a Agent) GetTransaction(ctx context.Context, traceID string, log *zap.SugaredLogger, transactionID string) (*models.Transaction, error) {

	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "algorand.GetTransaction")
	span.SetAttributes(attribute.String("transactionID", transactionID))
	defer span.End()

	log.Infow("algorand.GetTransaction", "traceid", traceID)

	//var transactionData models.Transaction
	var err error

	// Whatever we do, we still have to get data from Couch (for proposer and block hash, at least for the time being)
	couchTransaction, err := a.transactionCore.GetTransaction(ctx, transactionID)
	if err != nil {
		return nil, fmt.Errorf("unable to query couch for transaction %s: %w", transactionID, err)
	}

	// Try Indexer
	if a.indexerCore != nil {
		idxTransaction, err := a.indexerCore.GetTransaction(ctx, traceID, transactionID)
		if err != nil {
			log.Errorf("unable to get transaction data from indexer for transaction ID %s\n", transactionID)
		} else {
			//transactionData = transaction.Transaction{
			//	NewTransaction: transaction.NewTransaction{
			//		ID:          nil,
			//		Transaction: idxBlock,
			//		DocType:     "",
			//	},
			//	ID:             "",
			//	Rev:            "",
			//}

			//transactionData = transaction.Transaction{
			//	Transaction: idxBlock,
			//}
			//return &transactionData, nil
			return &idxTransaction, nil
		}
	}
	return &couchTransaction, nil
}
