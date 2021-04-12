package db

import (
	"context"
	"fmt"
	"github.com/go-kivik/kivik/v4"
	"github.com/kevguy/algosearch/backend/business/data/schema"
	"github.com/kevguy/algosearch/backend/foundation/web"
	"go.opentelemetry.io/otel"
)

// GetNumOfBlocks retrieves the number of blocks that exist in the database.
func (s Store) GetNumOfBlocks(ctx context.Context) (int64, error) {

	// http://0.0.0.0:5984/algo_global/_design/block/_view/blockByRoundNoCount?reduce=true&group_level=0

	ctx, span := otel.GetTracerProvider().
		Tracer("").
		Start(ctx, "block.GetNumOfBlocks")
	defer span.End()

	s.log.Infow("block.GetNumOfBlocks",
		"traceid", web.GetTraceID(ctx))

	exist, err := s.couchClient.DBExists(ctx, s.dbName)
	if err != nil || !exist {
		return 0, fmt.Errorf(s.dbName + " database check fails: %w", err)
	}
	db := s.couchClient.DB(s.dbName)

	// https://github.com/go-kivik/kivik/issues/246
	rows, err := db.Query(ctx, schema.BlockDDoc, "_view/" +schema.BlockViewByRoundCount, kivik.Options{
		"reduce": true,
		"group_level": 0,
	})
	if err != nil {
		return 0, fmt.Errorf("fetch data error: %w", err)
	}

	var count int64
	for rows.Next() {
		if err := rows.ScanValue(&count); err != nil {
			return 0, fmt.Errorf("can't find anything: %w", err)
		}
	}

	return count, nil
}
