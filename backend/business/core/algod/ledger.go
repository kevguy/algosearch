package algod

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
)

// GetSupply retrieves the current supply reported by the ledger from the Algod API.
func (c Core) GetSupply(ctx context.Context, traceID string) (*models.Supply, error) {

	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "algod.GetSupply")
	defer span.End()

	c.log.Infow("algod.GetSupply", "traceid", traceID)

	supply, err := c.algodClient.Supply().Do(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to query for ledger supply info")
	}

	return &supply, nil
}
