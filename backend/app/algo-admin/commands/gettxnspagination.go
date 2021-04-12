package commands

import (
	"context"
	"fmt"
	"github.com/kevguy/algosearch/backend/business/core/transaction"
	"github.com/kevguy/algosearch/backend/foundation/couchdb"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

func GetTransactionsPaginationCmd(traceID string, log *zap.SugaredLogger, couchCfg couchdb.Config, dbName, latestTxnID string, noOfItems, pageNo int64, order string) error {
	db, err := couchdb.Open(couchCfg)
	if err != nil {
		return errors.Wrap(err, ": connect to couchdb database")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer db.Close(ctx)
	defer cancel()

	transactionCore := transaction.NewCore(log, db, dbName)

	txns, numOfPages, numOfTxns, err := transactionCore.GetTransactionsPagination(ctx, latestTxnID, order, pageNo, noOfItems)
	fmt.Printf("Num of pages: %d", numOfPages)
	fmt.Printf("Num of txns: %d", numOfTxns)
	for _, item := range txns {
		fmt.Println(item.ID)
	}
	return nil
}
