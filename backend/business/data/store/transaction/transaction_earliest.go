package transaction

import (
	"context"
	"fmt"
	"github.com/go-kivik/kivik/v4"
	"github.com/kevguy/algosearch/backend/business/data/schema"
	"github.com/kevguy/algosearch/backend/foundation/web"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
)

// GetEarliestTransactionId retrieves the latest Transaction ID that can be found in the database.
func (s Store) GetEarliestTransactionId(ctx context.Context) (string, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "transaction.GetEarliestTransactionId")
	//span.SetAttributes(attribute.String("query", q))
	defer span.End()

	s.log.Infow("transaction.GetEarliestTransactionId", "traceid", web.GetTraceID(ctx))

	exist, err := s.couchClient.DBExists(ctx, schema.GlobalDbName)
	if err != nil || !exist {
		return "", errors.Wrap(err, schema.GlobalDbName+ " database check fails")
	}
	db := s.couchClient.DB(schema.GlobalDbName)

	rows, err := db.Query(ctx, schema.TransactionDDoc, "_view/" +schema.TransactionViewInLatest, kivik.Options{
		"include_docs": true,
		"descending": false,
		"limit": 1,
	})
	if err != nil {
		return "", fmt.Errorf("fetch data error: %w", err)
	}

	if rows.Err() != nil {
		return "", fmt.Errorf("row error, can't find anything: %w", err)
	}

	rows.Next()
	var doc Transaction
	if err := rows.ScanDoc(&doc); err != nil {
		// No docs can be found
		return "", fmt.Errorf("no docs can be found: %w", err)
	}

	return doc.Id, nil
}