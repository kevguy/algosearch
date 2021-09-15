package handlers

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	transaction2 "github.com/kevguy/algosearch/backend/business/data/store/transaction"
	"github.com/kevguy/algosearch/backend/business/sys/validate"
	"github.com/kevguy/algosearch/backend/foundation/web"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type transactionGroup struct {
	log         *zap.SugaredLogger
	store       transaction2.Store
	algodClient *algod.Client
}

// getTransaction retrieves a block from CouchDB based on the round number (num)
func (tG transactionGroup) getTransaction(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	_, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	id := web.Param(r, "id")
	// TODO: add trace ID
	transactionData, err := tG.store.GetTransaction(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "unable to get transaction %s", id)
	}

	return web.Respond(ctx, w, transactionData, http.StatusOK)
}

// getLatestSyncedTransaction retrieves the latest transaction from CouchDB.
func (tG transactionGroup) getLatestSyncedTransaction(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	_, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	transactionData, err := tG.store.GetLatestTransactionId(ctx)
	if err != nil {
		return errors.Wrapf(err, "unable to get latest synced transaction")
	}

	return web.Respond(ctx, w, transactionData, http.StatusOK)
}

// getEarliestSyncedTransaction retrieves the earliest transaction from CouchDB.
func (tG transactionGroup) getEarliestSyncedTransaction(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	_, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	transactionData, err := tG.store.GetEarliestTransactionId(ctx)
	if err != nil {
		return errors.Wrapf(err, "unable to get earliest synced transaction")
	}

	return web.Respond(ctx, w, transactionData, http.StatusOK)
}

func (tG transactionGroup) getTransactionsPagination(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	// limit
	limitQueries := web.Query(r, "limit")
	if len(limitQueries) == 0 {
		return validate.NewRequestError(fmt.Errorf("missing query parameter: limit"), http.StatusBadRequest)
	}
	limit, err := strconv.Atoi(limitQueries[0])
	if err != nil {
		return validate.NewRequestError(fmt.Errorf("invalid 'limit' format: %s", limitQueries[0]), http.StatusBadRequest)
	}

	// latest_txn
	latestTxnQueries := web.Query(r, "latest_txn")
	if len(latestTxnQueries) == 0 {
		return validate.NewRequestError(fmt.Errorf("missing query parameter: latest_txn"), http.StatusBadRequest)
	}
	latestTxn := latestTxnQueries[0]

	// page
	pageQueries := web.Query(r, "page")
	if len(pageQueries) == 0 {
		return validate.NewRequestError(fmt.Errorf("missing query parameter: page"), http.StatusBadRequest)
	}
	page, err := strconv.Atoi(pageQueries[0])
	if err != nil {
		return validate.NewRequestError(fmt.Errorf("invalid 'page' format: %s", latestTxnQueries[0]), http.StatusBadRequest)
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

	result, numOfPages, numOfTxns, err := tG.store.GetTransactionsPagination(ctx, latestTxn, order, int64(page), int64(limit))
	if err != nil {
		return errors.Wrap(err, "Error fetching pagination results")
	}

	type Payload struct {
		NumOfPages	int64 `json:"num_of_pages"`
		NumOfTxns	int64               `json:"num_of_txns"`
		Items []transaction2.Transaction `json:"items"`
	}

	return web.Respond(ctx, w, Payload{
		NumOfPages: numOfPages,
		NumOfTxns:  numOfTxns,
		Items:      result,
	}, http.StatusOK)
}
