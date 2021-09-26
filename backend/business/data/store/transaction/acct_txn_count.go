package transaction

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

func (s Store) GetTransactionCountByAcct(ctx context.Context, acctID, startKey, endKey string) (int64, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "transaction.GetTransactionByAcctCountBtnKeys")
	span.SetAttributes(attribute.String("acctID", acctID))
	span.SetAttributes(attribute.String("startKey", startKey))
	span.SetAttributes(attribute.String("endKey", endKey))
	defer span.End()

	s.log.Infow("transaction.GetTransactionCountBtnKeys",
		"traceid", web.GetTraceID(ctx),
		"acctID", acctID,
		"startKey", startKey,
		"endKey", endKey)

	exist, err := s.couchClient.DBExists(ctx, schema.GlobalDbName)
	if err != nil || !exist {
		return 0, fmt.Errorf(schema.GlobalDbName + " database check fails: %w", err)
	}
	db := s.couchClient.DB(schema.GlobalDbName)

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
	rows, err := db.Query(ctx, schema.TransactionDDoc, "_view/" +schema.TransactionViewByAccountCount, kivik.Options{
		//"start_key": []string{acctID, "1", strconv.FormatUint(earliestRoundTime, 10), startKey},
		//"end_key": []string{acctID, "1", strconv.FormatUint(latestRoundTime, 10), endKey},
		"start_key": []string{acctID, strconv.FormatUint(earliestRoundTime, 10), startKey},
		"end_key": []string{acctID, strconv.FormatUint(latestRoundTime, 10), endKey},
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
