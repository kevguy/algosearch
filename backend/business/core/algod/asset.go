package algod

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// GetAsset retrieves asset info from the Algod API based on the asset ID given
func (c Core) GetAsset(ctx context.Context, traceID string, assetID uint64) (*models.Asset, error) {

	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "algod.GetAsset")
	span.SetAttributes(attribute.Int64("asset", int64(assetID)))
	defer span.End()

	c.log.Infow("algod.GetAsset", "traceid", traceID)

	assetInfo, err := c.algodClient.GetAssetByID(assetID).Do(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to query for asset info")
	}

	return &assetInfo, nil
}
