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

// GetAsset retrieves asset info from the Algod API based on the asset ID given
func GetAsset(ctx context.Context, traceID string, log *zap.SugaredLogger, algodClient *algodv2.Client, assetID uint64) (*models.Asset, error) {

	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "algod.GetAsset")
	span.SetAttributes(attribute.Int64("asset", int64(assetID)))
	defer span.End()

	log.Infow("algod.GetAsset", "traceid", traceID)

	assetInfo, err := algodClient.GetAssetByID(assetID).Do(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to query for asset info")
	}

	return &assetInfo, nil
}
