package indexer

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// GetAsset retrieves asset info from the Indexer API based on the asset ID given.
func (c Core) GetAsset(ctx context.Context, assetID uint64) (models.Asset, error) {

	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "indexer.GetApplication")
	span.SetAttributes(attribute.Int64("asset", int64(assetID)))
	defer span.End()

	_, assetInfo, err := c.client.LookupAssetByID(assetID).Do(ctx)
	if err != nil {
		return models.Asset{}, errors.Wrap(err, "Unable to find application. Record may not exist in Postgre database.")
	}

	return assetInfo, nil
}
