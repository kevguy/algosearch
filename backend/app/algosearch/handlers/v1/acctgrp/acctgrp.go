package acctgrp

import (
	"context"
	"fmt"
	"github.com/kevguy/algosearch/backend/business/core/account"
	"github.com/kevguy/algosearch/backend/business/core/account/db"
	v1web "github.com/kevguy/algosearch/backend/business/web/v1"
	"github.com/kevguy/algosearch/backend/foundation/web"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

type Handlers struct {
	AcctCore account.Core
}

// GetAccount retrieves an account from CouchDB based on the account address (addr)
func (h Handlers) GetAccount(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	_, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	addr := web.Param(r, "addr")

	acctData, err := h.AcctCore.GetAccount(ctx, addr)
	if err != nil {
		return errors.Wrapf(err, "unable to get account %s", addr)
	}

	return web.Respond(ctx, w, acctData, http.StatusOK)
}

// GetLatestSyncedAccountAddr retrieves the latest account address from CouchDB.
func (h Handlers) GetLatestSyncedAccountAddr(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	_, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	accountData, err := h.AcctCore.GetLatestAccountID(ctx)
	if err != nil {
		return errors.Wrapf(err, "unable to get latest synced account address")
	}

	return web.Respond(ctx, w, accountData, http.StatusOK)
}

// GetEarliestSyncedAccountAddr retrieves the earliest account address from CouchDB.
func (h Handlers) GetEarliestSyncedAccountAddr(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	_, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	accountData, err := h.AcctCore.GetEarliestAccountID(ctx)
	if err != nil {
		return errors.Wrapf(err, "unable to get earliest synced account address")
	}

	return web.Respond(ctx, w, accountData, http.StatusOK)
}

func (h Handlers) GetAccountsPagination(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	// limit
	limitQueries := web.Query(r, "limit")
	if len(limitQueries) == 0 {
		return v1web.NewRequestError(fmt.Errorf("missing query parameter: limit"), http.StatusBadRequest)
	}
	limit, err := strconv.Atoi(limitQueries[0])
	if err != nil {
		return v1web.NewRequestError(fmt.Errorf("invalid 'limit' format: %s", limitQueries[0]), http.StatusBadRequest)
	}

	// latest_acct
	latestAcctQueries := web.Query(r, "latest_acct")
	if len(latestAcctQueries) == 0 {
		return v1web.NewRequestError(fmt.Errorf("missing query parameter: latest_acct"), http.StatusBadRequest)
	}
	latestAcctID := latestAcctQueries[0]

	// page
	pageQueries := web.Query(r, "page")
	if len(pageQueries) == 0 {
		return v1web.NewRequestError(fmt.Errorf("missing query parameter: page"), http.StatusBadRequest)
	}
	page, err := strconv.Atoi(pageQueries[0])
	if err != nil {
		//return v1web.NewRequestError(fmt.Errorf("invalid 'page' format: %s", latestTxnQueries[0]), http.StatusBadRequest)
		return v1web.NewRequestError(fmt.Errorf("invalid 'page' format: %s", pageQueries[0]), http.StatusBadRequest)
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

	result, numOfPages, numOfAccts, err := h.AcctCore.GetAccountsPagination(ctx, latestAcctID, order, int64(page), int64(limit))
	if err != nil {
		return fmt.Errorf("error fetching pagination results: %w", err)
	}

	type Payload struct {
		NumOfPages	int64 			`json:"num_of_pages"`
		NumOfAccts	int64 			`json:"num_of_accts"`
		Items 		[]db.Account 	`json:"items"`
	}

	return web.Respond(ctx, w, Payload{
		NumOfPages: numOfPages,
		NumOfAccts:  numOfAccts,
		Items:      result,
	}, http.StatusOK)
}

func (h Handlers) GetAcctCount(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	earliestAcctID, err := h.AcctCore.GetEarliestAccountID(ctx)
	if err != nil {
		return fmt.Errorf("fetching earliest account id overall in couch: %w", err)
	}

	latestAcctID, err := h.AcctCore.GetLatestAccountID(ctx)
	if err != nil {
		return fmt.Errorf("fetching latest account id overall in couch: %w", err)
	}

	count, err := h.AcctCore.GetAccountCountBtnKeys(ctx,
		earliestAcctID,
		latestAcctID)
	if err != nil {
		return fmt.Errorf("error fetching account count: %w", err)
	}

	return web.Respond(ctx, w, count, http.StatusOK)
}
