package commands

import (
	"context"
	"fmt"
	"github.com/kevguy/algosearch/backend/business/data/store/transaction"
	"github.com/kevguy/algosearch/backend/foundation/couchdb"
	"go.uber.org/zap"
	"time"
)

func GetTransactionInfoFromDbCmd(log *zap.SugaredLogger, couchCfg couchdb.Config) error {

	db, err := couchdb.Open(couchCfg)
	if err != nil {
		return fmt.Errorf("connect to couchdb database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	transactionStore := transaction.NewStore(log, db)

	earliestID, err := transactionStore.GetEarliestTransactionId(ctx)
	if err != nil {
		return fmt.Errorf("getting the earliest Transaction ID: %w", err)
	}
	fmt.Printf("Earliest Transaction ID: %s\n", earliestID)

	latestID, err := transactionStore.GetLatestTransactionId(ctx)
	if err != nil {
		return fmt.Errorf("getting the latest Transaction ID: %w", err)
	}
	fmt.Printf("Latest Transaction ID: %s\n", latestID)

	count, err := transactionStore.GetTransactionCountBtnKeys(ctx, latestID, earliestID)
	if err != nil {
		return fmt.Errorf("getting the transaction count: %w", err)
	}
	fmt.Printf("Number of Transactions: %d\n", count)

	return nil
}
