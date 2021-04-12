package indexer

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// GetTransaction retrieves a block from the Indexer API based on the round number
// given. The difference between this method and LookupRoundInJSON is the data will
// be processed into block.NewBlock format instead of models.Block.
//func GetRound(ctx context.Context, traceID string, log *zap.SugaredLogger, indexerClient *indexerv2.Client, roundNum uint64) (*block.NewBlock, error) {
func (c Core) GetTransaction(ctx context.Context, traceID string, transactionID string) (models.Transaction, error) {
	c.log.Infow("indexer.GetTransaction", "traceid", traceID)

	jsonBlock, err := c.LookupTransactionInJSON(ctx, transactionID)
	if err != nil {
		return models.Transaction{}, errors.Wrap(err, "Unable to look up block")
	}

	//blockData, err := ConvertBlockJSON(ctx, jsonBlock)
	//if err != nil {
	//	return nil, errors.Wrap(err, "Unable to look up block")
	//}
	//return &blockData, nil
	return jsonBlock, nil
}

// LookupTransactionInJSON searches for a transaction from the Indexer API based upon the
// transaction ID given.
func (c Core) LookupTransactionInJSON(ctx context.Context, transactionID string) (models.Transaction, error) {

	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "indexer.LookupTransactionInJSON")
	span.SetAttributes(attribute.String("transactionID", transactionID))
	defer span.End()

	transaction, err := c.client.LookupTransaction(transactionID).Do(ctx)
	if err != nil {
		return models.Transaction{}, errors.Wrap(err, "Unable to find block. Record may not exist in Postgre database.")
	}

	return transaction.Transaction, nil
}
