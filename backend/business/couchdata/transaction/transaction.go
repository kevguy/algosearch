// Package transaction contains transaction related CRUD functionality.
package transaction

import (
	"context"
	"fmt"
	"github.com/go-kivik/kivik/v4"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

const (
	TRANSACTIONS = "transactions"
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
func (s Store) AddTransaction(ctx context.Context, transaction Transaction) (string, string, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "transaction.AddTransaction")
	defer span.End()

	exist, err := s.couchClient.DBExists(ctx, TRANSACTIONS)
	if err != nil || !exist {
		return "", "", errors.Wrap(err, TRANSACTIONS+ " database check fails")
	}
	db := s.couchClient.DB(TRANSACTIONS)

	rev, err := db.Put(ctx, transaction.SuggestedParams.GenesisID, transaction)
	if err != nil {
		return "", "", errors.Wrap(err, TRANSACTIONS+ " database can't insert transaction id " + transaction.SuggestedParams.GenesisID)
	}
	return transaction.SuggestedParams.GenesisID, rev, nil
}

// GetTransaction adds a retrieves a transaction record from CouchDB.
func (s Store) GetTransaction(ctx context.Context, genesisHash string) (TransactionRecord, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "transaction.GetTransaction")
	defer span.End()

	exist, err := s.couchClient.DBExists(ctx, TRANSACTIONS)
	if err != nil || !exist {
		return TransactionRecord{}, errors.Wrap(err, TRANSACTIONS+ " database check fails")
	}
	db := s.couchClient.DB(TRANSACTIONS)

	row := db.Get(ctx, genesisHash)
	if row == nil {
		return TransactionRecord{}, errors.Wrap(err, TRANSACTIONS+ " get data empty")
	}

	var transaction TransactionRecord
	fmt.Printf("%v\n", row)
	err = row.ScanDoc(&transaction)
	if err != nil {
		return TransactionRecord{}, errors.Wrap(err, TRANSACTIONS+ "cannot unpack data from row")
	}

	return transaction, nil
}
