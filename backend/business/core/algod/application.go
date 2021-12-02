package algod

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// GetApplication retrieves application info from the Algod API based on the application ID given
func (c Core) GetApplication(ctx context.Context, traceID string, appID uint64) (*models.Application, error) {

	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "algod.GetApplication")
	span.SetAttributes(attribute.Int64("asset", int64(appID)))
	defer span.End()

	c.log.Infow("algod.GetApplication", "traceid", traceID)

	appInfo, err := c.algodClient.GetApplicationByID(appID).Do(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to query for application info")
	}

	return &appInfo, nil
}
