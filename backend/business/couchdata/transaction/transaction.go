// Package transaction contains transaction related CRUD functionality.
package transaction

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/go-kivik/kivik/v4"
	app "github.com/kevguy/algosearch/backend/business/algod"
	"github.com/kevguy/algosearch/backend/business/couchdata/schema"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

const (
	DocType = "txn"
)

type Store struct {
	log *zap.SugaredLogger
	couchClient *kivik.Client
}

// NewStore constructs a transaction store for api access.
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
		AssociatedAccounts: app.ExtractAccountAddrsFromTxn(transaction),
		AssociatedApplications: app.ExtractApplicationIdsFromTxn(transaction),
		AssociatedAssets: app.ExtractAssetIdsFromTxn(transaction),
	}
	//docId := fmt.Sprintf("%s.%s", DocType, doc.Id)
	exist, err := s.couchClient.DBExists(ctx, schema.GlobalDbName)
	if err != nil || !exist {
		return "", "", errors.Wrap(err, schema.GlobalDbName + " database check fails")
	}
	db := s.couchClient.DB(schema.GlobalDbName)

	rev, err := db.Put(ctx, doc.Id, doc)
	if err != nil {
		return "", "", errors.Wrap(err, schema.GlobalDbName + " database can't insert transaction id " + doc.Id)
	}
	return doc.Id, rev, nil
}

// AddTransactions bulk-adds transactions to CouchDB.
// It receives the []models.Transaction object and transform them into Transaction document objects and then
// insert them into the global CouchDB table.
func (s Store) AddTransactions(ctx context.Context, transactions []models.Transaction) (bool, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "transaction.AddTransactions")
	defer span.End()

	exist, err := s.couchClient.DBExists(ctx, schema.GlobalDbName)
	if err != nil || !exist {
		return false, errors.Wrap(err, schema.GlobalDbName + " database check fails")
	}
	db := s.couchClient.DB(schema.GlobalDbName)

	transactions_ := make([]interface{}, len(transactions))
	fmt.Println("Here are teh transactions")
	fmt.Printf("%v\n", transactions)

	// https://stackoverflow.com/questions/55755929/go-convert-interface-to-map
	// https://stackoverflow.com/questions/44094325/add-data-to-interface-in-struct
	for i := range transactions {
		doc := NewTransaction{
			ID: &transactions[i].Id,
			Transaction: transactions[i],
			DocType:     DocType,
			AssociatedAccounts: app.ExtractAccountAddrsFromTxn(transactions[i]),
			AssociatedApplications: app.ExtractApplicationIdsFromTxn(transactions[i]),
			AssociatedAssets: app.ExtractAssetIdsFromTxn(transactions[i]),
		}
		transactions_[i] = doc
		//fmt.Println("YYYYYYYYYY")
		//fmt.Printf("%v\n", transactions_[i])
		//v, _ := transactions_[i].(map[string]interface{})
		//fmt.Println("VVVVVVVV")
		//fmt.Printf("%v\n", v)
		//v["_id"] = transactions[i].Id
		//transactions_[i] = v
		//fmt.Println("looping")
		//fmt.Println(transactions_[i])
	}

	_, err = db.BulkDocs(ctx, transactions_)
	if err != nil {
		return false, errors.Wrap(err, "Can't bulk insert the transactions")
	}

	return true, nil
}

// GetTransaction retrieves a transaction record from CouchDB based upon the transaction ID given.
func (s Store) GetTransaction(ctx context.Context, transactionID string) (models.Transaction, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "transaction.GetTransaction")
	defer span.End()

	exist, err := s.couchClient.DBExists(ctx, schema.GlobalDbName)
	if err != nil || !exist {
		return models.Transaction{}, errors.Wrap(err, schema.GlobalDbName + " database check fails")
	}
	db := s.couchClient.DB(schema.GlobalDbName)

	docId := fmt.Sprintf("%s.%s", DocType, transactionID)
	row := db.Get(ctx, docId)
	if row == nil {
		return models.Transaction{}, errors.Wrap(err, schema.GlobalDbName + " get data empty")
	}

	var transaction Transaction
	fmt.Printf("%v\n", row)
	err = row.ScanDoc(&transaction)
	if err != nil {
		return models.Transaction{}, errors.Wrap(err, schema.GlobalDbName + "cannot unpack data from row")
	}

	return transaction.Transaction, nil
}

func (s Store) GetEarliestTransactionId(ctx context.Context) (string, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "transaction.GetEarliestTransactionId")
	//span.SetAttributes(attribute.String("query", q))
	defer span.End()

	exist, err := s.couchClient.DBExists(ctx, schema.GlobalDbName)
	if err != nil || !exist {
		return "", errors.Wrap(err, schema.GlobalDbName + " database check fails")
	}
	db := s.couchClient.DB(schema.GlobalDbName)

	rows, err := db.Query(ctx, schema.BlockDDoc, "_view/" + schema.TransactionViewByIdInLatest, kivik.Options{
		"include_docs": true,
		"descending": false,
		"limit": 1,
	})
	if err != nil {
		return "", errors.Wrap(err, "Fetch data error")
	}

	if rows.Err() != nil {
		return "", errors.Wrap(err, "rows error, Can't find anything")
	}

	rows.Next()
	var doc Transaction
	if err := rows.ScanDoc(&doc); err != nil {
		// No docs can be found
		return "", errors.Wrap(err, "Can't find anything")
	}

	return doc.Id, nil
}

func (s Store) GetLatestTransactionId(ctx context.Context) (string, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "transaction.GetLatestTransactionId")
	//span.SetAttributes(attribute.String("query", q))
	defer span.End()

	exist, err := s.couchClient.DBExists(ctx, schema.GlobalDbName)
	if err != nil || !exist {
		return "", errors.Wrap(err, schema.GlobalDbName + " database check fails")
	}
	db := s.couchClient.DB(schema.GlobalDbName)

	rows, err := db.Query(ctx, schema.BlockDDoc, "_view/" + schema.TransactionViewByIdInLatest, kivik.Options{
		"include_docs": true,
		"descending": true,
		"limit": 1,
	})
	if err != nil {
		return "", errors.Wrap(err, "Fetch data error")
	}

	if rows.Err() != nil {
		return "", errors.Wrap(err, "rows error, Can't find anything")
	}

	rows.Next()
	var doc Transaction
	if err := rows.ScanDoc(&doc); err != nil {
		// No docs can be found
		return "", errors.Wrap(err, "Can't find anything")
	}

	return doc.Id, nil
}

// https://stackoverflow.com/questions/11284383/couchdb-count-unique-document-field
// https://stackoverflow.com/questions/12944294/using-a-couchdb-view-can-i-count-groups-and-filter-by-key-range-at-the-same-tim
func (s Store) GetTransactionCountBtnKeys(ctx context.Context, startKey, endKey string) (int64, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "transaction.GetTransactionCountBtnKeys")
	//span.SetAttributes(attribute.String("query", q))
	defer span.End()

	exist, err := s.couchClient.DBExists(ctx, schema.GlobalDbName)
	if err != nil || !exist {
		return 0, errors.Wrap(err, schema.GlobalDbName + " database check fails")
	}
	db := s.couchClient.DB(schema.GlobalDbName)

	rows, err := db.Query(ctx, schema.BlockDDoc, "_view/" + schema.TransactionViewByIdInCount, kivik.Options{
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

func (s Store) GetTransactionsPagination(ctx context.Context, traceID string, log *zap.SugaredLogger, latestTransactionId string, order string, pageNo, limit int64) ([]Transaction, int64, int64, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "block.GetBlocksPagination")
	span.SetAttributes(attribute.String("latestTransactionId", latestTransactionId))
	span.SetAttributes(attribute.Int64("pageNo", pageNo))
	span.SetAttributes(attribute.Int64("limit", limit))
	defer span.End()

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

	rows, err := db.Query(ctx, schema.TransactionDDoc, "_view/" + schema.TransactionViewByIdInLatest, options)
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
