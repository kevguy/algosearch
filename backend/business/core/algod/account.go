package algod

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// GetAccount retrieves account info from the Algod API based on the account address given
func (c Core) GetAccount(ctx context.Context, traceID string, address string) (*models.Account, error) {

	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "algod.GetAccount")
	span.SetAttributes(attribute.String("account", address))
	defer span.End()

	c.log.Infow("algod.GetAccount", "traceid", traceID)

	accountInfo, err := c.algodClient.AccountInformation(address).Do(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to query for account info")
	}

	return &accountInfo, nil
}
