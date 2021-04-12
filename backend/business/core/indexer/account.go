package indexer

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// GetAccount retrieves account info from the Indexer API based on the account address
// given.
func (c Core) GetAccount(ctx context.Context, accountID string) (models.Account, error) {

	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "indexer.GetAccount")
	span.SetAttributes(attribute.String("account", accountID))
	defer span.End()

	_, accountInfo, err := c.client.LookupAccountByID(accountID).Do(ctx)
	if err != nil {
		return models.Account{}, errors.Wrap(err, "Unable to find account. Record may not exist in Postgre database.")
	}

	return accountInfo, nil
}
