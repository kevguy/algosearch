// Package account provides the core business API of handling
// everything account related.
package account

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/go-kivik/kivik/v4"
	"github.com/kevguy/algosearch/backend/business/core/account/db"
	"go.uber.org/zap"
)

// Core manages the set of API's for block access.
type Core struct {
	store db.Store
}

// NewCore constructs a core for product api access.
func NewCore(log *zap.SugaredLogger, couchClient *kivik.Client, dbName string) Core {
	return Core{
		store: db.NewStore(log, couchClient, dbName),
	}
}

func (c Core) AddAccount(ctx context.Context, account models.Account) (string, string, error) {
	return c.store.AddAccount(ctx, account)
}

func (c Core) AddAccounts(ctx context.Context, accounts []models.Account) (bool, error) {
	return c.store.AddAccounts(ctx, accounts)
}

func (c Core) GetAccount(ctx context.Context, accountAddr string) (models.Account, error) {
	return c.store.GetAccount(ctx, accountAddr)
}

func (c Core) GetEarliestAccountID(ctx context.Context) (string, error) {
	return c.store.GetEarliestAccountID(ctx)
}

func (c Core) GetLatestAccountID(ctx context.Context) (string, error) {
	return c.store.GetLatestAccountID(ctx)
}

func (c Core) GetAccountCountBtnKeys(ctx context.Context, startKey, endKey string) (int64, error) {
	return c.store.GetAccountCountBtnKeys(ctx, startKey, endKey)
}

func (c Core) GetAccountsPagination(ctx context.Context, latestAccountID string, order string, pageNo, limit int64) ([]db.Account, int64, int64, error) {
	return c.store.GetAccountsPagination(ctx, latestAccountID, order, pageNo, limit)
}
