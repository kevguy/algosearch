package db

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

// GetEarliestAcctTransaction retrieves the earliest Transaction for an account that can be found in the database.
func (s Store) GetEarliestAcctTransaction(ctx context.Context, acctID string) (Transaction, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "transaction.GetEarliestAcctTransaction")
	span.SetAttributes(attribute.String("acctID", acctID))
	defer span.End()

	s.log.Infow("transaction.GetEarliestAcctTransaction",
		"traceid", web.GetTraceID(ctx),
		"acctID", acctID)

	return s.getEarliestLatestAcctTransaction(ctx, acctID, true)
}

// GetLatestAcctTransaction retrieves the latest Transaction for an account that can be found in the database.
func (s Store) GetLatestAcctTransaction(ctx context.Context, acctID string) (Transaction, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "transaction.GetLatestAcctTransaction")
	span.SetAttributes(attribute.String("acctID", acctID))
	defer span.End()

	s.log.Infow("transaction.GetLatestAcctTransaction",
		"traceid", web.GetTraceID(ctx),
		"acctID", acctID)

	return s.getEarliestLatestAcctTransaction(ctx, acctID, false)

}

// getEarliestLatestAcctTransaction retrieves the latest/earliest Transaction for an account that can be found
// in the database, depending on how you define the `earliest` parameter.
func (s Store) getEarliestLatestAcctTransaction(ctx context.Context, acctID string, earliest bool) (Transaction, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "transaction.getEarliestLatestAcctTransaction")
	span.SetAttributes(attribute.String("acctID", acctID))
	defer span.End()

	s.log.Infow("transaction.getEarliestLatestAcctTransaction",
		"traceid", web.GetTraceID(ctx),
		"acctID", acctID)

	exist, err := s.couchClient.DBExists(ctx, s.dbName)
	if err != nil || !exist {
		return Transaction{}, errors.Wrap(err, s.dbName+ " database check fails")
	}
	db := s.couchClient.DB(s.dbName)

	options := kivik.Options{
		"include_docs": true,
		"limit": 1,
		//"inclusive_end": true,
		//"start_key": []string{acctID, "1"},
		//"end_key": []string{acctID, "2"},
	}

	if earliest {
		options["start_key"] = []string{acctID, "1"}
		options["descending"] = false
	} else {
		options["start_key"] = []string{acctID, "2"}
		options["descending"] = true
	}

	rows, err := db.Query(ctx, schema.TransactionDDoc, "_view/" + schema.TransactionViewByAccount, options)
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
