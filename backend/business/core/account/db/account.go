// Package db contains account related CRUD functionality.
package db

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/go-kivik/kivik/v4"
	"github.com/kevguy/algosearch/backend/business/data/schema"
	"github.com/kevguy/algosearch/backend/foundation/web"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

const (
	DocType = "acct"
)

// Store manages the set of API's for account access.
type Store struct {
	log *zap.SugaredLogger
	couchClient *kivik.Client
	dbName string
}

// NewStore constructs an account store for api access.
func NewStore(log *zap.SugaredLogger, couchClient *kivik.Client, dbName string) Store {
	return Store{
		log: log,
		couchClient: couchClient,
		dbName: dbName,
	}
}

// AddAccount adds an account to CouchDB.
// It receives the models.Account object and transform it into an Account document object and then
// insert it into the global CouchDB table.
func (s Store) AddAccount(ctx context.Context, account models.Account) (string, string, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "account.AddAccount")
	defer span.End()

	s.log.Infow("account.AddAccount", "traceid", web.GetTraceID(ctx))

	var doc = NewAccount{
		Account: account,
		DocType: DocType,
	}
	//docId := fmt.Sprintf("%s.%s", DocType, doc.Id)
	exist, err := s.couchClient.DBExists(ctx, s.dbName)
	if err != nil || !exist {
		return "", "", fmt.Errorf("%s database check fails: %w", s.dbName, err)
	}
	db := s.couchClient.DB(s.dbName)

	rev, err := db.Put(ctx, doc.Address, doc)
	if err != nil {
		return "", "", errors.Wrap(err, s.dbName+ " database can't insert account id " + doc.Address)
	}
	return doc.Address, rev, nil

}

// AddAccounts bulk-adds accounts to CouchDB.
// It receives the []models.Account object and transform them into Account document objects and then
// insert them into the global CouchDB table.
func (s Store) AddAccounts(ctx context.Context, accounts []models.Account) (bool, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "account.AddAccounts")
	defer span.End()

	s.log.Infow("account.AddAccounts", "traceid", web.GetTraceID(ctx))

	exist, err := s.couchClient.DBExists(ctx, s.dbName)
	if err != nil || !exist {
		return false, errors.Wrap(err, s.dbName+ " database check fails")
	}
	db := s.couchClient.DB(s.dbName)

	accounts_ := make([]interface{}, len(accounts))
	//fmt.Println("Here are the accounts")
	//fmt.Printf("%v\n", accounts)

	// https://stackoverflow.com/questions/55755929/go-convert-interface-to-map
	// https://stackoverflow.com/questions/44094325/add-data-to-interface-in-struct
	for i := range accounts {
		doc := NewAccount{
			ID:      &accounts[i].Address,
			Account: accounts[i],
			DocType: DocType,
		}
		accounts_[i] = doc
	}

	_, err = db.BulkDocs(ctx, accounts_)
	if err != nil {
		return false, errors.Wrap(err, "Can't bulk insert the accounts")
	}

	return true, nil
}

// GetAccount retrieves a account record from CouchDB based upon the account ID given.
func (s Store) GetAccount(ctx context.Context, accountAddr string) (models.Account, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "account.GetAccount")
	span.SetAttributes(attribute.String("accountAddr", accountAddr))
	defer span.End()

	s.log.Infow("account.GetAccount", "traceid", web.GetTraceID(ctx), "accountAddr", accountAddr)

	exist, err := s.couchClient.DBExists(ctx, s.dbName)
	if err != nil || !exist {
		return models.Account{}, errors.Wrap(err, s.dbName+ " database check fails")
	}
	db := s.couchClient.DB(s.dbName)

	docID := fmt.Sprintf("%s", accountAddr)
	row := db.Get(ctx, docID)
	if row == nil {
		return models.Account{}, errors.Wrap(err, s.dbName+ " get data empty")
	}

	var account Account
	err = row.ScanDoc(&account)
	if err != nil {
		return models.Account{}, errors.Wrap(err, s.dbName+ "cannot unpack data from row")
	}

	return account.Account, nil
}

func (s Store) GetEarliestAccountID(ctx context.Context) (string, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "account.GetEarliestAccountID")
	//span.SetAttributes(attribute.String("query", q))
	defer span.End()

	s.log.Infow("account.GetEarliestAccountID", "traceid", web.GetTraceID(ctx))

	exist, err := s.couchClient.DBExists(ctx, s.dbName)
	if err != nil || !exist {
		return "", errors.Wrap(err, s.dbName+ " database check fails")
	}
	db := s.couchClient.DB(s.dbName)

	rows, err := db.Query(ctx, schema.AccountDDoc, "_view/" +schema.AccountViewByIDInLatest, kivik.Options{
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
	var doc Account
	if err := rows.ScanDoc(&doc); err != nil {
		// No docs can be found
		return "", errors.Wrap(err, "Can't find anything")
	}

	return doc.Address, nil
}

func (s Store) GetLatestAccountID(ctx context.Context) (string, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "account.GetLatestAccountID")
	//span.SetAttributes(attribute.String("query", q))
	defer span.End()

	s.log.Infow("account.GetLatestAccountID", "traceid", web.GetTraceID(ctx))

	exist, err := s.couchClient.DBExists(ctx, s.dbName)
	if err != nil || !exist {
		return "", errors.Wrap(err, s.dbName+ " database check fails")
	}
	db := s.couchClient.DB(s.dbName)

	rows, err := db.Query(ctx, schema.AccountDDoc, "_view/" +schema.AccountViewByIDInLatest, kivik.Options{
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
	var doc Account
	if err := rows.ScanDoc(&doc); err != nil {
		// No docs can be found
		return "", errors.Wrap(err, "Can't find anything")
	}

	return doc.Address, nil
}

// GetAccountCountBtnKeys retrieves the number of keys between two keys
// References:
// 	https://stackoverflow.com/questions/11284383/couchdb-count-unique-document-field
// 	https://stackoverflow.com/questions/12944294/using-a-couchdb-view-can-i-count-groups-and-filter-by-key-range-at-the-same-tim
func (s Store) GetAccountCountBtnKeys(ctx context.Context, startKey, endKey string) (int64, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "account.GetAccountCountBtnKeys")
	span.SetAttributes(attribute.String("startKey", startKey))
	span.SetAttributes(attribute.String("endKey", endKey))
	defer span.End()

	s.log.Infow("account.GetAccountCountBtnKeys",
		"traceid", web.GetTraceID(ctx),
		"startKey", startKey,
		"endKey", endKey)

	if startKey == endKey {
		return 0, nil
	}

	exist, err := s.couchClient.DBExists(ctx, s.dbName)
	if err != nil || !exist {
		return 0, errors.Wrap(err, s.dbName+ " database check fails")
	}
	db := s.couchClient.DB(s.dbName)

	// curl 'http://admin:password@0.0.0.0:5984/algo_global/_design/acct/_view/acctByCount?startKey=2255PMXS65R54KKH5FQVV5UQZSAQCYL5U3OWQ2E5IZGOLK5XVTAVKNRPPQ&endKey=ZZYX3V6N74FGHGYLMSKJRVTRXT7GAZAQ47F4MOPX6S7FQRQU4FXZOLRQ2I'
	rows, err := db.Query(ctx, schema.AccountDDoc, "_view/" +schema.AccountViewByIDInCount, kivik.Options{
		"startKey": startKey,
		"endKey": endKey,
	})
	if err != nil {
		return 0, errors.Wrap(err, "Fetch data error")
	}

	var count int64
	for rows.Next() {
		if err := rows.ScanValue(&count); err != nil {
			return 0, errors.Wrap(err, "Can't find anything")
		}
	}

	return count, nil
}

func (s Store) GetAccountsPagination(ctx context.Context, latestAccountID string, order string, pageNo, limit int64) ([]Account, int64, int64, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "account.GetAccountsPagination")
	span.SetAttributes(attribute.String("latestAccountID", latestAccountID))
	span.SetAttributes(attribute.Int64("pageNo", pageNo))
	span.SetAttributes(attribute.Int64("limit", limit))
	defer span.End()

	s.log.Infow("account.GetAccountsPagination",
		"traceid", web.GetTraceID(ctx),
		"latestAccountID", latestAccountID,
		"pageNo", pageNo,
		"limit", limit)

	// Get the earliest account id
	earliestTxnID, err := s.GetEarliestAccountID(ctx)
	if err != nil {
		return nil, 0, 0, errors.Wrap(err, ": Get earliest synced account id")
	}

	numOfAccounts, err := s.GetAccountCountBtnKeys(ctx, earliestTxnID, latestAccountID)
	if err != nil {
		return nil, 0, 0, errors.Wrap(err, ": Get account count between keys")
	}

	// We can skip database check cuz GetEarliestAccountID already did it
	db := s.couchClient.DB(s.dbName)

	var numOfPages = numOfAccounts / limit
	if numOfAccounts % limit > 0 {
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
		options["start_key"] = latestAccountID

		// Use page number to calculate number of items to skip
		skip := (pageNo - 1) * limit
		options["skip"] = (pageNo - 1) * limit

		// Find the key to start reading and get the `page limit` number of records
		if (numOfAccounts - skip) > limit {
			options["limit"] = limit
		} else {
			options["limit"] = numOfAccounts - skip
		}
	} else {
		// Ascending order
		options["descending"] = false

		// Calculate the number of records to skip
		skip := (pageNo - 1) * limit
		options["skip"] = skip

		if (numOfAccounts - skip) > limit {
			options["limit"] =  numOfAccounts - skip
		} else {
			options["limit"] = limit
		}
	}

	rows, err := db.Query(ctx, schema.AccountDDoc, "_view/" +schema.AccountViewByIDInLatest, options)
	if err != nil {
		return nil, 0, 0, errors.Wrap(err, "Fetch data error")
	}

	var fetchedAccounts = []Account{}
	for rows.Next() {
		var account = Account{}
		if err := rows.ScanDoc(&account); err != nil {
			return nil, 0, 0, errors.Wrap(err, "unwrapping block")
		}
		fetchedAccounts = append(fetchedAccounts, account)
	}

	if rows.Err() != nil {
		return nil, 0, 0, errors.Wrap(err, "rows error, Can't find anything")
	}

	return fetchedAccounts, numOfPages, numOfAccounts, nil
}
