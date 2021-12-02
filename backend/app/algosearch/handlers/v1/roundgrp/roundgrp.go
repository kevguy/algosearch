package roundgrp

import (
	"context"
	"fmt"
	"github.com/kevguy/algosearch/backend/business/core/algod"
	"github.com/kevguy/algosearch/backend/business/core/block"
	"github.com/kevguy/algosearch/backend/business/core/block/db"
	v1web "github.com/kevguy/algosearch/backend/business/web/v1"
	"github.com/kevguy/algosearch/backend/foundation/web"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

type Handlers struct {
	BlockCore block.Core
	AlgodCore algod.Core
}

// GetCurrentRoundFromAPI retrieves the current round and returns the block data from Algod API
func (h Handlers) GetCurrentRoundFromAPI(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	blockData, err := h.AlgodCore.GetCurrentRound(ctx, v.TraceID)
	if err != nil {
		return errors.Wrap(err, "unable to get current round")
	}

	return web.Respond(ctx, w, blockData, http.StatusOK)
}

// GetRoundFromAPI retrieves a block from the Algod API based on the round number (num)
func (h Handlers) GetRoundFromAPI(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	numStr := web.Param(r, "num")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return v1web.NewRequestError(fmt.Errorf("invalid num format: %s", numStr), http.StatusBadRequest)
	}

	blockData, err := h.AlgodCore.GetRound(ctx, v.TraceID, uint64(num))
	if err != nil {
		return errors.Wrapf(err, "unable to get round %d", num)
	}

	return web.Respond(ctx, w, blockData, http.StatusOK)
}

// GetRound retrieves a block from CouchDB based on the round number (num)
func (h Handlers) GetRound(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	//v, err := web.GetValues(ctx)
	//if err != nil {
	//	return web.NewShutdownError("web value missing from context")
	//}

	numStr := web.Param(r, "num")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return v1web.NewRequestError(fmt.Errorf("invalid num format: %s", numStr), http.StatusBadRequest)
	}

	blockData, err := h.BlockCore.GetBlockByNum(ctx, uint64(num))
	if err != nil {
		return errors.Wrapf(err, "unable to get round %d", num)
	}

	return web.Respond(ctx, w, blockData, http.StatusOK)

}

// GetLatestSyncedRound the latest block from CouchDB.
func (h Handlers) GetLatestSyncedRound(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	//v, err := web.GetValues(ctx)
	//if err != nil {
	//	return web.NewShutdownError("web value missing from context")
	//}

	blockData, err := h.BlockCore.GetLatestBlock(ctx)
	if err != nil {
		return errors.Wrapf(err, "unable to get latest synced round")
	}

	return web.Respond(ctx, w, blockData, http.StatusOK)
}

// GetEarliestSyncedRound retrieves the earliest block from CouchDB.
func (h Handlers) GetEarliestSyncedRound(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	_, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	blockData, err := h.BlockCore.GetEarliestSyncedRoundNumber(ctx)
	if err != nil {
		return errors.Wrapf(err, "unable to get earliest synced round")
	}

	return web.Respond(ctx, w, blockData, http.StatusOK)
}

// GetRoundsPagination accepts the following parameters:
// - limit: number of items per page
// - latest_blk: the latest block number client wants to start with
// - page: number of pages
// - order: asc/desc
// The application counts from the latest_blk, calculates the number of pages using the number
// of items specified and retrieves the list of block for different pages.
// It returns the number of pages, number of blocks til the end and the list of blocks of interest
// as the response.
func (h Handlers) GetRoundsPagination(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	// limit
	limitQueries := web.Query(r, "limit")
	if len(limitQueries) == 0 {
		return v1web.NewRequestError(fmt.Errorf("missing query parameter: limit"), http.StatusBadRequest)
	}
	limit, err := strconv.Atoi(limitQueries[0])
	if err != nil {
		return v1web.NewRequestError(fmt.Errorf("invalid 'limit' format: %s", limitQueries[0]), http.StatusBadRequest)
	}

	// latest_blk
	latestBlkQueries := web.Query(r, "latest_blk")
	if len(latestBlkQueries) == 0 {
		return v1web.NewRequestError(fmt.Errorf("missing query parameter: latest_blk"), http.StatusBadRequest)
	}
	latestBlk, err := strconv.Atoi(latestBlkQueries[0])
	if err != nil {
		return v1web.NewRequestError(fmt.Errorf("invalid 'start' format: %s", latestBlkQueries[0]), http.StatusBadRequest)
	}

	// page
	pageQueries := web.Query(r, "page")
	if len(pageQueries) == 0 {
		return v1web.NewRequestError(fmt.Errorf("missing query parameter: page"), http.StatusBadRequest)
	}
	page, err := strconv.Atoi(pageQueries[0])
	if err != nil {
		return v1web.NewRequestError(fmt.Errorf("invalid 'page' format: %s", latestBlkQueries[0]), http.StatusBadRequest)
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
		return v1web.NewRequestError(fmt.Errorf("invalid 'sort' format: %s", orderQueries[0]), http.StatusBadRequest)
	}

	result, numOfPages, numOfBlks, err := h.BlockCore.GetBlocksPagination(ctx, int64(latestBlk), order, int64(page), int64(limit))
	if err != nil {
		return errors.Wrap(err, "Error fetching pagination results")
	}

	type Payload struct {
		NumOfPages	int64 `json:"num_of_pages"`
		NumOfBlks	int64   `json:"num_of_blks"`
		Items []db.Block `json:"items"`
	}

	return web.Respond(ctx, w, Payload{
		NumOfPages: numOfPages,
		NumOfBlks:  numOfBlks,
		Items:      result,
	}, http.StatusOK)
}

