package handlers

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/kevguy/algosearch/backend/business/algorand"
	"github.com/kevguy/algosearch/backend/business/couchdata/block"
	"github.com/kevguy/algosearch/backend/business/sys/validate"
	"github.com/kevguy/algosearch/backend/foundation/web"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type roundGroup struct {
	log			*zap.SugaredLogger
	store		block.Store
	algodClient	*algod.Client
}

// getCurrentRoundFromAPI retrieves the current round and returns the block data from Algod API
func (rG roundGroup) getCurrentRoundFromAPI(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	blockData, err := algorand.GetCurrentRound(ctx, v.TraceID, rG.log, rG.algodClient)
	if err != nil {
		return errors.Wrap(err, "unable to get current round")
	}

	return web.Respond(ctx, w, blockData, http.StatusOK)
}

// getRoundFromAPI retrieves a block from the Algod API based on the round number (num)
func (rG roundGroup) getRoundFromAPI(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	numStr := web.Param(r, "num")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return validate.NewRequestError(fmt.Errorf("invalid num format: %s", numStr), http.StatusBadRequest)
	}

	blockData, err := algorand.GetRound(ctx, v.TraceID, rG.log, rG.algodClient, uint64(num))
	if err != nil {
		return errors.Wrapf(err, "unable to get round %d", num)
	}

	return web.Respond(ctx, w, blockData, http.StatusOK)
}

// getRound retrieves a block from CouchDB based on the round number (num)
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

	blockData, err := rG.store.GetBlockByNum(ctx, v.TraceID, rG.log, uint64(num))
	if err != nil {
		return errors.Wrapf(err, "unable to get round %d", num)
	}

	return web.Respond(ctx, w, blockData, http.StatusOK)

}

// getLatestRound retrieves the latest block from CouchDB.
func (rG roundGroup) getLatestRound(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	blockData, err := rG.store.GetLatestBlock(ctx, v.TraceID, rG.log)
	if err != nil {
		return errors.Wrapf(err, "unable to get latest round")
	}

	return web.Respond(ctx, w, blockData, http.StatusOK)

}
