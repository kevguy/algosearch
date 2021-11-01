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

// GetEarliestAppTransaction retrieves the earliest Transaction for an account that can be found in the database.
func (s Store) GetEarliestAppTransaction(ctx context.Context, appID string) (Transaction, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "transaction.GetEarliestAppTransaction")
	span.SetAttributes(attribute.String("appID", appID))
	defer span.End()

	s.log.Infow("transaction.GetEarliestAppTransaction",
		"traceid", web.GetTraceID(ctx),
		"appID", appID)

	return s.getEarliestLatestAppTransaction(ctx, appID, true)
}

// GetLatestAppTransaction retrieves the latest Transaction for an account that can be found in the database.
func (s Store) GetLatestAppTransaction(ctx context.Context, appID string) (Transaction, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "transaction.GetLatestAppTransaction")
	span.SetAttributes(attribute.String("appID", appID))
	defer span.End()

	s.log.Infow("transaction.GetLatestAppTransaction",
		"traceid", web.GetTraceID(ctx),
		"appID", appID)

	return s.getEarliestLatestAppTransaction(ctx, appID, false)

}

// getEarliestLatestAppTransaction retrieves the latest/earliest Transaction for an account that can be found
// in the database, depending on how you define the `earliest` parameter.
func (s Store) getEarliestLatestAppTransaction(ctx context.Context, appID string, earliest bool) (Transaction, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "transaction.getEarliestLatestAppTransaction")
	span.SetAttributes(attribute.String("appID", appID))
	defer span.End()

	s.log.Infow("transaction.getEarliestLatestAppTransaction",
		"traceid", web.GetTraceID(ctx),
		"appID", appID)

	exist, err := s.couchClient.DBExists(ctx, schema.GlobalDbName)
	if err != nil || !exist {
		return Transaction{}, errors.Wrap(err, schema.GlobalDbName+ " database check fails")
	}
	db := s.couchClient.DB(schema.GlobalDbName)

	options := kivik.Options{
		"include_docs": true,
		"limit": 1,
		"start_key": fmt.Sprintf("\"[%s, 1]\"", appID),
	}

	if earliest == true {
		options["descending"] = false
	} else {
		options["descending"] = true
	}

	rows, err := db.Query(ctx, schema.TransactionDDoc, "_view/" + schema.TransactionViewByApplication, options)
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
