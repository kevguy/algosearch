package block

import (
	"context"
	"github.com/go-kivik/kivik/v4"
	"github.com/kevguy/algosearch/backend/business/core/block/db"
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

func (c Core) AddBlock(ctx context.Context, block db.NewBlock) (string, string, error) {
	return c.store.AddBlock(ctx, block)
}

func (c Core) AddBlocks(ctx context.Context, blocks []db.Block) (bool, error) {
	return c.store.AddBlocks(ctx, blocks)
}

func (c Core) GetBlockByHash(ctx context.Context, blockHash string) (db.Block, error) {
	return c.store.GetBlockByHash(ctx, blockHash)
}

func (c Core) GetBlockByNum(ctx context.Context, blockNum uint64) (db.Block, error) {
	return c.store.GetBlockByNum(ctx, blockNum)
}

func (c Core) GetEarliestSyncedRoundNumber(ctx context.Context) (uint64, error) {
	return c.store.GetEarliestSyncedRoundNumber(ctx)
}

func (c Core) GetLastSyncedRoundNumber(ctx context.Context) (uint64, error) {
	return c.store.GetLastSyncedRoundNumber(ctx)
}

func (c Core) GetLatestBlock(ctx context.Context) (db.Block, error) {
	return c.store.GetLatestBlock(ctx)
}

func (c Core) GetBlocksPagination(ctx context.Context, latestBlockNum int64, order string, pageNo int64, limit int64) ([]db.Block, int64, int64, error) {
	return c.store.GetBlocksPagination(ctx, latestBlockNum, order, pageNo, limit)
}

func (c Core) GetNumOfBlocks(ctx context.Context) (int64, error) {
	return c.store.GetNumOfBlocks(ctx)
}