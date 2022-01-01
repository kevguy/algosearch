package srchgrp

import (
	"context"
	"fmt"
	"github.com/kevguy/algosearch/backend/business/core/account"
	"github.com/kevguy/algosearch/backend/business/core/application"
	"github.com/kevguy/algosearch/backend/business/core/asset"
	"github.com/kevguy/algosearch/backend/business/core/block"
	"github.com/kevguy/algosearch/backend/business/core/transaction"
	v1web "github.com/kevguy/algosearch/backend/business/web/v1"
	"github.com/kevguy/algosearch/backend/foundation/web"
	"net/http"
	"strconv"
)

type Handlers struct {
	BlockCore 		block.Core
	TransactionCore transaction.Core
	AcctCore 		account.Core
	AssetCore 		asset.Core
	ApplicationCore	application.Core
}

func (h Handlers) SrchKey(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	_, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}
	keyQueries := web.Query(r, "key")
	if len(keyQueries) == 0 {
		return v1web.NewRequestError(fmt.Errorf("missing query parameter: limit"), http.StatusBadRequest)
	}

	var blockHashFound = false
	var blockRoundFound = false
	var txnFound = false
	var acctFound = false
	var assetFound = false
	var appFound = false

	// Search if block exists
	block, err := h.BlockCore.GetBlockByHash(ctx, keyQueries[0])
	if block.BlockHash != "" && block.BlockHash == keyQueries[0] {
		blockHashFound = true
	}
	if keyQueries[0] != "" {
		keyInInt, err := strconv.Atoi(keyQueries[0])
		if err == nil {
			block, err := h.BlockCore.GetBlockByNum(ctx, uint64(keyInInt))
			if err == nil && block.Round == uint64(keyInInt) {
				blockRoundFound = true
			}
		}
	}

	// Search if transaction address exists
	txn, err := h.TransactionCore.GetTransaction(ctx, keyQueries[0])
	if txn.Id != "" && txn.Id == keyQueries[0] {
		txnFound = true
	}

	// Search if account address exists
	acct, err := h.AcctCore.GetAccount(ctx, keyQueries[0])
	if acct.Address != "" && acct.Address == acct.Address {
		acctFound = true
	}

	// Search if asset exists
	asset, err := h.AssetCore.GetAsset(ctx, keyQueries[0])
	if asset.Index != 0 {
		assetFound = true
	}

	// Search if application exists
	if keyQueries[0] != "" {
		keyInInt, err := strconv.Atoi(keyQueries[0])
		if err == nil {
			app, err := h.ApplicationCore.GetApplication(ctx, fmt.Sprintf("%d", keyInInt))
			if err == nil && app.Id == uint64(keyInInt) {
				appFound = true
			}
		}
	}

	// TODO: Search by Asset Name

	// TODO: Search Group Tx ID

	type Response struct {
		BlockHashFound bool `json:"block_hash_found"`
		BlockRoundFound bool `json:"block_round_found"`
		TxnFound bool `json:"txn_found"`
		AcctFound bool `json:"acct_found"`
		AssetFound bool `json:"asset_found"`
		AppFound bool `json:"application_found"`
	}

	return web.Respond(ctx, w, Response{
		BlockHashFound:  blockHashFound,
		BlockRoundFound: blockRoundFound,
		TxnFound:        txnFound,
		AcctFound:       acctFound,
		AssetFound:      assetFound,
		AppFound:        appFound,
	}, http.StatusOK)
}
