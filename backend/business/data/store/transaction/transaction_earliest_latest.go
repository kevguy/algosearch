package transaction

import (
	"context"
	"fmt"
	"github.com/go-kivik/kivik/v4"
	"github.com/kevguy/algosearch/backend/business/data/schema"
	"github.com/kevguy/algosearch/backend/foundation/web"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// GetEarliestTransaction retrieves the latest Transaction that can be found in the database.
func (s Store) GetEarliestTransaction(ctx context.Context) (Transaction, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "transaction.GetEarliestTransaction")
	//span.SetAttributes(attribute.String("query", q))
	defer span.End()

	s.log.Infow("transaction.GetEarliestTransaction", "traceid", web.GetTraceID(ctx))

	return s.getEarliestLatestTransaction(ctx, true)
}

// GetLatestTransaction retrieves the latest Transaction that can be found in the database.
func (s Store) GetLatestTransaction(ctx context.Context) (Transaction, error) {
	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "transaction.GetLatestTransactionId")
	//span.SetAttributes(attribute.String("query", q))
	defer span.End()

	s.log.Infow("transaction.GetLatestTransactionId", "traceid", web.GetTraceID(ctx))

	return s.getEarliestLatestTransaction(ctx, false)
}

// GetEarliestTransaction retrieves the latest Transaction that can be found in the database.
func (s Store) getEarliestLatestTransaction(ctx context.Context, earliest bool) (Transaction, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "transaction.getEarliestLatestTransaction")
	span.SetAttributes(attribute.Bool("earliest", earliest))
	defer span.End()

	s.log.Infow("transaction.getEarliestLatestTransaction",
		"traceid", web.GetTraceID(ctx),
		"earliest", earliest)

	exist, err := s.couchClient.DBExists(ctx, schema.GlobalDbName)
	if err != nil || !exist {
		return Transaction{}, errors.Wrap(err, schema.GlobalDbName+ " database check fails")
	}
	db := s.couchClient.DB(schema.GlobalDbName)

	options := kivik.Options{
		"include_docs": true,
		"limit": 1,
	}

	if earliest == true {
		options["descending"] = false
	} else {
		options["descending"] = true
	}

	rows, err := db.Query(ctx, schema.TransactionDDoc, "_view/" +schema.TransactionViewInLatest, options)
	if err != nil {
		return Transaction{}, fmt.Errorf("fetch data error: %w", err)
	}

	if rows.Err() != nil {
		return Transaction{}, fmt.Errorf("row error, can't find anything: %w", err)
	}

	rows.Next()
	var doc Transaction
	if err := rows.ScanDoc(&doc); err != nil {
		// No docs can be found
		return Transaction{}, fmt.Errorf("no docs can be found: %w", err)
	}

	return doc, nil
}
