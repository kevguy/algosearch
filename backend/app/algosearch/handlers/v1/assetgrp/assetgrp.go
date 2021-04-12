package assetgrp

import (
	"context"
	"fmt"
	"github.com/kevguy/algosearch/backend/business/core/algod"
	v1web "github.com/kevguy/algosearch/backend/business/web/v1"
	"github.com/kevguy/algosearch/backend/foundation/web"

	//"github.com/kevguy/algosearch/backend/foundation/web"
	"net/http"
	"strconv"
)

type Handlers struct {
	AlgodCore algod.Core
}

func (h Handlers) GetAssetByIDFromAPI(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	idxStr := web.Param(r, "idx")
	idx, err := strconv.Atoi(idxStr)
	if err != nil {
		return v1web.NewRequestError(fmt.Errorf("invalid idx format: %s", idxStr), http.StatusBadRequest)
	}

	asset, err := h.AlgodCore.GetAsset(ctx, v.TraceID, uint64(idx))
	if err != nil {
		return fmt.Errorf("getting asset[%d] from algod api: %w", idx, err)
	}
	return web.Respond(ctx, w, asset, http.StatusOK)

}
