package handlers

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	algodBiz "github.com/kevguy/algosearch/backend/business/algod"
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
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	blockData, err := algodBiz.GetCurrentRound(ctx, v.TraceID, rG.log, rG.algodClient)
	if err != nil {
		return errors.Wrap(err, "unable to get current round")
	}

	return web.Respond(ctx, w, blockData, http.StatusOK)
}

// getRoundFromAPI retrieves a block from the Algod API based on the round number (num)
func (rG roundGroup) getRoundFromAPI(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	numStr := web.Param(r, "num")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return validate.NewRequestError(fmt.Errorf("invalid num format: %s", numStr), http.StatusBadRequest)
	}

	blockData, err := algodBiz.GetRound(ctx, v.TraceID, rG.log, rG.algodClient, uint64(num))
	if err != nil {
		return errors.Wrapf(err, "unable to get round %d", num)
	}

	return web.Respond(ctx, w, blockData, http.StatusOK)
}

// getRound retrieves a block from CouchDB based on the round number (num)
func (rG roundGroup) getRound(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
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
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	blockData, err := rG.store.GetLatestBlock(ctx, v.TraceID, rG.log)
	if err != nil {
		return errors.Wrapf(err, "unable to get latest round")
	}

	return web.Respond(ctx, w, blockData, http.StatusOK)
}

func (rG roundGroup) getRoundsPagination(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	// limit
	limitQueries := web.Query(r, "limit")
	if len(limitQueries) == 0 {
		return validate.NewRequestError(fmt.Errorf("missing query parameter: limit"), http.StatusBadRequest)
	}
	limit, err := strconv.Atoi(limitQueries[0])
	if err != nil {
		return validate.NewRequestError(fmt.Errorf("invalid 'limit' format: %s", limitQueries[0]), http.StatusBadRequest)
	}

	// latest_blk
	latestBlkQueries := web.Query(r, "latest_blk")
	if len(latestBlkQueries) == 0 {
		return validate.NewRequestError(fmt.Errorf("missing query parameter: start"), http.StatusBadRequest)
	}
	latestBlk, err := strconv.Atoi(latestBlkQueries[0])
	if err != nil {
		return validate.NewRequestError(fmt.Errorf("invalid 'start' format: %s", latestBlkQueries[0]), http.StatusBadRequest)
	}

	// page
	pageQueries := web.Query(r, "page")
	if len(pageQueries) == 0 {
		return validate.NewRequestError(fmt.Errorf("missing query parameter: page"), http.StatusBadRequest)
	}
	page, err := strconv.Atoi(pageQueries[0])
	if err != nil {
		return validate.NewRequestError(fmt.Errorf("invalid 'page' format: %s", latestBlkQueries[0]), http.StatusBadRequest)
	}

	// order
	var order string
	orderQueries := web.Query(r, "order")
	if len(orderQueries) == 0 {
		//return validate.NewRequestError(fmt.Errorf("missing query parameter: sort"), http.StatusBadRequest)
		order = "desc"
	} else {
		order = orderQueries[0]
	}
	if order != "asc" && order != "desc" {
		return validate.NewRequestError(fmt.Errorf("invalid 'sort' format: %s", orderQueries[0]), http.StatusBadRequest)
	}

	result, numOfPages, numOfBlks, err := rG.store.GetBlocksPagination(ctx, v.TraceID, rG.log, int64(latestBlk), order, int64(page), int64(limit))
	if err != nil {
		return errors.Wrap(err, "Error fetching pagination results")
	}

	type Payload struct {
		NumOfPages	int64 `json:"num_of_pages"`
		NumOfBlks	int64 `json:"num_of_blks"`
		Items []block.Block `json:"items"`
	}

	return web.Respond(ctx, w, Payload{
		NumOfPages: numOfPages,
		NumOfBlks:  numOfBlks,
		Items:      result,
	}, http.StatusOK)
}
