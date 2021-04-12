package commands

import (
	"context"
	"fmt"
	"github.com/kevguy/algosearch/backend/business/core/transaction"
	"github.com/kevguy/algosearch/backend/foundation/couchdb"
	"go.uber.org/zap"
	"time"
)

func GetTransactionsFromDbWithPaginationCmd(
	log *zap.SugaredLogger,
	couchCfg couchdb.Config,
	pageNo, limit int64, order string) error {

	db, err := couchdb.Open(couchCfg)
	if err != nil {
		return fmt.Errorf("connect to couchdb database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	transactionCore := transaction.NewCore(log, db)

	txns, noOfPages, noOfTxns, err := transactionCore.GetTransactionsPagination(ctx, "", order, pageNo, limit)
	if err != nil {
		return fmt.Errorf("getting transactions by pagination from db: %w", err)
	}

	fmt.Println("=====================================================")
	fmt.Printf("No. of Transaction ID: %d\n", noOfTxns)
	fmt.Printf("Limit: %d\n", limit)
	fmt.Printf("Order: %s\n", order)
	fmt.Printf("No. of Pages: %d\n", noOfPages)
	fmt.Printf("Transactions for Page %d:\n", pageNo)
	for idx, txn := range txns {
		fmt.Printf("\t%d - %s\n", idx + 1, txn.ID)
	}
	fmt.Println("=====================================================")

	return nil
}
