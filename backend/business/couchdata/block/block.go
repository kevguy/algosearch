// Package block contains block related CRUD functionality.
package block

import (
	"context"
	"fmt"
	"github.com/go-kivik/kivik/v4"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"strconv"
)

const (
	BLOCKS = "blocks"
)

type Store struct {
	log *zap.SugaredLogger
	couchClient *kivik.Client
}

// NewStore constructs a product store for api access.
func NewStore(log *zap.SugaredLogger, couchClient *kivik.Client) Store {
	return Store{
		log: log,
		couchClient: couchClient,
	}
}

// AddBlock adds a block to CouchDB.
func (s Store) AddBlock(ctx context.Context, block NewBlock) (string, string, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "block.AddBlock")
	defer span.End()

	exist, err := s.couchClient.DBExists(ctx, BLOCKS)
	if err != nil || !exist {
		return "", "", errors.Wrap(err, BLOCKS+ " database check fails")
	}
	db := s.couchClient.DB(BLOCKS)

	//docID, rev, err := db.CreateDoc(ctx, block, map[string]interface{}{
	//	"_id": strconv.FormatUint(block.Round, 10),
	//	"key": strconv.FormatUint(block.Round, 10),
	//})
	rev, err := db.Put(ctx, strconv.FormatUint(block.Round, 10), block)
	if err != nil {
		return "", "", errors.Wrap(err, BLOCKS+ " database can't insert block number " + string(block.Round))
	}
	return strconv.FormatUint(block.Round, 10), rev, nil
}

// AddBlock adds a block to CouchDB.
func (s Store) GetBlock(ctx context.Context, blockNum uint64) (Block, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "block.AddBlock")
	defer span.End()

	exist, err := s.couchClient.DBExists(ctx, BLOCKS)
	if err != nil || !exist {
		return Block{}, errors.Wrap(err, BLOCKS+ " database check fails")
	}
	db := s.couchClient.DB(BLOCKS)

	row := db.Get(ctx, strconv.FormatUint(blockNum, 10))
	if row == nil {
		return Block{}, errors.Wrap(err, BLOCKS+ " get data empty")
	}

	var block Block
	fmt.Printf("%v\n", row)
	err = row.ScanDoc(&block)
	if err != nil {
		return Block{}, errors.Wrap(err, BLOCKS+ "cannot unpack data from row")
	}

	return block, nil
}

func (s Store) AddBlocks(ctx context.Context, blocks []Block) (bool, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "block.AddBlock")
	defer span.End()

	exist, err := s.couchClient.DBExists(ctx, BLOCKS)
	if err != nil || !exist {
		return false, errors.Wrap(err, BLOCKS+ " database check fails")
	}
	db := s.couchClient.DB(BLOCKS)

	blocks_ := make([]interface{}, len(blocks))
	for i := range blocks {
		blocks_[i] = blocks[i]
	}

	_, err = db.BulkDocs(ctx, blocks_)
	if err != nil {
		return false, errors.Wrap(err, "Can't bulk insert the blocks")
	}

	return true, nil
}

// GetLastSyncedRoundNumber retrieves the last round number that is synced to CouchDB.
func (s Store) GetLastSyncedRoundNumber(ctx context.Context) (uint64, error) {
//func (s Store) GetLastSyncedRoundNumber(ctx context.Context, traceID string) (uint64, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "block.GetLastSyncedRoundNumber")
	//span.SetAttributes(attribute.String("query", q))
	defer span.End()

	exist, err := s.couchClient.DBExists(ctx, BLOCKS)
	if err != nil || !exist {
		return 0, errors.Wrap(err, BLOCKS+ " database check fails")
	}
	db := s.couchClient.DB(BLOCKS)

	rows, err := db.Query(ctx, "_design/latest", "_view/latest", kivik.Options{
		"include_docs": true,
		"descending": true,
		"limit": 1,
	})

	if rows.Err() != nil {
		return 0, errors.Wrap(err, "Can't find anything")
	}

	rows.Next()
	var doc Block
	if err := rows.ScanDoc(&doc); err != nil {
		// No docs can be found
		return 0, errors.Wrap(err, "Can't find anything")
	}

	return uint64(doc.Round), nil
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