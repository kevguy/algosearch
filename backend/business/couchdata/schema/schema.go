// Package schema contains the database schema, migrations and seeding data.
package schema

import (
	"context"
	_ "github.com/go-kivik/couchdb/v4" // The CouchDB driver
	kivik "github.com/go-kivik/kivik/v4"
	"github.com/kevguy/algosearch/backend/foundation/couchdb"
	"github.com/pkg/errors"
)

func createDB(ctx context.Context, db *kivik.Client, dbName string) error {
	// Create the "blocks" database
	exist, err := db.DBExists(ctx, dbName)
	if err != nil {
		return errors.Wrap(err, dbName + " database check fails")
	}
	if exist {
		err = db.DestroyDB(ctx, dbName)
		if err != nil {
			return errors.Wrap(err, dbName + " database deletion fails")
		}
	}
	err = db.CreateDB(ctx, dbName)
	if err != nil {
		return errors.Wrap(err, dbName + " database creation fails")
	}
	return nil
}

// InsertQueryViewForTransactionDB inserts the query view into the transactions databae.
// If it exists, continue. Else, create it.
func InsertQueryViewForTransactionDB(ctx context.Context, client *kivik.Client, dbName string) error {

	exist, err := client.DBExists(ctx, dbName)
	if err != nil || !exist {
		return errors.Wrap(err, dbName + " database check fails")
	}
	db := client.DB(dbName)
	//rows, err := db.Query(ctx, "_design/query", "_view/bytimestamp")
	//if err != nil {
	//	return errors.Wrap(err, dbName + " database and query by timestamp view failed to be queried")
	//}
	//if !rows.Next() {
	//	_, err = db.Put(context.TODO(), "_design/query", map[string]interface{}{
	//		"_id": "_design/query",
	//		"views": map[string]interface{}{
	//			"bytimestamp": map[string]interface{}{
	//				"map": "function(doc) { emit(doc.round); }",
	//			},
	//		},
	//	})
	//	if err != nil {
	//		return errors.Wrap(err, dbName + " database and query by timestamp view failed to be created")
	//	}
	//}
	_, err = db.Put(ctx, "_design/query", map[string]interface{}{
		"_id": "_design/query",
		"views": map[string]interface{}{
			"bytimestamp": map[string]interface{}{
				"map": "function(doc) { emit(doc.round); }",
			},
		},
	})
	if err != nil {
		return errors.Wrap(err, dbName + " database and query by timestamp view failed to be created")
	}
	return nil
}

// InsertLatestViewForBlocksDB inserts the latest view into the blocks database.
// If it exists, continue. Else, create it.
func InsertLatestViewForBlocksDB(ctx context.Context, client *kivik.Client, dbName string) error {

	exist, err := client.DBExists(ctx, dbName)
	if err != nil || !exist {
		return errors.Wrap(err, dbName + " database check fails")
	}
	db := client.DB(dbName)
	//rows, err := db.Query(ctx, "_design/latest", "_view/latest")
	//if err != nil {
	//	return errors.Wrap(err, dbName + " database and query by timestamp view failed to be queried")
	//}
	//if !rows.Next() {
	//	_, err = db.Put(context.TODO(), "_design/latest", map[string]interface{}{
	//		"_id": "_design/latest",
	//		"views": map[string]interface{}{
	//			"latest": map[string]interface{}{
	//				"map": "function(doc) { emit(doc.round); }",
	//			},
	//		},
	//	})
	//	if err != nil {
	//		return errors.Wrap(err, dbName + " database and query by timestamp view failed to be created")
	//	}
	//}
	_, err = db.Put(ctx, "_design/latest", map[string]interface{}{
		"_id": "_design/latest",
		"views": map[string]interface{}{
			"latest": map[string]interface{}{
				"map": "function(doc) { emit(doc.round); }",
			},
		},
	})
	if err != nil {
		return errors.Wrap(err, dbName + " database and query by timestamp view failed to be created")
	}
	return nil
}

// Migrate attempts to bring the schema for db up to date with the migrations
// defined in this package.
func Migrate(ctx context.Context, db *kivik.Client) error {
	if err := couchdb.StatusCheck(ctx, db); err != nil {
		return errors.Wrap(err, "status check database")
	}

	// Create the databases, delete the old ones if exist
	var dbList = []string{"blocks", "transactions", "addresses"}
	for _, dbName := range dbList {
		if err := createDB(ctx, db, dbName); err != nil {
			return errors.Wrap(err, dbName + " database creation fails")
		}
	}

	if err := InsertQueryViewForTransactionDB(ctx, db, "transactions"); err != nil {
		return errors.Wrap(err, "transactions database fails to create query view")
	}
	if err := InsertLatestViewForBlocksDB(ctx, db, "blocks"); err != nil {
		return errors.Wrap(err, "blocks database fails to create latest view")
	}

	return nil
}
