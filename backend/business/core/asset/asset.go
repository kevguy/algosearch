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
func NewCore(log *zap.SugaredLogger, couchClient *kivik.Client) Core {
	return Core{
		store: db.NewStore(log, couchClient),
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

func (c Core) GetEarliestAssetId(ctx context.Context) (string, error) {
	return c.store.GetEarliestAssetId(ctx)
}

func (c Core) GetLatestAssetId(ctx context.Context) (string, error) {
	return c.store.GetLatestAssetId(ctx)
}

func (c Core) GetAssetCountBtnKeys(ctx context.Context, startKey, endKey string) (int64, error) {
	return c.store.GetAssetCountBtnKeys(ctx, startKey, endKey)
}

func (c Core) GetAssetsPagination(ctx context.Context, latestAssetId string, order string, pageNo, limit int64) ([]db.Asset, int64, int64, error) {
	return c.store.GetAssetsPagination(ctx, latestAssetId, order, pageNo, limit)
}
