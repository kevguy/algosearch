package block

import (
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
