package commands

import (
	"context"
	"fmt"
	"github.com/kevguy/algosearch/backend/business/data/store/transaction"
	"github.com/kevguy/algosearch/backend/foundation/couchdb"
	"go.uber.org/zap"
	"time"
)

// GetTransactionsByAcctFromDbCmd gets all transactions by an account from database.
func GetTransactionsByAcctFromDbCmd(log *zap.SugaredLogger, couchCfg couchdb.Config, acctID string) error {

	// http://89.39.110.254:5984/algo_global/_design/txn/_view/txnByAcct?
	// include_docs=true&
	// inclusive_end=true&
	// start_key=%5B%222255PMXS65R54KKH5FQVV5UQZSAQCYL5U3OWQ2E5IZGOLK5XVTAVKNRPPQ%22%2C%20%221%22%5D&
	// end_key=%5B%222255PMXS65R54KKH5FQVV5UQZSAQCYL5U3OWQ2E5IZGOLK5XVTAVKNRPPQ%22%2C%20%222%22%5D&skip=0&
	// limit=101&
	// reduce=false

	db, err := couchdb.Open(couchCfg)
	if err != nil {
		return fmt.Errorf("connect to couchdb database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	transactionStore := transaction.NewStore(log, db)

	txns, err := transactionStore.GetTransactionsByAcct(ctx, acctID, "desc")
	if err != nil {
		return fmt.Errorf("getting transactions from account %s: %w", acctID, err)
	}

	fmt.Println("=====================================================")
	fmt.Println("#\tID\t\tRound Time")
	for idx, txn := range txns {
		fmt.Printf("%d\t%s\t%d\t%s\n", idx + 1, txn.ID, txn.RoundTime, time.Unix(int64(txn.RoundTime), 0).String())
	}
	fmt.Println("=====================================================")

	return nil
}
