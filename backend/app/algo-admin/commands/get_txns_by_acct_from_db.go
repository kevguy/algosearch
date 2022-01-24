package commands

import (
	"context"
	"fmt"
	"github.com/kevguy/algosearch/backend/business/core/transaction"
	"github.com/kevguy/algosearch/backend/foundation/couchdb"
	"go.uber.org/zap"
	"time"
)

// GetTransactionsByAcctFromDBCmd gets all transactions by an account from database.
func GetTransactionsByAcctFromDBCmd(log *zap.SugaredLogger, couchCfg couchdb.Config, dbName, acctID string) error {

	order := "asc" // "desc"

	db, err := couchdb.Open(couchCfg)
	if err != nil {
		return fmt.Errorf("connect to couchdb database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	transactionCore := transaction.NewCore(log, db, dbName)

	log.Infof("Getting List of Transactions")
	txns, err := transactionCore.GetTransactionsByAcct(ctx, acctID, order)
	if err != nil {
		return fmt.Errorf("getting transactions from account %s: %w", acctID, err)
	}

	log.Infof("Getting Earliest Transaction")
	earliestTxn, err := transactionCore.GetEarliestAcctTransaction(ctx, acctID)
	if err != nil {
		return fmt.Errorf("getting earliest transaction from account %s: %w", acctID, err)
	}
	//fmt.Printf("Earliest Transaction is %s\n", earliestTxn.ID)

	log.Infof("Getting Latest Transaction")
	latestTxn, err := transactionCore.GetLatestAcctTransaction(ctx, acctID)
	if err != nil {
		return fmt.Errorf("getting latest transaction from account %s: %w", acctID, err)
	}
	//fmt.Printf("Latest Transaction is %s\n", latestTxn.ID)

	log.Infof("Getting Transaction Count")
	count, err := transactionCore.GetTransactionCountByAcct(ctx, acctID, earliestTxn.ID, latestTxn.ID)
	if err != nil {
		return fmt.Errorf("getting transaction count from account %s: %w", acctID, err)
	}
	//fmt.Printf("Transaction Count is %s\n", count)

	fmt.Println("=====================================================")
	fmt.Printf("Account ID: %s\n", acctID)
	fmt.Printf("Number of Transactions: %d\n", count)
	fmt.Printf("Showing Transactions in %s order\n", order)
	fmt.Println("#\tID\t\t\t\t\t\t\tType\tRound Time (Epoch)\tRound Time (TimeStamp)")
	for idx, txn := range txns {
		fmt.Printf("%d\t%s\t%s\t%d\t\t%s\n", idx + 1, txn.ID, txn.Type, txn.RoundTime, time.Unix(int64(txn.RoundTime), 0).String())
	}
	fmt.Println("=====================================================")

	return nil
}
