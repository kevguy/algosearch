package commands

import (
	"context"
	"fmt"
	"github.com/kevguy/algosearch/backend/business/data/store/transaction"
	"github.com/kevguy/algosearch/backend/foundation/couchdb"
	"go.uber.org/zap"
	"time"
)

func GetTransactionsByAcctFromDbWithPaginationCmd(
	log *zap.SugaredLogger,
	couchCfg couchdb.Config,
	acctID string,
	pageNo, limit int64, order string) error {

	db, err := couchdb.Open(couchCfg)
	if err != nil {
		return fmt.Errorf("connect to couchdb database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	transactionStore := transaction.NewStore(log, db)

	txns, noOfPages, noOfTxns, err := transactionStore.GetTransactionsByAcctPagination(ctx, acctID, order, pageNo, limit)
	if err != nil {
		return fmt.Errorf("getting transactions by pagination from db: %w", err)
	}

	fmt.Println("=====================================================")
	fmt.Printf("Account ID: %s\n", acctID)
	fmt.Printf("Limit: %d\n", limit)
	fmt.Printf("Order: %s\n", order)
	fmt.Printf("No. of Pages: %d\n", noOfPages)
	fmt.Printf("No. of Transactions: %d\n", noOfTxns)
	//fmt.Printf("Number of Transactions: %d\n", count)
	fmt.Printf("Showing Transactions in %s order\n", order)
	fmt.Println("#\tID\t\t\t\t\t\t\tType\tRound Time (Epoch)\tRound Time (TimeStamp)")
	for idx, txn := range txns {
		fmt.Printf("%d\t%s\t%s\t%d\t\t%s\n", idx + 1, txn.ID, txn.Type, txn.RoundTime, time.Unix(int64(txn.RoundTime), 0).String())
	}
	fmt.Println("=====================================================")


	//fmt.Println("=====================================================")
	//fmt.Printf("No. of Transaction ID: %d\n", noOfTxns)
	//fmt.Printf("Limit: %d\n", limit)
	//fmt.Printf("Order: %s\n", order)
	//fmt.Printf("No. of Pages: %d\n", noOfPages)
	//fmt.Printf("Transactions for Page %d:\n", pageNo)
	//for idx, txn := range txns {
	//	fmt.Printf("\t%d - %s\n", idx + 1, txn.ID)
	//}
	//fmt.Println("=====================================================")

	return nil
}
