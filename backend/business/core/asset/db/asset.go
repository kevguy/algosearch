// Package db contains asset related CRUD functionality.
package db

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/go-kivik/kivik/v4"
	"github.com/kevguy/algosearch/backend/business/data/schema"
	"github.com/kevguy/algosearch/backend/foundation/web"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"strconv"
)

const (
	DocType = "asset"
)

type Store struct {
	log *zap.SugaredLogger
	couchClient *kivik.Client
	dbName string
}

// NewStore constructs an asset store for api access.
func NewStore(log *zap.SugaredLogger, couchClient *kivik.Client, dbName string) Store {
	return Store{
		log: log,
		couchClient: couchClient,
		dbName: dbName,
	}
}

// AddAsset adds an asset to CouchDB.
// It receives the models.Asset object and transform it into an Asset document object and then
// insert it into the global CouchDB table.
func (s Store) AddAsset(ctx context.Context, asset models.Asset) (string, string, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "asset.AddAsset")
	defer span.End()

	s.log.Infow("asset.AddAsset", "traceid", web.GetTraceID(ctx))

	var doc = NewAsset{
		Asset:   asset,
		DocType: DocType,
	}
	//docID := fmt.Sprintf("%s.%s", DocType, doc.Id)
	exist, err := s.couchClient.DBExists(ctx, s.dbName)
	if err != nil || !exist {
		return "", "", errors.Wrap(err, s.dbName+ " database check fails")
	}
	db := s.couchClient.DB(s.dbName)

	docID := strconv.FormatUint(doc.Index, 10)
	rev, err := db.Put(ctx, docID, doc)
	if err != nil {
		return "", "", errors.Wrap(err, s.dbName+ " database can't insert asset id " +docID)
	}
	return docID, rev, nil
}

// AddAssets bulk-adds assets to CouchDB.
// It receives the []models.Asset object and transform them into Asset document objects and then
// insert them into the global CouchDB table.
func (s Store) AddAssets(ctx context.Context, assets []models.Asset) (bool, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "asset.AddAssets")
	defer span.End()

	s.log.Infow("asset.AddAssets", "traceid", web.GetTraceID(ctx))

	exist, err := s.couchClient.DBExists(ctx, s.dbName)
	if err != nil || !exist {
		return false, errors.Wrap(err, s.dbName+ " database check fails")
	}
	db := s.couchClient.DB(s.dbName)

	assets_ := make([]interface{}, len(assets))
	//fmt.Println("Here are teh assets")
	//fmt.Printf("%v\n", assets)

	// https://stackoverflow.com/questions/55755929/go-convert-interface-to-map
	// https://stackoverflow.com/questions/44094325/add-data-to-interface-in-struct
	for i := range assets {
		docID := strconv.FormatUint(assets[i].Index, 10)
		doc := NewAsset{
			ID:      &docID,
			Asset:   assets[i],
			DocType: DocType,
		}
		assets_[i] = doc
		//fmt.Println("YYYYYYYYYY")
		//fmt.Printf("%v\n", assets_[i])
		//v, _ := assets_[i].(map[string]interface{})
		//fmt.Println("VVVVVVVV")
		//fmt.Printf("%v\n", v)
		//v["_id"] = assets[i].Id
		//assets_[i] = v
		//fmt.Println("looping")
		//fmt.Println(assets_[i])
	}

	_, err = db.BulkDocs(ctx, assets_)
	if err != nil {
		return false, errors.Wrap(err, "Can't bulk insert the assets")
	}

	return true, nil
}

// GetAsset retrieves a asset record from CouchDB based upon the asset ID given.
func (s Store) GetAsset(ctx context.Context, assetID string) (models.Asset, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "asset.GetAsset")
	span.SetAttributes(attribute.String("assetID", assetID))
	defer span.End()

	s.log.Infow("asset.GetAsset", "traceid", web.GetTraceID(ctx), "assetID", assetID)

	exist, err := s.couchClient.DBExists(ctx, s.dbName)
	if err != nil || !exist {
		return models.Asset{}, errors.Wrap(err, s.dbName+ " database check fails")
	}
	db := s.couchClient.DB(s.dbName)

	docID := fmt.Sprintf("%s.%s", DocType, assetID)
	row := db.Get(ctx, docID)
	if row == nil {
		return models.Asset{}, errors.Wrap(err, s.dbName+ " get data empty")
	}

	var asset Asset
	//fmt.Printf("%v\n", row)
	err = row.ScanDoc(&asset)
	if err != nil {
		return models.Asset{}, errors.Wrap(err, s.dbName+ "cannot unpack data from row")
	}

	return asset.Asset, nil
}

func (s Store) GetEarliestAssetID(ctx context.Context) (string, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "asset.GetEarliestAssetID")
	//span.SetAttributes(attribute.String("query", q))
	defer span.End()

	s.log.Infow("asset.GetEarliestAssetID", "traceid", web.GetTraceID(ctx))

	exist, err := s.couchClient.DBExists(ctx, s.dbName)
	if err != nil || !exist {
		return "", errors.Wrap(err, s.dbName+ " database check fails")
	}
	db := s.couchClient.DB(s.dbName)

	rows, err := db.Query(ctx, schema.BlockDDoc, "_view/" +schema.AssetViewByIDInLatest, kivik.Options{
		"include_docs": true,
		"descending": false,
		"limit": 1,
	})
	if err != nil {
		return "", errors.Wrap(err, "Fetch data error")
	}

	if rows.Err() != nil {
		return "", errors.Wrap(err, "rows error, Can't find anything")
	}

	rows.Next()
	var doc Asset
	if err := rows.ScanDoc(&doc); err != nil {
		// No docs can be found
		return "", errors.Wrap(err, "Can't find anything")
	}

	docID := fmt.Sprintf("%s.%s", DocType, doc.Index)
	return docID, nil
}

func (s Store) GetLatestAssetID(ctx context.Context) (string, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "asset.GetLatestAssetID")
	//span.SetAttributes(attribute.String("query", q))
	defer span.End()

	s.log.Infow("asset.GetLatestAssetID", "traceid", web.GetTraceID(ctx))

	exist, err := s.couchClient.DBExists(ctx, s.dbName)
	if err != nil || !exist {
		return "", errors.Wrap(err, s.dbName+ " database check fails")
	}
	db := s.couchClient.DB(s.dbName)

	rows, err := db.Query(ctx, schema.BlockDDoc, "_view/" +schema.AssetViewByIDInLatest, kivik.Options{
		"include_docs": true,
		"descending": true,
		"limit": 1,
	})
	if err != nil {
		return "", errors.Wrap(err, "Fetch data error")
	}

	if rows.Err() != nil {
		return "", errors.Wrap(err, "rows error, Can't find anything")
	}

	rows.Next()
	var doc Asset
	if err := rows.ScanDoc(&doc); err != nil {
		// No docs can be found
		return "", errors.Wrap(err, "Can't find anything")
	}

	docID := fmt.Sprintf("%s.%d", DocType, doc.Index)
	return docID, nil
}

// GetAssetCountBtnKeys retrieves the number of keys between two keys.
// References:
// 	https://stackoverflow.com/questions/11284383/couchdb-count-unique-document-field
// 	https://stackoverflow.com/questions/12944294/using-a-couchdb-view-can-i-count-groups-and-filter-by-key-range-at-the-same-tim
func (s Store) GetAssetCountBtnKeys(ctx context.Context, startKey, endKey string) (int64, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "asset.GetAssetCountBtnKeys")
	span.SetAttributes(attribute.String("startKey", startKey))
	span.SetAttributes(attribute.String("endKey", endKey))
	defer span.End()

	s.log.Infow("asset.GetAssetCountBtnKeys",
		"traceid", web.GetTraceID(ctx),
		"startKey", startKey,
		"endKey", endKey)

	exist, err := s.couchClient.DBExists(ctx, s.dbName)
	if err != nil || !exist {
		return 0, errors.Wrap(err, s.dbName+ " database check fails")
	}
	db := s.couchClient.DB(s.dbName)

	rows, err := db.Query(ctx, schema.BlockDDoc, "_view/" +schema.AssetViewByIDInCount, kivik.Options{
		"start_key": startKey,
		"end_key": endKey,
	})
	if err != nil {
		return 0, errors.Wrap(err, "Fetch data error")
	}

	type Payload struct {
		Key *string `json:"key"`
		Value int64 `json:"value"`
	}

	var payload Payload
	for rows.Next() {
		if err := rows.ScanDoc(&payload); err != nil {
			return 0, errors.Wrap(err, "Can't find anything")
		}
	}

	return payload.Value, nil
}

func (s Store) GetAssetsPagination(ctx context.Context, latestAssetID string, order string, pageNo, limit int64) ([]Asset, int64, int64, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "block.GetBlocksPagination")
	span.SetAttributes(attribute.String("latestAssetID", latestAssetID))
	span.SetAttributes(attribute.Int64("pageNo", pageNo))
	span.SetAttributes(attribute.Int64("limit", limit))
	defer span.End()

	s.log.Infow("asset.GetAssetsPagination",
		"traceid", web.GetTraceID(ctx),
		"latestAssetID", latestAssetID,
		"pageNo", pageNo,
		"limit", limit)

	// Get the earliest asset id
	earliestTxnID, err := s.GetEarliestAssetID(ctx)
	if err != nil {
		return nil, 0, 0, errors.Wrap(err, ": Get earliest synced asset id")
	}

	numOfAssets, err := s.GetAssetCountBtnKeys(ctx, earliestTxnID, latestAssetID)
	if err != nil {
		return nil, 0, 0, errors.Wrap(err, ": Get asset count between keys")
	}

	// We can skip database check cuz GetEarliestAssetID already did it
	db := s.couchClient.DB(s.dbName)

	var numOfPages = numOfAssets / limit
	if numOfAssets % limit > 0 {
		numOfPages += 1
	}

	if pageNo < 1 || pageNo > numOfPages {
		return nil, 0, 0, errors.Wrapf(err, "page number is less than 1 or exceeds page limit: %d", numOfPages)
	}

	options := kivik.Options{
		"include_docs": true,
		"limit": limit,
	}

	if order == "desc" {
		// Descending order
		options["descending"] = true

		// Start with latest block number
		options["start_key"] = latestAssetID

		// Use page number to calculate number of items to skip
		skip := (pageNo - 1) * limit
		options["skip"] = (pageNo - 1) * limit

		// Find the key to start reading and get the `page limit` number of records
		if (numOfAssets - skip) > limit {
			options["limit"] = limit
		} else {
			options["limit"] = numOfAssets - skip
		}
	} else {
		// Ascending order
		options["descending"] = false

		// Calculate the number of records to skip
		skip := (pageNo - 1) * limit
		options["skip"] = skip

		if (numOfAssets - skip) > limit {
			options["limit"] =  numOfAssets - skip
		} else {
			options["limit"] = limit
		}
	}

	rows, err := db.Query(ctx, schema.AssetDDoc, "_view/" +schema.AssetViewByIDInLatest, options)
	if err != nil {
		return nil, 0, 0, errors.Wrap(err, "Fetch data error")
	}

	var fetchedAssets = []Asset{}
	for rows.Next() {
		var asset = Asset{}
		if err := rows.ScanDoc(&asset); err != nil {
			return nil, 0, 0, errors.Wrap(err, "unwrapping block")
		}
		fetchedAssets = append(fetchedAssets, asset)
	}

	if rows.Err() != nil {
		return nil, 0, 0, errors.Wrap(err, "rows error, Can't find anything")
	}

	return fetchedAssets, numOfPages, numOfAssets, nil
}
