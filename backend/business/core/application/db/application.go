// Package db contains application related CRUD functionality.
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
	DocType = "app"
)

type Store struct {
	log *zap.SugaredLogger
	couchClient *kivik.Client
	dbName string
}

// NewStore constructs an application store for api access.
func NewStore(log *zap.SugaredLogger, couchClient *kivik.Client, dbName string) Store {
	return Store{
		log: log,
		couchClient: couchClient,
		dbName: dbName,
	}
}

// AddApplication adds an application to CouchDB.
// It receives the models.Application object and transform it into a Application document object and then
// insert it into the global CouchDB table.
func (s Store) AddApplication(ctx context.Context, application models.Application) (string, string, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "application.AddApplication")
	defer span.End()

	s.log.Infow("application.AddApplication", "traceid", web.GetTraceID(ctx))

	var doc = NewApplication{
		Application: application,
		DocType:     DocType,
	}
	//docID := fmt.Sprintf("%s.%s", DocType, doc.Id)
	exist, err := s.couchClient.DBExists(ctx, s.dbName)
	if err != nil || !exist {
		return "", "", errors.Wrap(err, s.dbName+ " database check fails")
	}
	db := s.couchClient.DB(s.dbName)

	docID := strconv.FormatUint(doc.Id, 10)
	rev, err := db.Put(ctx, docID, doc)
	if err != nil {
		return "", "", errors.Wrap(err, s.dbName+ " database can't insert application id " +docID)
	}
	return docID, rev, nil
}

// AddApplications bulk-adds applications to CouchDB.
// It receives the []models.Application object and transform them into Application document objects and then
// insert them into the global CouchDB table.
func (s Store) AddApplications(ctx context.Context, applications []models.Application) (bool, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "application.AddApplications")
	defer span.End()

	s.log.Infow("application.AddApplications", "traceid", web.GetTraceID(ctx))

	exist, err := s.couchClient.DBExists(ctx, s.dbName)
	if err != nil || !exist {
		return false, errors.Wrap(err, s.dbName+ " database check fails")
	}
	db := s.couchClient.DB(s.dbName)

	applications_ := make([]interface{}, len(applications))
	//fmt.Println("Here are teh applications")
	//fmt.Printf("%v\n", applications)

	// https://stackoverflow.com/questions/55755929/go-convert-interface-to-map
	// https://stackoverflow.com/questions/44094325/add-data-to-interface-in-struct
	for i := range applications {
		docID := strconv.FormatUint(applications[i].Id, 10)
		doc := NewApplication{
			ID:          &docID,
			Application: applications[i],
			DocType:     DocType,
		}
		applications_[i] = doc
	}

	_, err = db.BulkDocs(ctx, applications_)
	if err != nil {
		return false, errors.Wrap(err, "Can't bulk insert the applications")
	}

	return true, nil
}

// GetApplication retrieves an application record from CouchDB based upon the application ID given.
func (s Store) GetApplication(ctx context.Context, applicationID string) (models.Application, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "application.GetApplication")
	span.SetAttributes(attribute.String("applicationID", applicationID))
	defer span.End()

	s.log.Infow("application.GetApplication", "traceid", web.GetTraceID(ctx), "applicationID", applicationID)

	exist, err := s.couchClient.DBExists(ctx, s.dbName)
	if err != nil || !exist {
		return models.Application{}, errors.Wrap(err, s.dbName+ " database check fails")
	}
	db := s.couchClient.DB(s.dbName)

	docID := fmt.Sprintf("%s.%s", DocType, applicationID)
	row := db.Get(ctx, docID)
	if row == nil {
		return models.Application{}, errors.Wrap(err, s.dbName+ " get data empty")
	}

	var application Application
	fmt.Printf("%v\n", row)
	err = row.ScanDoc(&application)
	if err != nil {
		return models.Application{}, errors.Wrap(err, s.dbName+ "cannot unpack data from row")
	}

	return application.Application, nil
}

func (s Store) GetEarliestApplicationID(ctx context.Context) (string, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "application.GetEarliestApplicationID")
	//span.SetAttributes(attribute.String("query", q))
	defer span.End()

	s.log.Infow("application.GetEarliestApplicationID", "traceid", web.GetTraceID(ctx))

	exist, err := s.couchClient.DBExists(ctx, s.dbName)
	if err != nil || !exist {
		return "", errors.Wrap(err, s.dbName+ " database check fails")
	}
	db := s.couchClient.DB(s.dbName)

	rows, err := db.Query(ctx, schema.BlockDDoc, "_view/" +schema.ApplicationViewByIDInLatest, kivik.Options{
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
	var doc Application
	if err := rows.ScanDoc(&doc); err != nil {
		// No docs can be found
		return "", errors.Wrap(err, "Can't find anything")
	}

	docID := fmt.Sprintf("%s.%s", DocType, doc.Id)
	return docID, nil
}

func (s Store) GetLatestApplicationID(ctx context.Context) (string, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "application.GetLatestApplicationID")
	//span.SetAttributes(attribute.String("query", q))
	defer span.End()

	s.log.Infow("application.GetLatestApplicationID", "traceid", web.GetTraceID(ctx))

	exist, err := s.couchClient.DBExists(ctx, s.dbName)
	if err != nil || !exist {
		return "", errors.Wrap(err, s.dbName+ " database check fails")
	}
	db := s.couchClient.DB(s.dbName)

	rows, err := db.Query(ctx, schema.BlockDDoc, "_view/" +schema.ApplicationViewByIDInLatest, kivik.Options{
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
	var doc Application
	if err := rows.ScanDoc(&doc); err != nil {
		// No docs can be found
		return "", errors.Wrap(err, "Can't find anything")
	}

	docID := fmt.Sprintf("%s.%s", DocType, doc.Id)
	return docID, nil
}

// GetApplicationCountBtnKeys retrieves the number of keys between two keys.
// References:
// 	https://stackoverflow.com/questions/11284383/couchdb-count-unique-document-field
// 	https://stackoverflow.com/questions/12944294/using-a-couchdb-view-can-i-count-groups-and-filter-by-key-range-at-the-same-tim
func (s Store) GetApplicationCountBtnKeys(ctx context.Context, startKey, endKey string) (int64, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "application.GetApplicationCountBtnKeys")
	span.SetAttributes(attribute.String("startKey", startKey))
	span.SetAttributes(attribute.String("endKey", endKey))
	defer span.End()

	s.log.Infow("application.GetApplicationCountBtnKeys",
		"traceid", web.GetTraceID(ctx),
		"startKey", startKey,
		"endKey", endKey)

	exist, err := s.couchClient.DBExists(ctx, s.dbName)
	if err != nil || !exist {
		return 0, errors.Wrap(err, s.dbName+ " database check fails")
	}
	db := s.couchClient.DB(s.dbName)

	rows, err := db.Query(ctx, schema.BlockDDoc, "_view/" +schema.ApplicationViewByIDInCount, kivik.Options{
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

func (s Store) GetApplicationsPagination(ctx context.Context, latestApplicationID string, order string, pageNo, limit int64) ([]Application, int64, int64, error) {

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "block.GetBlocksPagination")
	span.SetAttributes(attribute.String("latestApplicationID", latestApplicationID))
	span.SetAttributes(attribute.Int64("pageNo", pageNo))
	span.SetAttributes(attribute.Int64("limit", limit))
	defer span.End()

	s.log.Infow("application.GetApplicationsPagination",
		"traceid", web.GetTraceID(ctx),
		"latestApplicationID", latestApplicationID,
		"pageNo", pageNo,
		"limit", limit)

	// Get the earliest application id
	earliestTxnID, err := s.GetEarliestApplicationID(ctx)
	if err != nil {
		return nil, 0, 0, errors.Wrap(err, ": Get earliest synced application id")
	}

	numOfApplications, err := s.GetApplicationCountBtnKeys(ctx, earliestTxnID, latestApplicationID)
	if err != nil {
		return nil, 0, 0, errors.Wrap(err, ": Get application count between keys")
	}

	// We can skip database check cuz GetEarliestApplicationID already did it
	db := s.couchClient.DB(s.dbName)

	var numOfPages = numOfApplications / limit
	if numOfApplications % limit > 0 {
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
		options["start_key"] = latestApplicationID

		// Use page number to calculate number of items to skip
		skip := (pageNo - 1) * limit
		options["skip"] = (pageNo - 1) * limit

		// Find the key to start reading and get the `page limit` number of records
		if (numOfApplications - skip) > limit {
			options["limit"] = limit
		} else {
			options["limit"] = numOfApplications - skip
		}
	} else {
		// Ascending order
		options["descending"] = false

		// Calculate the number of records to skip
		skip := (pageNo - 1) * limit
		options["skip"] = skip

		if (numOfApplications - skip) > limit {
			options["limit"] =  numOfApplications - skip
		} else {
			options["limit"] = limit
		}
	}

	rows, err := db.Query(ctx, schema.ApplicationDDoc, "_view/" +schema.ApplicationViewByIDInLatest, options)
	if err != nil {
		return nil, 0, 0, errors.Wrap(err, "Fetch data error")
	}

	var fetchedApplications = []Application{}
	for rows.Next() {
		var application = Application{}
		if err := rows.ScanDoc(&application); err != nil {
			return nil, 0, 0, errors.Wrap(err, "unwrapping block")
		}
		fetchedApplications = append(fetchedApplications, application)
	}

	if rows.Err() != nil {
		return nil, 0, 0, errors.Wrap(err, "rows error, Can't find anything")
	}

	return fetchedApplications, numOfPages, numOfApplications, nil
}
