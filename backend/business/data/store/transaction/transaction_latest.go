package transaction

import (
	"context"
	"fmt"
	"github.com/go-kivik/kivik/v4"
	"github.com/kevguy/algosearch/backend/business/data/schema"
	"github.com/kevguy/algosearch/backend/foundation/web"
	"go.opentelemetry.io/otel"
)

// GetLatestTransactionId retrieves the latest Transaction ID that can be found in the database.
func (s Store) GetLatestTransactionId(ctx context.Context) (string, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "transaction.GetLatestTransactionId")
	//span.SetAttributes(attribute.String("query", q))
	defer span.End()

	s.log.Infow("transaction.GetLatestTransactionId", "traceid", web.GetTraceID(ctx))

	exist, err := s.couchClient.DBExists(ctx, schema.GlobalDbName)
	if err != nil || !exist {
		return "", fmt.Errorf(schema.GlobalDbName + " database check fails: %w", err)
	}
	db := s.couchClient.DB(schema.GlobalDbName)

	rows, err := db.Query(ctx, schema.TransactionDDoc, "_view/" +schema.TransactionViewInLatest, kivik.Options{
		"include_docs": true,
		"descending": true,
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
