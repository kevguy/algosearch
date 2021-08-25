// Package transaction contains transaction related CRUD functionality.
package transaction

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/go-kivik/kivik/v4"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

const (
	BlocksDb = "algo_global"
	DocType = "transaction"
)

type Store struct {
	log *zap.SugaredLogger
	couchClient *kivik.Client
}

// NewStore constructs a product store for api access.
func NewStore(log *zap.SugaredLogger, couchClient *kivik.Client) Store {
	return Store{
		log: log,
		couchClient: couchClient,
	}
}

// AddTransaction adds a transaction to CouchDB.
// It receives the models.Transaction object and transform it into a Transaction document object and then
// insert it into the global CouchDB table.
func (s Store) AddTransaction(ctx context.Context, transaction models.Transaction) (string, string, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "transaction.AddTransaction")
	defer span.End()

	var doc = NewTransaction{
		Transaction: transaction,
		DocType:     DocType,
	}
	docId := fmt.Sprintf("%s.%s", DocType, doc.Id)
	exist, err := s.couchClient.DBExists(ctx, BlocksDb)
	if err != nil || !exist {
		return "", "", errors.Wrap(err, BlocksDb+ " database check fails")
	}
	db := s.couchClient.DB(BlocksDb)

	rev, err := db.Put(ctx, docId, doc)
	if err != nil {
		return "", "", errors.Wrap(err, BlocksDb+ " database can't insert transaction id " + doc.Id)
	}
	return docId, rev, nil
}

// GetTransaction adds a retrieves a transaction record from CouchDB based upon the transaction ID given.
func (s Store) GetTransaction(ctx context.Context, transactionID string) (models.Transaction, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "transaction.GetTransaction")
	defer span.End()

	exist, err := s.couchClient.DBExists(ctx, BlocksDb)
	if err != nil || !exist {
		return models.Transaction{}, errors.Wrap(err, BlocksDb + " database check fails")
	}
	db := s.couchClient.DB(BlocksDb)

	docId := fmt.Sprintf("%s.%s", DocType, transactionID)
	row := db.Get(ctx, docId)
	if row == nil {
		return models.Transaction{}, errors.Wrap(err, BlocksDb + " get data empty")
	}

	var transaction Transaction
	fmt.Printf("%v\n", row)
	err = row.ScanDoc(&transaction)
	if err != nil {
		return models.Transaction{}, errors.Wrap(err, BlocksDb + "cannot unpack data from row")
	}

	return transaction.Transaction, nil
}

// GetBlocksPagination retrieves a list of blocks based upon the following parameters:
// latestBlockNum: the latest block number that user knows about
// order: desc/asc
// pageNo: the number of pages the user wants to look at
// limit: number of blocks per page
// https://docs.couchdb.org/en/main/ddocs/views/pagination.html
/**
func (s Store) GetTransactionsPagination(ctx context.Context, traceID string, log *zap.SugaredLogger, latestTransactionID string, order string, pageNo int64, limit int64) ([]Transaction, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "transaction.GetTransactionsPagination")
	span.SetAttributes(attribute.String("latestTransactionID", latestTransactionID))
	span.SetAttributes(attribute.Int64("pageNo", pageNo))
	span.SetAttributes(attribute.Int64("limit", limit))
	defer span.End()

	exist, err := s.couchClient.DBExists(ctx, BlocksDb)
	if err != nil || !exist {
		return nil, errors.Wrap(err, BlocksDb+ " database check fails")
	}
	db := s.couchClient.DB(BlocksDb)

	// We can basically treat latestBlockNum as number of blocks
	var numOfPages int64 = latestBlockNum / pageNo
	if latestBlockNum % pageNo > 0 {
		numOfPages += 1
	}

	if pageNo < 1 || pageNo > numOfPages {
		return nil, errors.Wrapf(err, "page number is less than 1 or exceeds page limit: %d", numOfPages)
	}

	options := kivik.Options{
		"include_docs": true,
		"limit": limit,
	}

	if order == "desc" {
		options["descending"] = true
		options["start_key"] = latestBlockNum
		skip := (pageNo - 1) * limit
		options["skip"] = (pageNo - 1) * limit
		if (latestBlockNum - skip) > limit {
			options["limit"] = limit
		} else {
			options["limit"] = latestBlockNum - skip
		}
	} else {
		options["descending"] = false
		skip := (pageNo - 1) * limit
		options["skip"] = skip
		if (skip + limit - latestBlockNum) > 0 {
			options["limit"] =  latestBlockNum - skip
		} else {
			options["limit"] = limit
		}
	}

	//rows, err := db.Query(ctx, "_design/latest", "_view/latest", kivik.Options{
	//	"include_docs": true,
	//	"descending": true,
	//	"limit": limit,
	//	"skip": lastBlockNum - limit,
	//})
	rows, err := db.Query(ctx, "_design/latest", "_view/latest", options)
	if err != nil {
		return nil, errors.Wrap(err, "Fetch data error")
	}

	var fetchedBlocks = []Block{}
	for rows.Next() {
		var block = Block{}
		if err := rows.ScanDoc(&block); err != nil {
			return nil, errors.Wrap(err, "unwrapping block")
		}
		fetchedBlocks = append(fetchedBlocks, block)
	}

	if rows.Err() != nil {
		return nil, errors.Wrap(err, "rows error, Can't find anything")
	}

	return fetchedBlocks, nil
}
*/
