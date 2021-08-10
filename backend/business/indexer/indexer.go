package indexer

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	indexerv2 "github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"github.com/kevguy/algosearch/backend/business/couchdata/block"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"strconv"
)

// GetRound retrieves a block from the Indexer API based on the round number
// given. The difference between this method and LookupRoundInJSON is the data will
// be processed into block.NewBlock format instead of models.Block.
func GetRound(ctx context.Context, traceID string, log *zap.SugaredLogger, indexerClient *indexerv2.Client, roundNum uint64) (*block.NewBlock, error) {
	log.Infow("indexer.GetRound", "traceid", traceID)

	jsonBlock, err := LookupRoundInJSON(ctx, indexerClient, roundNum)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to look up block")
	}

	blockData, err := ConvertBlockJSON(ctx, jsonBlock)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to look up block")
	}
	return &blockData, nil
}

// LookupRoundInJSON searches for a block from the Indexer API based upon the
// round number given.
func LookupRoundInJSON(ctx context.Context, indexerClient *indexerv2.Client, roundNum uint64) (models.Block, error) {

	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "indexer.LookupRoundInJSON")
	span.SetAttributes(attribute.String("blockNum", strconv.FormatUint(roundNum, 10)))
	defer span.End()

	block, err := indexerClient.LookupBlock(roundNum).Do(ctx)
	if err != nil {
		return models.Block{}, errors.Wrap(err, "Unable to find block. Record may not exist in Postgre database.")
	}

	return block, nil
}

func ConvertBlockJSON(ctx context.Context, jsonBlock models.Block) (block.NewBlock, error) {

	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "indexer.ConvertBlockJSON")
	defer span.End()

	var genesisHashStr = base64.StdEncoding.EncodeToString(jsonBlock.GenesisHash[:])

	var newBLock = block.NewBlock{
		GenesisHash:        genesisHashStr,
		GenesisID:          jsonBlock.GenesisId,
		PrevBlockHash:      "",
		Rewards:            block.Rewards{},
		Round:              0,
		Seed:               "",
		Timestamp:          0,
		Transactions:       nil,
		TransactionsRoot:   "",
		TransactionCounter: 0,
		UpgradeState:       block.UpgradeState{},
		UpgradeVote:        block.UpgradeVote{},
		Proposer:           "",
		BlockHash:          "",
	}
}
