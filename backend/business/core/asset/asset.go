// Package asset provides the core business API of handling
// everything asset related.
package asset

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/go-kivik/kivik/v4"
	"github.com/kevguy/algosearch/backend/business/core/asset/db"
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

func (c Core) AddAsset(ctx context.Context, asset models.Asset) (string, string, error) {
	return c.store.AddAsset(ctx, asset)
}

func (c Core) AddAssets(ctx context.Context, assets []models.Asset) (bool, error) {
	return c.store.AddAssets(ctx, assets)
}

func (c Core) GetAsset(ctx context.Context, assetID string) (models.Asset, error) {
	return c.store.GetAsset(ctx, assetID)
}

func (c Core) GetEarliestAssetID(ctx context.Context) (string, error) {
	return c.store.GetEarliestAssetID(ctx)
}

func (c Core) GetLatestAssetID(ctx context.Context) (string, error) {
	return c.store.GetLatestAssetID(ctx)
}

func (c Core) GetAssetCountBtnKeys(ctx context.Context, startKey, endKey string) (int64, error) {
	return c.store.GetAssetCountBtnKeys(ctx, startKey, endKey)
}

func (c Core) GetAssetsPagination(ctx context.Context, latestAssetID string, order string, pageNo, limit int64) ([]db.Asset, int64, int64, error) {
	return c.store.GetAssetsPagination(ctx, latestAssetID, order, pageNo, limit)
}
