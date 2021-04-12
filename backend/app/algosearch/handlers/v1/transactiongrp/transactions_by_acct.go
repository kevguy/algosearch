package transactiongrp

import (
	"context"
	"fmt"
	"github.com/kevguy/algosearch/backend/business/core/transaction/db"
	v1web "github.com/kevguy/algosearch/backend/business/web/v1"
	"github.com/kevguy/algosearch/backend/foundation/web"
	"net/http"
	"strconv"
)

func (h Handlers) GetTransactionsByAcctIDCount(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	acctID := web.Param(r, "acct_id")

	earliestTransaction, err := h.TransactionCore.GetEarliestTransaction(ctx)
	if err != nil {
		return fmt.Errorf("fetching earliest transaction overall in couch: %w", err)
	}

	latestTransaction, err := h.TransactionCore.GetLatestTransaction(ctx)
	if err != nil {
		return fmt.Errorf("fetching latest transaction overall in couch: %w", err)
	}

	count, err := h.TransactionCore.GetTransactionCountByAcct(ctx,
		acctID,
		earliestTransaction.ID,
		latestTransaction.ID)
	if err != nil {
		return fmt.Errorf("error fetching transaction count by acct ID[%q]: %w", acctID, err)
	}

	return web.Respond(ctx, w, count, http.StatusOK)
}


func (h Handlers) GetTransactionsByAcctID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	acctID := web.Param(r, "acct_id")

	// limit
	limitQueries := web.Query(r, "limit")
	if len(limitQueries) == 0 {
		return v1web.NewRequestError(fmt.Errorf("missing query parameter: limit"), http.StatusBadRequest)
	}
	limit, err := strconv.Atoi(limitQueries[0])
	if err != nil {
		return v1web.NewRequestError(fmt.Errorf("invalid 'limit' format: %s", limitQueries[0]), http.StatusBadRequest)
	}

	// page
	pageNoQueries := web.Query(r, "page")
	if len(pageNoQueries) == 0 {
		return v1web.NewRequestError(fmt.Errorf("missing query parameter: page"), http.StatusBadRequest)
	}
	pageNo, err := strconv.Atoi(pageNoQueries[0])
	if err != nil {
		//return v1web.NewRequestError(fmt.Errorf("invalid 'page' format: %s", latestTxnQueries[0]), http.StatusBadRequest)
		return v1web.NewRequestError(fmt.Errorf("invalid 'page' format: %s", pageNoQueries[0]), http.StatusBadRequest)
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

	result, numOfPages, numOfTxns, err := h.TransactionCore.GetTransactionsByAcctPagination(ctx, acctID, order, int64(pageNo), int64(limit))
	if err != nil {
		return fmt.Errorf("error fetching pagination results: %w", err)
	}

	type Payload struct {
		NumOfPages	int64 `json:"num_of_pages"`
		NumOfTxns	int64 `json:"num_of_txns"`
		Items []db.Transaction `json:"items"`
	}

	return web.Respond(ctx, w, Payload{
		NumOfPages: numOfPages,
		NumOfTxns:  numOfTxns,
		Items:      result,
	}, http.StatusOK)
}
