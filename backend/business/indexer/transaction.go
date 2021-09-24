package indexer

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	indexerv2 "github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// GetTransaction retrieves a block from the Indexer API based on the round number
// given. The difference between this method and LookupRoundInJSON is the data will
// be processed into block.NewBlock format instead of models.Block.
//func GetRound(ctx context.Context, traceID string, log *zap.SugaredLogger, indexerClient *indexerv2.Client, roundNum uint64) (*block.NewBlock, error) {
func GetTransaction(ctx context.Context, traceID string, log *zap.SugaredLogger, indexerClient *indexerv2.Client, transactionID string) (models.Transaction, error) {
	log.Infow("indexer.GetTransaction", "traceid", traceID)

	jsonBlock, err := LookupTransactionInJSON(ctx, indexerClient, transactionID)
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
func LookupTransactionInJSON(ctx context.Context, indexerClient *indexerv2.Client, transactionID string) (models.Transaction, error) {

	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "indexer.LookupTransactionInJSON")
	span.SetAttributes(attribute.String("transactionID", transactionID))
	defer span.End()

	transaction, err := indexerClient.LookupTransaction(transactionID).Do(ctx)
	if err != nil {
		return models.Transaction{}, errors.Wrap(err, "Unable to find block. Record may not exist in Postgre database.")
	}

	return transaction.Transaction, nil
}