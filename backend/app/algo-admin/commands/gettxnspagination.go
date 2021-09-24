package commands

import (
	"context"
	"fmt"
	"github.com/kevguy/algosearch/backend/business/data/store/transaction"
	"github.com/kevguy/algosearch/backend/foundation/couchdb"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

func GetTransactionsPaginationCmd(traceID string, log *zap.SugaredLogger, couchCfg couchdb.Config, latestTxnId string, noOfItems, pageNo int64, order string) error {
	db, err := couchdb.Open(couchCfg)
	if err != nil {
		return errors.Wrap(err, ": connect to couchdb database")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer db.Close(ctx)
	defer cancel()

	transactionStore := transaction.NewStore(log, db)

	txns, numOfPages, numOfTxns, err := transactionStore.GetTransactionsPagination(ctx, latestTxnId, order, pageNo, noOfItems)
	fmt.Printf("Num of pages: %d", numOfPages)
	fmt.Printf("Num of txns: %d", numOfTxns)
	for _, item := range txns {
		fmt.Println(item.ID)
	}
	return nil
}