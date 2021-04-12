package commands

import (
	"context"
	"fmt"
	"github.com/kevguy/algosearch/backend/business/data/schema"
	"github.com/kevguy/algosearch/backend/foundation/couchdb"
	"github.com/pkg/errors"
	"time"
)

// ErrHelp provides context that help was given.
var ErrHelp = errors.New("provided help")

// Migrate creates the schema in the database.
func Migrate(cfg couchdb.Config, dbName string) error {
	db, err := couchdb.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer db.Close(ctx)
	defer cancel()

	if err := schema.Migrate(ctx, db, dbName); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	fmt.Println("migrations complete")
	return nil
}
