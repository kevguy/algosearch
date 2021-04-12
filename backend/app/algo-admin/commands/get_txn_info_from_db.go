package commands

import (
	"context"
	"fmt"
	"github.com/kevguy/algosearch/backend/business/core/transaction"
	"github.com/kevguy/algosearch/backend/foundation/couchdb"
	"go.uber.org/zap"
	"time"
)

// GetTransactionInfoFromDBCmd gets general transaction info from database, inc. earliest and latest transaction IDs
// and number of transactions
func GetTransactionInfoFromDBCmd(log *zap.SugaredLogger, couchCfg couchdb.Config, dbName string) error {

	// http://23.45.678.90:5984/algo_global/_design/txn/_view/txnByIdCount?
	// inclusive_end=true&
	// start_key=%5B%221560901701%22%2C%20%22MUCKTXOIUQ3RDM3UQOY56Z42MDBJSUGERRNPYKYMC7QYRRH3LHOA%22%5D&
	// end_key=%5B%221562014302%22%2C%20%22EWM2NMC33DJTRHCKQ3UYJXOWDK2UU4DPI2ZKYXX6ESEAUMFSE2PA%22%5D&
	// reduce=true&
	// group_level=0&
	// skip=0&
	// limit=101
	db, err := couchdb.Open(couchCfg)
	if err != nil {
		return fmt.Errorf("connect to couchdb database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	transactionCore := transaction.NewCore(log, db, dbName)

	earliestTxn, err := transactionCore.GetEarliestTransaction(ctx)
	if err != nil {
		return fmt.Errorf("getting the earliest Transaction ID: %w", err)
	}
	earliestID := earliestTxn.ID

	latestTxn, err := transactionCore.GetLatestTransaction(ctx)
	if err != nil {
		return fmt.Errorf("getting the latest Transaction ID: %w", err)
	}
	latestID := latestTxn.ID

	count, err := transactionCore.GetTransactionCountBtnKeys(ctx, earliestID, latestID)
	if err != nil {
		return fmt.Errorf("getting the transaction count: %w", err)
	}
	fmt.Println("=====================================================")
	fmt.Printf("Earliest Transaction ID: %s\n", earliestID)
	fmt.Printf("Latest Transaction ID: %s\n", latestID)
	fmt.Printf("Number of Transactions: %d\n", count)
	fmt.Println("=====================================================")

	return nil
}
