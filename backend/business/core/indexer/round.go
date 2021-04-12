package indexer

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"strconv"
)

// GetRound retrieves a block from the Indexer API based on the round number
// given. The difference between this method and LookupRoundInJSON is the data will
// be processed into block.NewBlock format instead of models.Block.
//func GetRound(ctx context.Context, traceID string, log *zap.SugaredLogger, indexerClient *indexerv2.Client, roundNum uint64) (*block.NewBlock, error) {
func (c Core) GetRound(ctx context.Context, traceID string, roundNum uint64) (models.Block, error) {
	c.log.Infow("indexer.GetRound", "traceid", traceID)

	jsonBlock, err := c.LookupRoundInJSON(ctx, roundNum)
	if err != nil {
		return models.Block{}, errors.Wrap(err, "Unable to look up block")
	}

	//blockData, err := ConvertBlockJSON(ctx, jsonBlock)
	//if err != nil {
	//	return nil, errors.Wrap(err, "Unable to look up block")
	//}
	//return &blockData, nil
	return jsonBlock, nil
}

// LookupRoundInJSON searches for a block from the Indexer API based upon the
// round number given.
func (c Core) LookupRoundInJSON(ctx context.Context, roundNum uint64) (models.Block, error) {

	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "indexer.LookupRoundInJSON")
	span.SetAttributes(attribute.String("blockNum", strconv.FormatUint(roundNum, 10)))
	defer span.End()

	block, err := c.client.LookupBlock(roundNum).Do(ctx)
	if err != nil {
		return models.Block{}, errors.Wrap(err, "Unable to find block. Record may not exist in Postgre database.")
	}

	return block, nil
}
