package application

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/go-kivik/kivik/v4"
	"github.com/kevguy/algosearch/backend/business/core/application/db"
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

func (c Core) AddApplication(ctx context.Context, application models.Application) (string, string, error) {
	return c.store.AddApplication(ctx, application)
}

func (c Core) AddApplications(ctx context.Context, applications []models.Application) (bool, error) {
	return c.store.AddApplications(ctx, applications)
}

func (c Core) GetApplication(ctx context.Context, applicationID string) (models.Application, error) {
	return c.store.GetApplication(ctx, applicationID)
}

func (c Core) GetEarliestApplicationId(ctx context.Context) (string, error) {
	return c.store.GetEarliestApplicationId(ctx)
}

func (c Core) GetLatestApplicationId(ctx context.Context) (string, error) {
	return c.store.GetLatestApplicationId(ctx)
}

func (c Core) GetApplicationCountBtnKeys(ctx context.Context, startKey, endKey string) (int64, error) {
	return c.store.GetApplicationCountBtnKeys(ctx, startKey, endKey)
}

func (c Core) GetApplicationsPagination(ctx context.Context, latestApplicationId string, order string, pageNo, limit int64) ([]db.Application, int64, int64, error) {
	return c.store.GetApplicationsPagination(ctx, latestApplicationId, order, pageNo, limit)
}
