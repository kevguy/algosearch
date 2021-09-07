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

// GetAccount retrieves account info from the Algod API based on the account address given
func GetAccount(ctx context.Context, traceID string, log *zap.SugaredLogger, algodClient *algodv2.Client, address string) (*models.Account, error) {

	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "algod.GetAccount")
	span.SetAttributes(attribute.String("account", address))
	defer span.End()

	log.Infow("algod.GetAccount", "traceid", traceID)

	accountInfo, err := algodClient.AccountInformation(address).Do(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to query for account info")
	}

	return &accountInfo, nil
}
