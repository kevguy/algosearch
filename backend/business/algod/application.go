package algod

import (
	"context"
	algodv2 "github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// GetApplication retrieves application info from the Algod API based on the application ID given
func GetApplication(ctx context.Context, traceID string, log *zap.SugaredLogger, algodClient *algodv2.Client, appID uint64) (*models.Application, error) {

	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "algod.GetApplication")
	span.SetAttributes(attribute.Int64("asset", int64(appID)))
	defer span.End()

	log.Infow("algod.GetApplication", "traceid", traceID)

	appInfo, err := algodClient.GetApplicationByID(appID).Do(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to query for application info")
	}

	return &appInfo, nil
}
