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

func (s Store) GetTransactionsByAsset(ctx context.Context, assetID string, order string) ([]Transaction, error) {
	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "transaction.GetTransactionsByAsset")
	span.SetAttributes(attribute.String("assetID", assetID))
	defer span.End()

	s.log.Infow("transaction.GetTransactionsByAsset",
		"traceid", web.GetTraceID(ctx),
		"assetID", assetID)

	exist, err := s.couchClient.DBExists(ctx, s.dbName)
	if err != nil || !exist {
		return nil, errors.Wrap(err, s.dbName+ " database check fails")
	}
	db := s.couchClient.DB(s.dbName)

	options := kivik.Options{
		"include_docs": true,
		//"limit": limit,
		"reduce": false,
		"inclusive_end": true,
		//"start_key": []string{assetID, "1"},
		//"end_key": []string{assetID, "2"},
	}

	if order == "asc" {
		options["start_key"] = []string{assetID, "1"}
		options["end_key"] = []string{assetID, "2"}
		options["descending"] = false
	} else {
		// assuming it's "desc"
		options["start_key"] = []string{assetID, "2"}
		options["end_key"] = []string{assetID, "1"}
		options["descending"] = true
	}
	fmt.Println(options["start_key"])
	fmt.Println(options["end_key"])

	rows, err := db.Query(ctx, schema.TransactionDDoc, "_view/" +schema.TransactionViewByAsset, options)
	if err != nil {
		return nil, errors.Wrap(err, "Fetch data error")
	}

	var fetchedTransactions []Transaction
	var count = 0
	for rows.Next() {
		if count != 0 {
			var transaction = Transaction{}
			if err := rows.ScanDoc(&transaction); err != nil {
				return nil, errors.Wrap(err, "unwrasseting block")
			}
			fetchedTransactions = append(fetchedTransactions, transaction)
		}
		count += 1
	}

	if rows.Err() != nil {
		return nil, errors.Wrap(err, "rows error, Can't find anything")
	}

	return fetchedTransactions, nil
}
