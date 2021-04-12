package indexer

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// GetApplication retrieves application info from the Indexer API based on the application ID given.
func (c Core) GetApplication(ctx context.Context, appID uint64) (models.Application, error) {

	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "indexer.GetApplication")
	span.SetAttributes(attribute.Int64("application", int64(appID)))
	defer span.End()

	appInfo, err := c.client.LookupApplicationByID(appID).Do(ctx)
	if err != nil {
		return models.Application{}, errors.Wrap(err, "Unable to find application. Record may not exist in Postgre database.")
	}

	return appInfo.Application, nil
}
