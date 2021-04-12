package ledgergrp

import (
	"context"
	"fmt"
	"github.com/kevguy/algosearch/backend/business/core/algod"
	"github.com/kevguy/algosearch/backend/foundation/web"
	"net/http"
)

type Handlers struct {
	AlgodCore algod.Core
}

func (h Handlers) GetLedgerSupplyFromAPI(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	supply, err := h.AlgodCore.GetSupply(ctx, v.TraceID)
	if err != nil {
		return fmt.Errorf("getting ledger supplt info from algod api: %w", err)
	}
	return web.Respond(ctx, w, supply, http.StatusOK)

}
