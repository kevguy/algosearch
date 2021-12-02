package transactiongrp

import (
	"context"
	"fmt"
	"github.com/kevguy/algosearch/backend/business/core/transaction"
	"github.com/kevguy/algosearch/backend/business/core/transaction/db"
	v1web "github.com/kevguy/algosearch/backend/business/web/v1"
	"github.com/kevguy/algosearch/backend/foundation/web"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

type Handlers struct {
	TransactionCore transaction.Core
}

// GetTransaction retrieves a block from CouchDB based on the round number (num)
func (h Handlers) GetTransaction(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	_, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	id := web.Param(r, "id")
	// TODO: add trace ID
	transactionData, err := h.TransactionCore.GetTransaction(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "unable to get transaction %s", id)
	}

	return web.Respond(ctx, w, transactionData, http.StatusOK)
}

// GetLatestSyncedTransaction retrieves the latest transaction from CouchDB.
func (h Handlers) GetLatestSyncedTransaction(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	_, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	transactionData, err := h.TransactionCore.GetLatestTransaction(ctx)
	if err != nil {
		return errors.Wrapf(err, "unable to get latest synced transaction")
	}

	return web.Respond(ctx, w, transactionData, http.StatusOK)
}

// GetEarliestSyncedTransaction retrieves the earliest transaction from CouchDB.
func (h Handlers) GetEarliestSyncedTransaction(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	_, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	transactionData, err := h.TransactionCore.GetEarliestTransaction(ctx)
	if err != nil {
		return errors.Wrapf(err, "unable to get earliest synced transaction")
	}

	return web.Respond(ctx, w, transactionData, http.StatusOK)
}

func (h Handlers) GetTransactionsPagination(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	// limit
	limitQueries := web.Query(r, "limit")
	if len(limitQueries) == 0 {
		return v1web.NewRequestError(fmt.Errorf("missing query parameter: limit"), http.StatusBadRequest)
	}
	limit, err := strconv.Atoi(limitQueries[0])
	if err != nil {
		return v1web.NewRequestError(fmt.Errorf("invalid 'limit' format: %s", limitQueries[0]), http.StatusBadRequest)
	}

	// latest_txn
	latestTxnQueries := web.Query(r, "latest_txn")
	if len(latestTxnQueries) == 0 {
		return v1web.NewRequestError(fmt.Errorf("missing query parameter: latest_txn"), http.StatusBadRequest)
	}
	latestTxn := latestTxnQueries[0]

	// page
	pageQueries := web.Query(r, "page")
	if len(pageQueries) == 0 {
		return v1web.NewRequestError(fmt.Errorf("missing query parameter: page"), http.StatusBadRequest)
	}
	page, err := strconv.Atoi(pageQueries[0])
	if err != nil {
		return v1web.NewRequestError(fmt.Errorf("invalid 'page' format: %s", latestTxnQueries[0]), http.StatusBadRequest)
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

	result, numOfPages, numOfTxns, err := h.TransactionCore.GetTransactionsPagination(ctx, latestTxn, order, int64(page), int64(limit))
	if err != nil {
		return errors.Wrap(err, "Error fetching pagination results")
	}

	type Payload struct {
		NumOfPages	int64 `json:"num_of_pages"`
		NumOfTxns	int64               `json:"num_of_txns"`
		Items []db.Transaction `json:"items"`
	}

	return web.Respond(ctx, w, Payload{
		NumOfPages: numOfPages,
		NumOfTxns:  numOfTxns,
		Items:      result,
	}, http.StatusOK)
}
