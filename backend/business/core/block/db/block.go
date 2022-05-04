// Package db contains block related CRUD functionality.
package db

import (
	"context"
	"fmt"
	"github.com/go-kivik/kivik/v4"
	"github.com/kevguy/algosearch/backend/business/data/schema"
	"github.com/kevguy/algosearch/backend/foundation/web"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

const (
	DocType = "block"
)

type Store struct {
	log         *zap.SugaredLogger
	couchClient *kivik.Client
	dbName      string
}

// NewStore constructs a block store for api access.
func NewStore(log *zap.SugaredLogger, couchClient *kivik.Client, dbName string) Store {
	return Store{
		log:         log,
		couchClient: couchClient,
		dbName:      dbName,
	}
}

// AddBlock adds a block to CouchDB using block hash as ID.
func (s Store) AddBlock(ctx context.Context, block NewBlock) (string, string, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "block.AddBlock")
	defer span.End()

	s.log.Infow("block.AddBlock", "traceid", web.GetTraceID(ctx))

	var doc = NewBlockDoc{
		NewBlock: block,
		DocType:  DocType,
	}
	exist, err := s.couchClient.DBExists(ctx, s.dbName)
	if err != nil || !exist {
		return "", "", errors.Wrap(err, s.dbName+" database check fails")
	}
	db := s.couchClient.DB(s.dbName)

	//docID, rev, err := db.CreateDoc(ctx, block, map[string]interface{}{
	//	"_id": block.BlockHash,
	//	"key": strconv.FormatUint(block.Round, 10),
	//})
	fmt.Println("Block hash")
	fmt.Println(doc.BlockHash)
	fmt.Println(doc.Round)
	var docID = ""
	// This is to handle private network
	if doc.BlockHash == "" {
		docID = fmt.Sprintf("%d", doc.Round)
	} else {
		docID = doc.BlockHash
	}
	rev, err := db.Put(ctx, docID, doc)
	if err != nil {
		return "", "", errors.Wrapf(err, s.dbName+" database can't insert block number %d", block.Round)
	}
	//return strconv.FormatUint(block.Round, 10), rev, nil
	return block.BlockHash, rev, nil
}

// AddBlocks add blocks to CouchDB using their block hashes as IDs.
func (s Store) AddBlocks(ctx context.Context, blocks []Block) (bool, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "block.AddBlocks")
	defer span.End()

	s.log.Infow("block.AddBlocks", "traceid", web.GetTraceID(ctx))

	exist, err := s.couchClient.DBExists(ctx, s.dbName)
	if err != nil || !exist {
		return false, errors.Wrap(err, s.dbName+" database check fails")
	}
	db := s.couchClient.DB(s.dbName)

	blocks_ := make([]interface{}, len(blocks))
	for i := range blocks {
		blocks_[i] = blocks[i]
		v, _ := blocks_[i].(map[string]interface{})
		v["_id"] = blocks[i].BlockHash
		blocks_[i] = v
	}

	_, err = db.BulkDocs(ctx, blocks_)
	if err != nil {
		return false, errors.Wrap(err, "Can't bulk insert the blocks")
	}

	return true, nil
}

// GetBlockByHash retrieves a block from CouchDB based upon the block hash.
func (s Store) GetBlockByHash(ctx context.Context, blockHash string) (Block, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "block.GetBlockByHash")
	span.SetAttributes(attribute.String("blockHash", blockHash))
	defer span.End()

	s.log.Infow("block.GetBlockByHash", "traceid", web.GetTraceID(ctx), "blockHash", blockHash)

	exist, err := s.couchClient.DBExists(ctx, s.dbName)
	if err != nil || !exist {
		return Block{}, errors.Wrap(err, s.dbName+" database check fails")
	}
	db := s.couchClient.DB(s.dbName)

	//row := db.Get(ctx, strconv.FormatUint(blockNum, 10))
	row := db.Get(ctx, blockHash)
	if row == nil {
		return Block{}, errors.Wrap(err, s.dbName+" get data empty")
	}

	var block Block
	fmt.Printf("%v\n", row)
	err = row.ScanDoc(&block)
	if err != nil {
		return Block{}, errors.Wrap(err, s.dbName+"cannot unpack data from row")
	}

	return block, nil
}

// GetBlockByNum gets a block from CouchDB based on round number.
func (s Store) GetBlockByNum(ctx context.Context, blockNum uint64) (Block, error) {

	//ctx, span := otel.GetTracerProvider().
	//	Tracer("").
	//	Start(ctx, "block.GetBlockByNum")
	//span.SetAttributes(attribute.Any("block-num", blockNum))
	//defer span.End()

	//s.log.Infow("block.GetBlockByNum", "traceid", traceID, "blockNum", blockNum)

	exist, err := s.couchClient.DBExists(ctx, s.dbName)
	if err != nil || !exist {
		return Block{}, errors.Wrap(err, s.dbName+" database check fails")
	}
	db := s.couchClient.DB(s.dbName)

	rows, err := db.Query(ctx, schema.BlockDDoc, "_view/"+schema.BlockViewByRoundNo, kivik.Options{
		"include_docs": true,
		"key":          blockNum,
		"limit":        1,
	})
	if err != nil {
		return Block{}, errors.Wrap(err, "Fetch data error")
	}

	if rows.Err() != nil {
		return Block{}, errors.Wrap(err, "rows error, Can't find anything")
	}

	rows.Next()
	var doc Block
	if err := rows.ScanDoc(&doc); err != nil {
		// No docs can be found
		return Block{}, errors.Wrap(err, "Can't find anything")
	}

	return doc, nil
}

// GetEarliestSyncedRoundNumber retrieves the earliest round number that is synced to CouchDB.
func (s Store) GetEarliestSyncedRoundNumber(ctx context.Context) (uint64, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "block.GetEarliestSyncedRoundNumber")
	//span.SetAttributes(attribute.String("query", q))
	defer span.End()

	s.log.Infow("block.GetEarliestSyncedRoundNumber", "traceid", web.GetTraceID(ctx))

	exist, err := s.couchClient.DBExists(ctx, s.dbName)
	if err != nil || !exist {
		return 0, errors.Wrap(err, s.dbName+" database check fails")
	}
	db := s.couchClient.DB(s.dbName)

	rows, err := db.Query(ctx, schema.BlockDDoc, "_view/"+schema.BlockViewByRoundNo, kivik.Options{
		"include_docs": true,
		"descending":   false,
		"limit":        1,
	})
	if err != nil {
		return 0, errors.Wrap(err, "Fetch data error")
	}

	if rows.Err() != nil {
		return 0, errors.Wrap(err, "rows error, Can't find anything")
	}

	rows.Next()
	var doc Block
	if err := rows.ScanDoc(&doc); err != nil {
		// No docs can be found
		return 0, errors.Wrap(err, "Can't find anything")
	}

	return doc.Round, nil
	//rows, err := db.Query(context.TODO(), "_design/foo", "_view/bar", kivik.Options{
	//	"startkey": `"foo"`,                           // Quotes are necessary so the
	//	"endkey":   `"foo` + kivik.EndKeySuffix + `"`, // key is a valid JSON object
	//})
	//if err != nil {
	//	panic(err)
	//}
	//for rows.Next() {
	//	var doc interface{}
	//	if err := rows.ScanDoc(&doc); err != nil {
	//		panic(err)
	//	}
	//	/* do something with doc */
	//}
	//if rows.Err() != nil {
	//	panic(rows.Err())
	//}
}

// GetLastSyncedRoundNumber retrieves the last round number that is synced to CouchDB.
func (s Store) GetLastSyncedRoundNumber(ctx context.Context) (uint64, bool, error) {
	//func (s Store) GetLastSyncedRoundNumber(ctx context.Context, traceID string) (uint64, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "block.GetLastSyncedRoundNumber")
	//span.SetAttributes(attribute.String("query", q))
	defer span.End()

	s.log.Infow("block.GetLastSyncedRoundNumber", "traceid", web.GetTraceID(ctx))

	exist, err := s.couchClient.DBExists(ctx, s.dbName)
	if err != nil || !exist {
		return 0, false, errors.Wrap(err, s.dbName+" database check fails")
	}
	db := s.couchClient.DB(s.dbName)

	rows, err := db.Query(ctx, schema.BlockDDoc, "_view/"+schema.BlockViewByRoundNo, kivik.Options{
		"include_docs": true,
		"descending":   true,
		"limit":        1,
	})
	if err != nil {
		return 0, false, errors.Wrap(err, "Fetch data error")
	}

	if rows.Err() != nil {
		return 0, false, errors.Wrap(err, "rows error, Can't find anything")
	}

	rows.Next()
	var doc Block
	if err := rows.ScanDoc(&doc); err != nil {
		// No docs can be found
		return 0, false, errors.Wrap(err, "Can't find anything")
	}

	return doc.Round, true, nil
}

// GetLatestBlock retrieves the block that is last synced to Couch.
func (s Store) GetLatestBlock(ctx context.Context) (Block, error) {

	//log.Infow("block.GetLatestBlock", "traceid", traceID)

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "block.GetLatestBlock")
	defer span.End()

	s.log.Infow("block.GetLatestBlock", "traceid", web.GetTraceID(ctx))

	exist, err := s.couchClient.DBExists(ctx, s.dbName)
	if err != nil || !exist {
		return Block{}, errors.Wrap(err, s.dbName+" database check fails")
	}
	db := s.couchClient.DB(s.dbName)

	rows, err := db.Query(ctx, schema.BlockDDoc, "_view/"+schema.BlockViewByRoundNo, kivik.Options{
		"include_docs": true,
		"descending":   true,
		"limit":        1,
	})
	if err != nil {
		return Block{}, errors.Wrap(err, "Fetch data error")
	}

	if rows.Err() != nil {
		return Block{}, errors.Wrap(err, "rows error, Can't find anything")
	}

	rows.Next()
	var doc Block
	if err := rows.ScanDoc(&doc); err != nil {
		// No docs can be found
		return Block{}, errors.Wrap(err, "Can't find anything")
	}

	return doc, nil
}

// GetBlocksPagination retrieves a list of blocks based upon the following parameters:
// latestBlockNum: the latest block number that user knows about
// order: desc/asc
// pageNo: the number of pages the user wants to look at
// limit: number of blocks per page
// https://docs.couchdb.org/en/main/ddocs/views/pagination.html
func (s Store) GetBlocksPagination(ctx context.Context, latestBlockNum int64, order string, pageNo int64, limit int64) ([]Block, int64, int64, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "block.GetBlocksPagination")
	span.SetAttributes(attribute.Int64("latestBlockNum", latestBlockNum))
	span.SetAttributes(attribute.Int64("pageNo", pageNo))
	span.SetAttributes(attribute.Int64("limit", limit))
	defer span.End()

	s.log.Infow("block.GetBlocksPagination",
		"traceid", web.GetTraceID(ctx),
		"latestBlockNum", latestBlockNum,
		"pageNo", pageNo,
		"limit", limit)

	// Get the earliest block number
	earliestBlkNum, err := s.GetEarliestSyncedRoundNumber(ctx)
	if err != nil {
		return nil, 0, 0, errors.Wrapf(err, ": Get earliest synced round number")
	}
	// We can skip database check cuz GetEarliestSyncedRoundNumber already did it
	db := s.couchClient.DB(s.dbName)

	// We can basically treat latestBlockNum as number of blocks
	var numOfBlks = latestBlockNum - int64(earliestBlkNum) + 1
	var numOfPages = numOfBlks / limit
	if numOfBlks%limit > 0 {
		numOfPages += 1
	}

	if pageNo < 1 || pageNo > numOfPages {
		return nil, 0, 0, errors.Wrapf(err, "page number is less than 1 or exceeds page limit: %d", numOfPages)
	}

	options := kivik.Options{
		"include_docs": true,
		"limit":        limit,
	}

	if order == "desc" {
		// Descending order
		options["descending"] = true

		// Start with latest block number
		options["start_key"] = latestBlockNum

		// Use page number to calculate number of items to skip
		skip := (pageNo - 1) * limit
		options["skip"] = (pageNo - 1) * limit

		// Find the key to start reading and get the `page limit` number of records
		if ((latestBlockNum - int64(earliestBlkNum) + 1) - skip) > limit {
			options["limit"] = limit
		} else {
			options["limit"] = latestBlockNum - skip
		}
	} else {
		// Ascending order
		options["descending"] = false

		// Calculate the number of records to skip
		skip := (pageNo - 1) * limit
		options["skip"] = skip

		if (int64(earliestBlkNum) + skip + limit - 1 - latestBlockNum) > 0 {
			options["limit"] = latestBlockNum - skip
		} else {
			options["limit"] = limit
		}
	}

	//rows, err := db.Query(ctx, "_design/latest", "_view/latest", kivik.Options{
	//	"include_docs": true,
	//	"descending": true,
	//	"limit": limit,
	//	"skip": lastBlockNum - limit,
	//})
	rows, err := db.Query(ctx, schema.BlockDDoc, "_view/"+schema.BlockViewByRoundNo, options)
	if err != nil {
		return nil, 0, 0, errors.Wrap(err, "Fetch data error")
	}

	var fetchedBlocks = []Block{}
	for rows.Next() {
		var block = Block{}
		if err := rows.ScanDoc(&block); err != nil {
			return nil, 0, 0, errors.Wrap(err, "unwrapping block")
		}
		fetchedBlocks = append(fetchedBlocks, block)
	}

	if rows.Err() != nil {
		return nil, 0, 0, errors.Wrap(err, "rows error, Can't find anything")
	}

	return fetchedBlocks, numOfPages, numOfBlks, nil
}

// GetBlockTxnSpeed compares the latest 10 blocks it could find in the database
// and finds the average transaction speed
func (s Store) GetBlockTxnSpeed(ctx context.Context) (float64, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "block.GetBlockTxnTime")
	defer span.End()

	s.log.Infow("block.GetBlockTxnTime", "traceid", web.GetTraceID(ctx))

	exist, err := s.couchClient.DBExists(ctx, s.dbName)
	if err != nil || !exist {
		return 0.0, errors.Wrap(err, s.dbName+" database check fails")
	}
	db := s.couchClient.DB(s.dbName)

	// Get latest 10 blocks
	options := kivik.Options{
		"include_docs": true,
		"limit":        10,
		"descending":   true,
	}

	rows, err := db.Query(ctx, schema.BlockDDoc, "_view/"+schema.BlockViewByRoundNo, options)
	if err != nil {
		return 0.0, errors.Wrap(err, "Fetch data error")
	}

	var fetchedBlocks = []Block{}
	for rows.Next() {
		var block = Block{}
		if err := rows.ScanDoc(&block); err != nil {
			return 0.0, errors.Wrap(err, "unwrapping block")
		}
		fetchedBlocks = append(fetchedBlocks, block)
	}

	if len(fetchedBlocks) <= 1 {
		return 0.0, nil
	} else {
		var timeDiffs []uint64
		for idx, block := range fetchedBlocks {
			if idx > 1 {
				timeDiffs = append(timeDiffs, fetchedBlocks[idx-1].Timestamp-block.Timestamp)
			}
		}

		var sum uint64 = 0
		for _, item := range timeDiffs {
			sum += item
		}
		average := float64(sum) / float64(len(timeDiffs))
		return average, nil
	}
}
