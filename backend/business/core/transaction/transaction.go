package transaction

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/go-kivik/kivik/v4"
	"github.com/kevguy/algosearch/backend/business/core/transaction/db"
	"go.uber.org/zap"
)

// Core manages the set of API's for block access.
type Core struct {
	store db.Store
}

// NewCore constructs a core for product api access.
func NewCore(log *zap.SugaredLogger, couchClient *kivik.Client) Core {
	return Core{
		store: db.NewStore(log, couchClient),
	}
}

func (c Core) AddTransaction(ctx context.Context, transaction models.Transaction) (string, string, error) {
	return c.store.AddTransaction(ctx, transaction)
}

func (c Core) AddTransactions(ctx context.Context, transactions []models.Transaction) (bool, error) {
	return c.store.AddTransactions(ctx, transactions)
}

func (c Core) GetTransaction(ctx context.Context, transactionID string) (models.Transaction, error) {
	return c.store.GetTransaction(ctx, transactionID)
}

func (c Core) GetTransactionCountBtnKeys(ctx context.Context, startKey, endKey string) (int64, error) {
	return c.store.GetTransactionCountBtnKeys(ctx, startKey, endKey)
}

func (c Core) GetEarliestTransaction(ctx context.Context) (db.Transaction, error) {
	return c.store.GetEarliestTransaction(ctx)
}

func (c Core) GetLatestTransaction(ctx context.Context) (db.Transaction, error) {
	return c.store.GetLatestTransaction(ctx)
}

func (c Core) GetTransactionsPagination(ctx context.Context, startTransactionId, order string, pageNo, limit int64) ([]db.Transaction, int64, int64, error) {
	return c.store.GetTransactionsPagination(ctx, startTransactionId, order, pageNo, limit)
}

func (c Core) GetEarliestAcctTransaction(ctx context.Context, acctID string) (db.Transaction, error) {
	return c.store.GetEarliestAcctTransaction(ctx, acctID)
}

func (c Core) GetLatestAcctTransaction(ctx context.Context, acctID string) (db.Transaction, error) {
	return c.store.GetLatestAcctTransaction(ctx, acctID)
}

func (c Core) GetTransactionCountByAcct(ctx context.Context, acctID, startKey, endKey string) (int64, error) {
	return c.store.GetTransactionCountByAcct(ctx, acctID, startKey, endKey)
}

func (c Core) GetTransactionsByAcctPagination(ctx context.Context, acctID, order string, pageNo, limit int64) ([]db.Transaction, int64, int64, error) {
	return c.store.GetTransactionsByAcctPagination(ctx, acctID, order, pageNo, limit)
}

func (c Core) GetTransactionsByAcct(ctx context.Context, acctID string, order string) ([]db.Transaction, error) {
	return c.store.GetTransactionsByAcct(ctx, acctID, order)
}

func (c Core) GetEarliestAppTransaction(ctx context.Context, appID string) (db.Transaction, error) {
	return c.store.GetEarliestAppTransaction(ctx, appID)
}

func (c Core) GetLatestAppTransaction(ctx context.Context, appID string) (db.Transaction, error) {
	return c.store.GetLatestAppTransaction(ctx, appID)
}

func (c Core) GetTransactionsByApp(ctx context.Context, appID string, order string) ([]db.Transaction, error) {
	return c.store.GetTransactionsByApp(ctx, appID, order)
}

func (c Core) GetEarliestAssetTransaction(ctx context.Context, assetID string) (db.Transaction, error) {
	return c.store.GetEarliestAssetTransaction(ctx, assetID)
}

func (c Core) GetLatestAssetTransaction(ctx context.Context, assetID string) (db.Transaction, error) {
	return c.store.GetLatestAssetTransaction(ctx, assetID)
}

func (c Core) GetTransactionsByAsset(ctx context.Context, assetID string, order string) ([]db.Transaction, error) {
	return c.store.GetTransactionsByAsset(ctx, assetID, order)
}
