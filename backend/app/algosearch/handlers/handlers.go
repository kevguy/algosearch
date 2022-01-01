// Package handlers contains the full set of handler functions and routes
// supported by the web api.
package handlers

import (
	"context"
	"expvar"
	"github.com/kevguy/algosearch/backend/app/algosearch/handlers/v1/acctgrp"
	"github.com/kevguy/algosearch/backend/app/algosearch/handlers/v1/assetgrp"
	"github.com/kevguy/algosearch/backend/app/algosearch/handlers/v1/ledgergrp"
	"github.com/kevguy/algosearch/backend/app/algosearch/handlers/v1/roundgrp"
	"github.com/kevguy/algosearch/backend/app/algosearch/handlers/v1/srchgrp"
	"github.com/kevguy/algosearch/backend/app/algosearch/handlers/v1/transactiongrp"
	"github.com/kevguy/algosearch/backend/app/algosearch/handlers/v1/wsgrp"
	"github.com/kevguy/algosearch/backend/business/core/account"
	algod2 "github.com/kevguy/algosearch/backend/business/core/algod"
	"github.com/kevguy/algosearch/backend/business/core/application"
	"github.com/kevguy/algosearch/backend/business/core/asset"
	block2 "github.com/kevguy/algosearch/backend/business/core/block"
	transaction2 "github.com/kevguy/algosearch/backend/business/core/transaction"
	"github.com/kevguy/algosearch/backend/foundation/websocket"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"github.com/go-kivik/kivik/v4"
	"github.com/kevguy/algosearch/backend/app/algosearch/handlers/apidoc/swaggergrp"
	"github.com/kevguy/algosearch/backend/app/algosearch/handlers/debug/checkgrp"
	"github.com/kevguy/algosearch/backend/app/algosearch/handlers/debug/samplegrp"
	"github.com/kevguy/algosearch/backend/business/sys/auth"

	"github.com/kevguy/algosearch/backend/business/web/v1/mid"
	"github.com/kevguy/algosearch/backend/foundation/web"
	"go.uber.org/zap"
)

// Options represent optional parameters.
type Options struct {
	corsOrigin string
}

// WithCORS provides configuration options for CORS.
func WithCORS(origin string) func(opts *Options) {
	return func(opts *Options) {
		opts.corsOrigin = origin
	}
}


// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Shutdown chan os.Signal
	APIProtocol		string
	APIHost			string
	Log				*zap.SugaredLogger
	Auth			*auth.Auth
	AlgodClient		*algod.Client
	IndexerClient	*indexer.Client
	CouchClient		*kivik.Client
	Hub				*websocket.Hub
}

// APIMux constructs an http.Handler with all application routes defined.
func APIMux(cfg APIMuxConfig, options ...func(opts *Options)) http.Handler {
	var opts Options
	for _, option := range options {
		option(&opts)
	}

	// Construct the web.App which holds all routes as well as common Middleware.
	app := web.NewApp(
		cfg.Shutdown,
		mid.Logger(cfg.Log),
		mid.Errors(cfg.Log),
		mid.Metrics(),
		mid.Panics(),
	)

	// Register the swagger assets.
	// For the endpoint /swagger/*,
	// files will be served inside the swagger folder
	fs := http.FileServer(http.Dir("swagger"))
	fs = http.StripPrefix("/swagger/", fs)
	f := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		fs.ServeHTTP(w, r)
		return nil
	}
	app.Handle(http.MethodGet, "", "/swagger/*", f, mid.Cors("*"))



	// Accept CORS 'OPTIONS' preflight requests if config has been provided.
	// Don't forget to apply the CORS middleware to the routes that need it.
	// Example Config: `conf:"default:https://MY_DOMAIN.COM"`
	if opts.corsOrigin != "" {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			return nil
		}
		app.Handle(http.MethodOptions, "", "/*", h)
	}

	// Load the routes for the different versions of the API.
	v1(app, cfg)

	return app
}

// DebugStandardLibraryMux registers all the debug routes from the standard library
// into a new mux bypassing the use of the DefaultServerMux. Using the
// DefaultServerMux would be a security risk since a dependency could inject a
// handler into our service without us knowing it.
func DebugStandardLibraryMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Register all the standard library debug endpoints.
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}

// DebugMux registers all the debug standard library routes and then custom
// debug application routes for the service. This bypassing the use of the
// DefaultServerMux. Using the DefaultServerMux would be a security risk since
// a dependency could inject a handler into our service without us knowing it.
func DebugMux(build string, log *zap.SugaredLogger, couchClient *kivik.Client, algodClient *algod.Client) http.Handler {
	mux := DebugStandardLibraryMux()

	// Register debug check endpoints.
	cgh := checkgrp.Handlers{
		Build: build,
		Log:   log,
		CouchClient: couchClient,
		AlgodClient: algodClient,
	}
	mux.HandleFunc("/debug/readiness", cgh.Readiness)
	mux.HandleFunc("/debug/liveness", cgh.Liveness)

	return mux
}

// v1 binds all the version 1 routes.
func v1(app *web.App, cfg APIMuxConfig) {
	const version = "v1"

	// Register sample endpoints
	samg := samplegrp.Handlers{}
	app.Handle(http.MethodGet, "", "/", samg.SendOK)
	app.Handle(http.MethodGet, "", "/test", samg.SendOK)
	app.Handle(http.MethodGet, "", "/test-error", samg.SendError)

	// Register the index page for the website.
	sg, err := swaggergrp.NewIndex(cfg.APIProtocol, cfg.APIHost, "algosearch")
	if err != nil {
		cfg.Log.Errorf("loading index template: %v", err)
		//return nil, errors.Wrap(err, "loading index template")
	}
	app.Handle(http.MethodGet, "", "/api/doc", sg.ServeDoc)

	algodCore := algod2.NewCore(cfg.Log, cfg.AlgodClient)
	blockCore := block2.NewCore(cfg.Log, cfg.CouchClient)
	txnCore := transaction2.NewCore(cfg.Log, cfg.CouchClient)
	acctCore := account.NewCore(cfg.Log, cfg.CouchClient)
	assetCore := asset.NewCore(cfg.Log, cfg.CouchClient)
	appCore := application.NewCore(cfg.Log, cfg.CouchClient)

// Register round endpoints
	rG := roundgrp.Handlers{
		BlockCore: blockCore,
		AlgodCore: algodCore,
	}
	app.Handle(http.MethodGet, version, "/algod/current-round", rG.GetCurrentRoundFromAPI, mid.Cors("*"))
	app.Handle(http.MethodGet, version, "/algod/rounds/:num", rG.GetRoundFromAPI, mid.Cors("*"))
	app.Handle(http.MethodGet, version, "/current-round", rG.GetLatestSyncedRound, mid.Cors("*"))
	app.Handle(http.MethodGet, version, "/earliest-round-num", rG.GetEarliestSyncedRound, mid.Cors("*"))
	app.Handle(http.MethodGet, version, "/rounds/:num", rG.GetRound, mid.Cors("*"))
	app.Handle(http.MethodGet, version, "/rounds", rG.GetRoundsPagination, mid.Cors("*"))

	// Register transaction endpoints
	tG := transactiongrp.Handlers{
		TransactionCore: txnCore,
	}
	app.Handle(http.MethodGet, version, "/current-txn", tG.GetLatestSyncedTransaction, mid.Cors("*"))
	app.Handle(http.MethodGet, version, "/earliest-txn", tG.GetEarliestSyncedTransaction, mid.Cors("*"))
	app.Handle(http.MethodGet, version, "/transactions/:id", tG.GetTransaction, mid.Cors("*"))
	app.Handle(http.MethodGet, version, "/transactions/acct/:acct_id", tG.GetTransactionsByAcctID, mid.Cors("*"))
	app.Handle(http.MethodGet, version, "/transactions/acct/:acct_id/count", tG.GetTransactionsByAcctIDCount, mid.Cors("*"))
	app.Handle(http.MethodGet, version, "/transactions", tG.GetTransactionsPagination, mid.Cors("*"))

	// Register account endpoints
	aG := acctgrp.Handlers{
		AcctCore: acctCore,
	}
	app.Handle(http.MethodGet, version, "/accounts/latest", aG.GetLatestSyncedAccountAddr, mid.Cors("*"))
	app.Handle(http.MethodGet, version, "/accounts/earliest", aG.GetEarliestSyncedAccountAddr, mid.Cors("*"))
	app.Handle(http.MethodGet, version, "/accounts/count", aG.GetAcctCount, mid.Cors("*"))
	app.Handle(http.MethodGet, version, "/accounts/:addr", aG.GetAccount, mid.Cors("*"))
	app.Handle(http.MethodGet, version, "/accounts", aG.GetAccountsPagination, mid.Cors("*"))

	asG := assetgrp.Handlers{AlgodCore: algodCore}
	app.Handle(http.MethodGet, version, "/algod/assets/:idx", asG.GetAssetByIDFromAPI, mid.Cors("*"))

	lG := ledgergrp.Handlers{
		AlgodCore: algodCore,
	}
	app.Handle(http.MethodGet, version, "/algod/ledger/supply", lG.GetLedgerSupplyFromAPI, mid.Cors("*"))

	sG := srchgrp.Handlers{
		BlockCore: blockCore,
		TransactionCore: txnCore,
		AcctCore: acctCore,
		AssetCore: assetCore,
		ApplicationCore: appCore,
	}
	app.Handle(http.MethodGet, version, "/search/:key", sG.SrchKey, mid.Cors("*"))

	// Register websocket endpoints
	wsG := wsgrp.Handlers{
		Hub: cfg.Hub,
	}
	app.Handle(http.MethodGet, "", "/wstest", wsG.ServeHomePage)
	app.Handle(http.MethodGet, "", "/ws", wsG.ServeWS)
	app.Handle(http.MethodGet, version, "/test-socket", wsG.SendDummy)
}
