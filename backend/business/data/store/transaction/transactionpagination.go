package transaction

import (
	"context"
	"github.com/go-kivik/kivik/v4"
	"github.com/kevguy/algosearch/backend/business/data/schema"
	"github.com/kevguy/algosearch/backend/foundation/web"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

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

	exist, err := s.couchClient.DBExists(ctx, schema.GlobalDbName)
	if err != nil || !exist {
		return 0, errors.Wrap(err, schema.GlobalDbName+ " database check fails")
	}
	db := s.couchClient.DB(schema.GlobalDbName)

	rows, err := db.Query(ctx, schema.TransactionDDoc, "_view/" +schema.TransactionViewByIdCount, kivik.Options{
		"start_key": startKey,
		"end_key": endKey,
	})
	if err != nil {
		return 0, errors.Wrap(err, "Fetch data error")
	}

	type Payload struct {
		Key *string `json:"key"`
		Value int64 `json:"value"`
	}

	var payload Payload
	for rows.Next() {
		if err := rows.ScanDoc(&payload); err != nil {
			return 0, errors.Wrap(err, "Can't find anything")
		}
	}

	return payload.Value, nil
}

func (s Store) GetTransactionsPagination(ctx context.Context, latestTransactionId, order string, pageNo, limit int64) ([]Transaction, int64, int64, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "block.GetBlocksPagination")
	span.SetAttributes(attribute.String("latestTransactionId", latestTransactionId))
	span.SetAttributes(attribute.Int64("pageNo", pageNo))
	span.SetAttributes(attribute.Int64("limit", limit))
	defer span.End()

	s.log.Infow("transaction.GetTransactionsPagination",
		"traceid", web.GetTraceID(ctx),
		"latestTranasctionId", latestTransactionId,
		"pageNo", pageNo,
		"limit", limit)

	// Get the earliest transaction id
	earliestTxnId, err := s.GetEarliestTransactionId(ctx)
	if err != nil {
		return nil, 0, 0, errors.Wrap(err, ": Get earliest synced transaction id")
	}

	numOfTransactions, err := s.GetTransactionCountBtnKeys(ctx, earliestTxnId, latestTransactionId)
	if err != nil {
		return nil, 0, 0, errors.Wrap(err, ": Get transaction count between keys")
	}

	// We can skip database check cuz GetEarliestTransactionId already did it
	db := s.couchClient.DB(schema.GlobalDbName)

	var numOfPages int64 = numOfTransactions / limit
	if numOfTransactions % limit > 0 {
		numOfPages += 1
	}

	if pageNo < 1 || pageNo > numOfPages {
		return nil, 0, 0, errors.Wrapf(err, "page number is less than 1 or exceeds page limit: %d", numOfPages)
	}

	options := kivik.Options{
		"include_docs": true,
		"limit": limit,
	}

	if order == "desc" {
		// Descending order
		options["descending"] = true

		// Start with latest block number
		options["start_key"] = latestTransactionId

		// Use page number to calculate number of items to skip
		skip := (pageNo - 1) * limit
		options["skip"] = (pageNo - 1) * limit

		// Find the key to start reading and get the `page limit` number of records
		if (numOfTransactions - skip) > limit {
			options["limit"] = limit
		} else {
			options["limit"] = numOfTransactions - skip
		}
	} else {
		// Ascending order
		options["descending"] = false

		// Calculate the number of records to skip
		skip := (pageNo - 1) * limit
		options["skip"] = skip

		if (numOfTransactions - skip) > limit {
			options["limit"] =  numOfTransactions - skip
		} else {
			options["limit"] = limit
		}
	}

	rows, err := db.Query(ctx, schema.TransactionDDoc, "_view/" +schema.TransactionViewInLatest, options)
	if err != nil {
		return nil, 0, 0, errors.Wrap(err, "Fetch data error")
	}

	var fetchedTransactions = []Transaction{}
	for rows.Next() {
		var transaction = Transaction{}
		if err := rows.ScanDoc(&transaction); err != nil {
			return nil, 0, 0, errors.Wrap(err, "unwrapping block")
		}
		fetchedTransactions = append(fetchedTransactions, transaction)
	}

	if rows.Err() != nil {
		return nil, 0, 0, errors.Wrap(err, "rows error, Can't find anything")
	}

	return fetchedTransactions, numOfPages, numOfTransactions, nil
}
