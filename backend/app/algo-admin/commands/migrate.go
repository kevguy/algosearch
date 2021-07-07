package commands

import (
	"context"
	"fmt"
	"github.com/kevguy/algosearch/backend/business/couchdata/schema"
	"github.com/kevguy/algosearch/backend/foundation/couchdb"
	"github.com/pkg/errors"
	"time"
)

// ErrHelp provides context that help was given.
var ErrHelp = errors.New("provided help")

// Migrate creates the schema in the database.
func Migrate(cfg couchdb.Config) error {
	db, err := couchdb.Open(cfg)
	if err != nil {
		return errors.Wrap(err, "connect to couchdb database")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer db.Close(ctx)
	defer cancel()

	if err := schema.Migrate(ctx, db); err != nil {
		return errors.Wrap(err, "migrate couchdb database")
	}

	fmt.Println("migrations complete")
	return nil
}
