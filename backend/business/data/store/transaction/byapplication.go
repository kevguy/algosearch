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

func (s Store) GetTransactionsByApplication(ctx context.Context, applicationID uint64, order string) ([]Transaction, error) {
	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "transaction.GetTransactionsByApplication")
	span.SetAttributes(attribute.Any("applicationID", applicationID))
	defer span.End()

	s.log.Infow("transaction.GetTransactionsByApplication",
		"traceid", web.GetTraceID(ctx),
		"applicationID", applicationID)

	exist, err := s.couchClient.DBExists(ctx, schema.GlobalDbName)
	if err != nil || !exist {
		return nil, errors.Wrap(err, schema.GlobalDbName+ " database check fails")
	}
	db := s.couchClient.DB(schema.GlobalDbName)

	options := kivik.Options{
		"include_docs": true,
		//"limit": limit,
		"start_key": fmt.Sprintf("[\"i%d\"]", applicationID),
		"end_key": fmt.Sprintf("[\"%d\", 2]", applicationID),
	}

	rows, err := db.Query(ctx, schema.TransactionDDoc, "_view/" +schema.TransactionViewByApplication, options)
	if err != nil {
		return nil, errors.Wrap(err, "Fetch data error")
	}

	var fetchedTransactions = []Transaction{}
	var count = 0
	for rows.Next() {
		if count != 0 {
			var transaction = Transaction{}
			if err := rows.ScanDoc(&transaction); err != nil {
				return nil, errors.Wrap(err, "unwrapping block")
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
