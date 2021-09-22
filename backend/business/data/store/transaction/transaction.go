// Package transaction contains transaction related CRUD functionality.
package transaction

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/go-kivik/kivik/v4"
	app "github.com/kevguy/algosearch/backend/business/algod"
	"github.com/kevguy/algosearch/backend/business/data/schema"
	"github.com/kevguy/algosearch/backend/foundation/web"
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

	s.log.Infow("transaction.AddTransaction", "traceid", web.GetTraceID(ctx))

	var doc = NewTransaction{
		Transaction:            transaction,
		DocType:                DocType,
		AssociatedAccounts:     app.ExtractAccountAddrsFromTxn(transaction),
		AssociatedApplications: app.ExtractApplicationIdsFromTxn(transaction),
		AssociatedAssets:       app.ExtractAssetIdsFromTxn(transaction),
	}
	//docId := fmt.Sprintf("%s.%s", DocType, doc.Id)
	exist, err := s.couchClient.DBExists(ctx, schema.GlobalDbName)
	if err != nil || !exist {
		return "", "", errors.Wrap(err, schema.GlobalDbName+ " database check fails")
	}
	db := s.couchClient.DB(schema.GlobalDbName)

	rev, err := db.Put(ctx, doc.Id, doc)
	if err != nil {
		return "", "", errors.Wrap(err, schema.GlobalDbName+ " database can't insert transaction id " + doc.Id)
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

	s.log.Infow("transaction.AddTransactions", "traceid", web.GetTraceID(ctx))

	exist, err := s.couchClient.DBExists(ctx, schema.GlobalDbName)
	if err != nil || !exist {
		return false, errors.Wrap(err, schema.GlobalDbName+ " database check fails")
	}
	db := s.couchClient.DB(schema.GlobalDbName)

	transactions_ := make([]interface{}, len(transactions))
	//fmt.Println("Here are teh transactions")
	//fmt.Printf("%v\n", transactions)

	// https://stackoverflow.com/questions/55755929/go-convert-interface-to-map
	// https://stackoverflow.com/questions/44094325/add-data-to-interface-in-struct
	for i := range transactions {
		doc := NewTransaction{
			ID:                     &transactions[i].Id,
			Transaction:            transactions[i],
			DocType:                DocType,
			AssociatedAccounts:     app.ExtractAccountAddrsFromTxn(transactions[i]),
			AssociatedApplications: app.ExtractApplicationIdsFromTxn(transactions[i]),
			AssociatedAssets:       app.ExtractAssetIdsFromTxn(transactions[i]),
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
	span.SetAttributes(attribute.String("transactionID", transactionID))
	defer span.End()

	s.log.Infow("transaction.GetTransaction", "traceid", web.GetTraceID(ctx), "transactionID", transactionID)

	exist, err := s.couchClient.DBExists(ctx, schema.GlobalDbName)
	if err != nil || !exist {
		return models.Transaction{}, errors.Wrap(err, schema.GlobalDbName+ " database check fails")
	}
	db := s.couchClient.DB(schema.GlobalDbName)

	docId := fmt.Sprintf("%s.%s", DocType, transactionID)
	row := db.Get(ctx, docId)
	if row == nil {
		return models.Transaction{}, errors.Wrap(err, schema.GlobalDbName+ " get data empty")
	}

	var transaction Transaction
	fmt.Printf("%v\n", row)
	err = row.ScanDoc(&transaction)
	if err != nil {
		return models.Transaction{}, errors.Wrap(err, schema.GlobalDbName+ "cannot unpack data from row")
	}

	return transaction.Transaction, nil
}

func (s Store) GetEarliestTransactionId(ctx context.Context) (string, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "transaction.GetEarliestTransactionId")
	//span.SetAttributes(attribute.String("query", q))
	defer span.End()

	s.log.Infow("transaction.GetEarliestTransactionId", "traceid", web.GetTraceID(ctx))

	exist, err := s.couchClient.DBExists(ctx, schema.GlobalDbName)
	if err != nil || !exist {
		return "", errors.Wrap(err, schema.GlobalDbName+ " database check fails")
	}
	db := s.couchClient.DB(schema.GlobalDbName)

	rows, err := db.Query(ctx, schema.BlockDDoc, "_view/" +schema.TransactionViewInLatest, kivik.Options{
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

	s.log.Infow("transaction.GetLatestTransactionId", "traceid", web.GetTraceID(ctx))

	exist, err := s.couchClient.DBExists(ctx, schema.GlobalDbName)
	if err != nil || !exist {
		return "", errors.Wrap(err, schema.GlobalDbName+ " database check fails")
	}
	db := s.couchClient.DB(schema.GlobalDbName)

	rows, err := db.Query(ctx, schema.BlockDDoc, "_view/" +schema.TransactionViewInLatest, kivik.Options{
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
