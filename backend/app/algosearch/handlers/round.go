package handlers

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/kevguy/algosearch/backend/business/algorand"
	"github.com/kevguy/algosearch/backend/business/sys/validate"
	"github.com/kevguy/algosearch/backend/foundation/web"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type roundGroup struct {
	log *zap.SugaredLogger
	algodClient *algod.Client
}

func (rG roundGroup) getCurrentRound(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	block, err := algorand.GetCurrentRound(ctx, v.TraceID, rG.log, rG.algodClient)
	if err != nil {
		return errors.Wrap(err, "unable to get current round")
	}

	return web.Respond(ctx, w, block, http.StatusOK)
}

func (rG roundGroup) getRound(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	numStr := web.Param(r, "num")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return validate.NewRequestError(fmt.Errorf("invalid num format: %s", numStr), http.StatusBadRequest)
	}

	block, err := algorand.GetRound(ctx, v.TraceID, rG.log, rG.algodClient, uint64(num))
	if err != nil {
		return errors.Wrap(err, "unable to get current round")
	}

	return web.Respond(ctx, w, block, http.StatusOK)
}
