package transaction

import (
	"context"
	"github.com/go-kivik/kivik/v4"
	"github.com/kevguy/algosearch/backend/business/data/schema"
	"github.com/kevguy/algosearch/backend/foundation/web"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"strconv"
)

func (s Store) GetTransactionsByAcctPagination(ctx context.Context, acctID, order string, pageNo, limit int64) ([]Transaction, int64, int64, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "transaction.GetTransactionsByAcctPagination")
	span.SetAttributes(attribute.Int64("pageNo", pageNo))
	span.SetAttributes(attribute.Int64("limit", limit))
	defer span.End()

	s.log.Infow("transaction.GetTransactionsByAcctPagination",
		"traceid", web.GetTraceID(ctx),
		"pageNo", pageNo,
		"limit", limit)

	// Get the earliest transaction id
	earliestTxn, err := s.GetEarliestAcctTransaction(ctx, acctID)
	if err != nil {
		return nil, 0, 0, errors.Wrap(err, ": Get earliest synced transaction id")
	}
	earliestTxnId := earliestTxn.ID
	earliestRoundTime := earliestTxn.RoundTime

	// Get the latest transaction id
	latestTxn, err := s.GetLatestAcctTransaction(ctx, acctID)
	if err != nil {
		return nil, 0, 0, errors.Wrap(err, ": Get latest synced transaction id")
	}
	latestTxnId := latestTxn.ID
	latestRoundTime := latestTxn.RoundTime

	numOfTransactions, err := s.GetTransactionCountByAcct(ctx, acctID, earliestTxnId, latestTxnId)
	if err != nil {
		return nil, 0, 0, errors.Wrap(err, ": Get transaction count between keys")
	}
	s.log.Infow("transaction.GetTransactionsPagination",
		"numOfTransactions", numOfTransactions)

	// We can skip database check cuz GetEarliestTransactionId already did it
	db := s.couchClient.DB(schema.GlobalDbName)

	var numOfPages int64 = numOfTransactions / limit
	if numOfTransactions % limit > 0 {
		numOfPages += 1
	}

	if pageNo < 1 || pageNo > numOfPages {
		return nil, 0, 0, errors.Wrapf(err, "page number is less than 1 or exceeds page limit: %d", numOfPages)
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

		// Start with latest block number we managed to find for the time being
		options["start_key"] = []string{
			acctID,
			"1",
			strconv.FormatUint(latestRoundTime, 10),
			latestTxnId,
		}

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

		//Start with earliest block number found
		// Start with latest block number we managed to find for the time being
		options["start_key"] = []string{
			acctID,
			"1",
			strconv.FormatUint(earliestRoundTime, 10),
			earliestTxnId,
		}

		// Calculate the number of records to skip
		skip := (pageNo - 1) * limit
		options["skip"] = skip

		if (numOfTransactions - skip) < limit {
			options["limit"] =  numOfTransactions - skip
		} else {
			options["limit"] = limit
		}
	}

	s.log.Infof("transaction.GetTransactionsByAcctPagination: staritng to query")
	s.log.Infow("transaction.GetTransactionsByAcctPagination",
		"include_docs", options["include_docs"],
		"limit", options["limit"],
		"descending", options["descending"],
		"start_key", options["start_key"],
		"skip", options["skip"])
	rows, err := db.Query(ctx, schema.TransactionDDoc, "_view/" +schema.TransactionViewByAccount, options)
	if err != nil {
		return nil, 0, 0, errors.Wrap(err, "Fetch data error")
	}

	var fetchedTransactions []Transaction
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
