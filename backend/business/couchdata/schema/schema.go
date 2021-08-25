// Package schema contains the database schema, migrations and seeding data.
package schema

import (
	"context"
	_ "github.com/go-kivik/couchdb/v4" // The CouchDB driver
	kivik "github.com/go-kivik/kivik/v4"
	"github.com/kevguy/algosearch/backend/foundation/couchdb"
	"github.com/pkg/errors"
)

const (
	GlobalDbName = "algo_global"

	BlockDDoc                = "_design/block"
	BlockViewByRoundInLatest = "blockByLatest"

	TransactionDDoc = "_design/txn"
	TransactionLatestView = "txnByLatest"

	AccountDDoc = "_design/acct"
	AccountLatestView = "byLatest"

	AssetDDoc = "_design/asset"
	AssetLatestView = "byLatest"

	ApplicationDDoc = "_design/app"
	ApplicationLatestView = "byLatest"
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

// InsertBlockViewsForGlobalDB creates a the latest view for the block design document. It stores
// block data.
func InsertBlockViewsForGlobalDB(ctx context.Context, client *kivik.Client, dbName string) error {
	// Check if DB exists
	exist, err := client.DBExists(ctx, dbName)
	if err != nil || !exist {
		return errors.Wrap(err, dbName + " database check fails")
	}
	db := client.DB(dbName)

	_, err = db.Query(ctx, BlockDDoc, "_view/" + BlockViewByRoundInLatest)
	//if err != nil {
	//	return errors.Wrap(err, dbName + " database and query by timestamp view failed to be queried")
	//}
	//if !rows.Next() {
		_, err = db.Put(context.TODO(), BlockDDoc, map[string]interface{}{
			"_id": BlockDDoc,
			"views": map[string]interface{}{
				// https://docs.couchdb.org/en/main/ddocs/views/joins.html
				BlockViewByRoundInLatest: map[string]interface{}{
					"map": "function(doc) { if (doc.doc_type === 'block')  { emit(doc.round, {_id: doc.BlockHash}); } }",
				},
			},
		})
		if err != nil {
			return errors.Wrap(err, dbName + " database and query by timestamp view failed to be created")
		}
	//}
	return nil
}

// InsertTransactionViewsForGlobalDB creates a the latest view for the transaction design document. It stores
// transaction data.
func InsertTransactionViewsForGlobalDB(ctx context.Context, client *kivik.Client, dbName string) error {
	// Check if DB exists
	exist, err := client.DBExists(ctx, dbName)
	if err != nil || !exist {
		return errors.Wrap(err, dbName + " database check fails")
	}
	db := client.DB(dbName)

	_, err = db.Query(ctx, TransactionDDoc, "_view/" + TransactionLatestView)
	//if err != nil {
	//	return errors.Wrap(err, dbName + " database and query by timestamp view failed to be queried")
	//}
	//if !rows.Next() {
		_, err = db.Put(context.TODO(), "_design/" + TransactionLatestView, map[string]interface{}{
			"_id": TransactionDDoc,
			"views": map[string]interface{}{
				TransactionLatestView: map[string]interface{}{
					"map": `function(doc) { 
						if (doc.type === 'txn') {
							emit(doc.round, doc);
						}
					}`,
				},
			},
		})
		if err != nil {
			return errors.Wrap(err, dbName + " database and query by timestamp view failed to be created")
		}
	//}
	return nil
}


// https://stackoverflow.com/questions/5422622/couchdb-views-tied-between-two-databases
// https://stackoverflow.com/questions/6380045/couchdb-join-two-documents
// https://stackoverflow.com/questions/24264898/combine-multiple-documents-in-a-couchdb-view
// https://stackoverflow.com/questions/23358813/join-two-different-documents-in-couchdb-using-futon-map-program
// https://stackoverflow.com/questions/37487069/multiple-joins-in-couchdb
// https://stackoverflow.com/questions/50422224/get-cloudant-couchdb-database-documents-by-passing-a-list-of-keys
// https://www.cmlenz.net/archives/2007/10/couchdb-joins

// InsertBlockViewForGlobalDB inserts a view that is specifically for blocks. It stores the block data in
// using the round number as ID, in the format of `block_ROUND_NO`.
func InsertBlockViewForGlobalDB(ctx context.Context, client *kivik.Client, dbName string) error {
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
	//var dbList = []string{"blocks", "transactions", "addresses"}
	var dbList = []string{GlobalDbName}
	for _, dbName := range dbList {
		if err := createDB(ctx, db, dbName); err != nil {
			return errors.Wrap(err, dbName + " database creation fails")
		}
	}

	// Block views
	if err := InsertBlockViewsForGlobalDB(ctx, db, GlobalDbName); err != nil {
		return errors.Wrap(err, "database fails to create view(s) for blocks")
	}

	// Transaction views
	if err := InsertTransactionViewsForGlobalDB(ctx, db, GlobalDbName); err != nil {
		return errors.Wrap(err, "database fails to create view(s) for transactions")
	}

	// To-Be-Deleted
	//if err := InsertQueryViewForTransactionDB(ctx, db, "transactions"); err != nil {
	//	return errors.Wrap(err, "transactions database fails to create query view")
	//}
	//if err := InsertLatestViewForBlocksDB(ctx, db, "blocks"); err != nil {
	//	return errors.Wrap(err, "blocks database fails to create latest view")
	//}

	return nil
}
