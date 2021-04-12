package db

import (
	"context"
	"fmt"
	"github.com/go-kivik/kivik/v4"
	"github.com/kevguy/algosearch/backend/business/data/schema"
	"github.com/kevguy/algosearch/backend/foundation/web"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"strconv"
)

// GetTransactionCountBtnKeys gets the count between a transaction and another. The transactions are arranged
// in chronological order in the view.
// https://stackoverflow.com/questions/11284383/couchdb-count-unique-document-field
// https://stackoverflow.com/questions/12944294/using-a-couchdb-view-can-i-count-groups-and-filter-by-key-range-at-the-same-tim
func (s Store) GetTransactionCountBtnKeys(ctx context.Context, startKey, endKey string) (int64, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "transaction.GetTransactionCountBtnKeys")
	span.SetAttributes(attribute.String("startKey", startKey))
	span.SetAttributes(attribute.String("endKey", endKey))
	defer span.End()

	s.log.Infow("transaction.GetTransactionCountBtnKeys",
		"traceid", web.GetTraceID(ctx),
		"startKey", startKey,
		"endKey", endKey)

	exist, err := s.couchClient.DBExists(ctx, s.dbName)
	if err != nil || !exist {
		return 0, fmt.Errorf(s.dbName + " database check fails: %w", err)
	}
	db := s.couchClient.DB(s.dbName)

	// Get transaction by ID
	earliestTxn, err := s.GetTransaction(ctx, startKey)
	if err != nil {
		return 0, fmt.Errorf("fetch earliest transaction error: %w", err)
	}
	earliestRoundTime := earliestTxn.RoundTime

	latestTxn, err := s.GetTransaction(ctx, endKey)
	if err != nil {
		return 0, fmt.Errorf("fetch latest transaction error: %w", err)
	}
	latestRoundTime := latestTxn.RoundTime

	// https://github.com/go-kivik/kivik/issues/246
	rows, err := db.Query(ctx, schema.TransactionDDoc, "_view/" +schema.TransactionViewByIDCount, kivik.Options{
		"start_key": []string{strconv.FormatUint(earliestRoundTime, 10), startKey},
		"end_key": []string{strconv.FormatUint(latestRoundTime, 10), endKey},
		"inclusive_end": true,
		"reduce": true,
		"group_level": 0,
		"skip": 0,
		//"limit": 101,
	})
	if err != nil {
		return 0, fmt.Errorf("fetch data error: %w", err)
	}

	var count int64
	for rows.Next() {
		if err := rows.ScanValue(&count); err != nil {
			return 0, fmt.Errorf("can't find anything: %w", err)
		}
	}

	return count, nil
}
