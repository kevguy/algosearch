// Package schema contains the database schema, migrations and seeding data.
package schema

import (
	"context"
	"fmt"
	_ "github.com/go-kivik/couchdb/v4" // The CouchDB driver
	kivik "github.com/go-kivik/kivik/v4"
	"github.com/kevguy/algosearch/backend/foundation/couchdb"
	"github.com/pkg/errors"
)

const (
	GlobalDbName = "algo_global"

	BlockDDoc          = "_design/block"
	BlockViewByRoundNo = "blockByRoundNo"
	BlockViewByRoundCount = "blockByRoundNoCount"

	TransactionDDoc            = "_design/txn"
	TransactionViewInLatest    = "txnInLatest"
	TransactionViewById      = "txnById"
	TransactionViewByIdCount = "txnByIdCount"
	TransactionViewByAccount = "txnByAcct"
	TransactionViewByAccountCount		= "txnByAcctCount"
	TransactionViewByAsset				= "txnByAsset"
	TransactionViewByAssetCount			= "txnByAssetCount"
	TransactionViewByApplication		= "txnByApp"
	TransactionViewByApplicationCount	= "txnByAppCount"

	AccountDDoc             = "_design/acct"
	AccountViewByIdInLatest = "acctByLatest"
	AccountViewByIdInCount  = "acctByCount"

	AssetDDoc             = "_design/asset"
	AssetViewByIdInLatest = "assetByLatest"
	AssetViewByIdInCount  = "assetByCount"

	ApplicationDDoc             = "_design/app"
	ApplicationViewByIdInLatest = "appByLatest"
	ApplicationViewByIdInCount  = "appByCount"
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

	rows, err := db.Query(ctx, BlockDDoc, "_view/" +BlockViewByRoundNo)
	//if err != nil {
	//	return errors.Wrap(err, dbName + " database and query by timestamp view failed to be queried")
	//}
	if rows == nil || !rows.Next() {
		_, err = db.Put(context.TODO(), BlockDDoc, map[string]interface{}{
			"_id": BlockDDoc,
			"views": map[string]interface{}{
				// https://docs.couchdb.org/en/main/ddocs/views/joins.html
				BlockViewByRoundNo: map[string]interface{}{
					"map": `function(doc) { 
						if (doc.doc_type === 'block')  {
							// emit(doc.round, {_id: doc.BlockHash});
							emit(doc.round, null);
						}
					}`,
				},
				BlockViewByRoundCount: map[string]interface{}{
					"map": `function(doc) {
						if (doc.doc_type === 'block') {
							emit(doc.round, 1);
						}
					}`,
					"reduce": "_sum",
				},
			},
		})
		if err != nil && err.Error() != "Conflict: Document update conflict." {
			return fmt.Errorf("%s database and query by timestamp block view failed to be created: %w", dbName, err)
		}
	}
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

	rows, err := db.Query(ctx, TransactionDDoc, "_view/" +TransactionViewInLatest)
	//if err != nil {
	//	return errors.Wrap(err, dbName + " database and query by timestamp view failed to be queried")
	//}
	//if !rows.Next() {
	if rows == nil || !rows.Next() {
		_, err = db.Put(context.TODO(), TransactionDDoc, map[string]interface{}{
			"_id": TransactionDDoc,
			"views": map[string]interface{}{
				TransactionViewInLatest: map[string]interface{}{
					"map": `function(doc) { 
						if (doc.doc_type === 'txn') {
							// emit(doc.id, {_id: doc.id});
							emit([` + "`" + `${doc["round-time"]}` + "`" + `, doc.id], null);
						}
					}`,
				},
				TransactionViewById: map[string]interface{}{
					"map": `function(doc) { 
						if (doc.doc_type === 'txn') {
							// emit(doc.id, {_id: doc.id});
							emit(doc.id, null);
						}
					}`,
				},
				// https://stackoverflow.com/questions/11284383/couchdb-count-unique-document-field
				TransactionViewByIdCount: map[string]interface{}{
					"map": `function(doc) {
						if (doc.doc_type === 'txn') {
							// emit(doc.id, 1);
							emit([` + "`" + `${doc["round-time"]}` + "`" + `, doc.id], 1);
						}
					}`,
					"reduce": "_sum",
				},
				TransactionViewByAccount: map[string]interface{} {
					"map": `function(doc) {
						if (doc.doc_type === 'acct') {
							emit([doc._id, "0"], null);
						} else if (doc.doc_type === 'txn') {
							doc.associated_accounts.forEach(acct => {
								emit([acct, "1", ` + "`" + `${doc["round-time"]}` + "`" + `, doc.id], null);
							})
						}
					}`,
				},
				// https://stackoverflow.com/questions/13216640/couchdb-getting-number-of-keys-in-given-key-range
				TransactionViewByAccountCount: map[string]interface{} {
					"map": `function(doc) {
						if (doc.doc_type === 'txn') {
							doc.associated_accounts.forEach(acct => {
								emit([acct, ` + "`" + `${doc["round-time"]}` + "`" + `, doc.id], 1);
							})
						}
					}`,
					"reduce": `function(keys, values, rereduce) {
						return sum(values);
					}`,
				},
				TransactionViewByAsset: map[string]interface{} {
					"map": `function(doc) {
						if (doc.doc_type === 'asset') {
							emit([doc._id, "0"], null);
						} else if (doc.doc_type === 'txn') {
							doc.associated_assets.forEach(asset => {
								emit([asset, "1", ` + "`" + `${doc["round-time"]}` + "`" + `, doc.id], null);
							})
						}
					}`,
				},
				TransactionViewByAssetCount: map[string]interface{} {
					"map": `function(doc) {
						if (doc.doc_type === 'txn') {
							doc.associated_assets.forEach(asset => {
								emit([asset, ` + "`" + `${doc["round-time"]}` + "`" + `, doc.id], 1);
							})
						}
					}`,
					"reduce": `function(keys, values, rereduce) {
						return sum(values);
					}`,
				},
				TransactionViewByApplication: map[string]interface{} {
					"map": `function(doc) {
						if (doc.doc_type === 'app') {
							emit([doc._id, "0"], null);
						} else if (doc.doc_type === 'txn') {
							doc.associated_applications.forEach(app => {
								emit([app, "1", ` + "`" + `${doc["round-time"]}` + "`" + `, doc.id], null);
							})
						}
					}`,
				},
				TransactionViewByApplicationCount: map[string]interface{} {
					"map": `function(doc) {
						if (doc.doc_type === 'txn') {
							doc.associated_applications.forEach(app => {
								emit([app, ` + "`" + `${doc["round-time"]}` + "`" + `, doc.id], 1);
							})
						}
					}`,
					"reduce": `function(keys, values, rereduce) {
						return sum(values);
					}`,
				},
			},
		})
		if err != nil && err.Error() != "Conflict: Document update conflict." {
			return fmt.Errorf("%s database and query by timestamp transaction view failed to be created: %w", dbName, err)
		}
	}
	return nil
}

// InsertAcctViewsForGlobalDB creates a the latest view for the acct design document. It stores
// transaction data.
func InsertAcctViewsForGlobalDB(ctx context.Context, client *kivik.Client, dbName string) error {
	// Check if DB exists
	exist, err := client.DBExists(ctx, dbName)
	if err != nil || !exist {
		return errors.Wrap(err, dbName + " database check fails")
	}
	db := client.DB(dbName)

	rows, err := db.Query(ctx, AccountDDoc, "_view/" +AccountViewByIdInLatest)
	//if err != nil {
	//	return errors.Wrap(err, dbName + " database and query by timestamp view failed to be queried")
	//}
	//if !rows.Next() {
	if rows == nil || !rows.Next() {
	_, err = db.Put(context.TODO(), AccountDDoc, map[string]interface{}{
		"_id": AccountDDoc,
		"views": map[string]interface{}{
			AccountViewByIdInLatest: map[string]interface{}{
				"map": `function(doc) { 
					if (doc.doc_type === 'acct') {
						// emit(doc.id, {_id: doc.id});
						emit(doc._id, null);
					}
				}`,
			},
			AccountViewByIdInCount: map[string]interface{}{
				"map": `function(doc) {
					if (doc.doc_type === 'acct') {
						emit(doc._id, 1);
					}
				}`,
				"reduce": "_sum",
			},
		},
	})
	if err != nil && err.Error() != "Conflict: Document update conflict." {
		return fmt.Errorf("%s database and query by timestamp account view failed to be created: %w", dbName, err)
	}
	}
	return nil
}

// InsertAssetViewsForGlobalDB creates a the latest view for the asset design document. It stores
// asset data.
func InsertAssetViewsForGlobalDB(ctx context.Context, client *kivik.Client, dbName string) error {
	// Check if DB exists
	exist, err := client.DBExists(ctx, dbName)
	if err != nil || !exist {
		return errors.Wrap(err, dbName + " database check fails")
	}
	db := client.DB(dbName)

	rows, err := db.Query(ctx, AssetDDoc, "_view/" +AssetViewByIdInLatest)
	//if err != nil {
	//	return errors.Wrap(err, dbName + " database and query by timestamp view failed to be queried")
	//}
	//if !rows.Next() {
	if rows == nil || !rows.Next() {
	_, err = db.Put(context.TODO(), AssetDDoc, map[string]interface{}{
		"_id": AssetDDoc,
		"views": map[string]interface{}{
			AssetViewByIdInLatest: map[string]interface{}{
				"map": `function(doc) { 
					if (doc.doc_type === 'asset') {
						// emit(doc.id, {_id: doc.id});
						emit(doc.id, null);
					}
				}`,
			},
			AssetViewByIdInCount: map[string]interface{}{
				"map": `function(doc) {
					if (doc.doc_type === 'asset') {
						emit(doc.id, 1);
					}
				}`,
				"reduce": "_sum",
			},
		},
	})
	if err != nil && err.Error() != "Conflict: Document update conflict." {
		return fmt.Errorf("%s database and query by timestamp asset view failed to be created: %w", dbName, err)
	}
	}
	return nil
}

// InsertApplicationViewsForGlobalDB creates a the latest view for the app design document. It stores
// application data.
func InsertApplicationViewsForGlobalDB(ctx context.Context, client *kivik.Client, dbName string) error {
	// Check if DB exists
	exist, err := client.DBExists(ctx, dbName)
	if err != nil || !exist {
		return errors.Wrap(err, dbName + " database check fails")
	}
	db := client.DB(dbName)

	rows, err := db.Query(ctx, ApplicationDDoc, "_view/" +ApplicationViewByIdInLatest)
	//if err != nil {
	//	return errors.Wrap(err, dbName + " database and query by timestamp view failed to be queried")
	//}
	//if !rows.Next() {
	if rows == nil || !rows.Next() {
	_, err = db.Put(context.TODO(), ApplicationDDoc, map[string]interface{}{
		"_id": AssetDDoc,
		"views": map[string]interface{}{
			ApplicationViewByIdInLatest: map[string]interface{}{
				"map": `function(doc) {
					if (doc.doc_type === 'app') {
						// emit(doc.id, {_id: doc.id});
						emit(doc.id, null);
					}
				}`,
			},
			ApplicationViewByIdInCount: map[string]interface{}{
				"map": `function(doc) {
					if (doc.doc_type === 'app') {
						emit(doc.id, 1);
					}
				}`,
				"reduce": "_sum",
			},
		},
	})
	if err != nil && err.Error() != "Conflict: Document update conflict." {
		return fmt.Errorf("%s database and query by timestamp application view failed to be created: %w", dbName, err)
	}
	}
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
		return fmt.Errorf("status check database: %w", err)
	}

	// Create the databases, delete the old ones if exist
	//var dbList = []string{"blocks", "transactions", "addresses"}
	//var dbList = []string{GlobalDbName}
	//for _, dbName := range dbList {
	//	if err := createDB(ctx, db, dbName); err != nil {
	//		return errors.Wrap(err, dbName + " database creation fails")
	//	}
	//}

	// Block views
	fmt.Println("Block views")
	if err := InsertBlockViewsForGlobalDB(ctx, db, GlobalDbName); err != nil {
		fmt.Printf("database fails to create view(s) for blocks: %s", err)
		return fmt.Errorf("database fails to create view(s) for blocks: %w", err)
	}

	// Transaction views
	fmt.Println("Transaction views")
	if err := InsertTransactionViewsForGlobalDB(ctx, db, GlobalDbName); err != nil {
		fmt.Printf("database fails to create view(s) for transactions: %s", err)
		return fmt.Errorf("database fails to create view(s) for transactions: %w", err)
	}

	// Account views
	fmt.Println("Account views")
	if err := InsertAcctViewsForGlobalDB(ctx, db, GlobalDbName); err != nil {
		fmt.Printf("database fails to create view(s) for accounts: %s", err)
		return fmt.Errorf("database fails to create view(s) for accounts: %w", err)
	}

	// Asset views
	fmt.Println("Asset views")
	if err := InsertAssetViewsForGlobalDB(ctx, db, GlobalDbName); err != nil {
		fmt.Printf("database fails to create view(s) for assets: %s", err)
		return fmt.Errorf("database fails to create view(s) for assets: %w", err)
	}

	// Application views
	fmt.Println("Application views")
	if err := InsertApplicationViewsForGlobalDB(ctx, db, GlobalDbName); err != nil {
		fmt.Printf("database fails to create view(s) for applications: %s", err)
		return fmt.Errorf("database fails to create view(s) for applications: %w", err)
	}

	// TODO: To-Be-Deleted
	//if err := InsertQueryViewForTransactionDB(ctx, db, "transactions"); err != nil {
	//	return errors.Wrap(err, "transactions database fails to create query view")
	//}
	//if err := InsertLatestViewForBlocksDB(ctx, db, "blocks"); err != nil {
	//	return errors.Wrap(err, "blocks database fails to create latest view")
	//}

	return nil
}
