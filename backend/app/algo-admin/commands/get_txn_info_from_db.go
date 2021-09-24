package commands

import (
	"context"
	"fmt"
	"github.com/kevguy/algosearch/backend/business/data/store/transaction"
	"github.com/kevguy/algosearch/backend/foundation/couchdb"
	"go.uber.org/zap"
	"time"
)

// GetTransactionInfoFromDbCmd gets general transaction info from database, inc. earliest and latest transaction IDs
// and number of transactions
func GetTransactionInfoFromDbCmd(log *zap.SugaredLogger, couchCfg couchdb.Config) error {

	// http://89.39.110.254:5984/algo_global/_design/txn/_view/txnByIdCount?
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

	transactionStore := transaction.NewStore(log, db)

	earliestID, err := transactionStore.GetEarliestTransactionId(ctx)
	if err != nil {
		return fmt.Errorf("getting the earliest Transaction ID: %w", err)
	}

	latestID, err := transactionStore.GetLatestTransactionId(ctx)
	if err != nil {
		return fmt.Errorf("getting the latest Transaction ID: %w", err)
	}

	count, err := transactionStore.GetTransactionCountBtnKeys(ctx, earliestID, latestID)
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