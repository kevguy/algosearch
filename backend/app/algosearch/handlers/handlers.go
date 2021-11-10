// Package handlers contains the full set of handler functions and routes
// supported by the web api.
package handlers

import (
	"context"
	"expvar"
	block2 "github.com/kevguy/algosearch/backend/business/core/block"
	transaction2 "github.com/kevguy/algosearch/backend/business/core/transaction"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"github.com/go-kivik/kivik/v4"
	"github.com/kevguy/algosearch/backend/app/algosearch/handlers/apidoc/swaggergrp"
	"github.com/kevguy/algosearch/backend/app/algosearch/handlers/debug/samplegrp"
	"github.com/kevguy/algosearch/backend/app/algosearch/handlers/debug/checkgrp"
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
	sg, err := swaggergrp.NewIndex(cfg.APIProtocol, cfg.APIHost, "cal-engine-swagger")
	if err != nil {
		cfg.Log.Errorf("loading index template: %v", err)
		//return nil, errors.Wrap(err, "loading index template")
	}
	app.Handle(http.MethodGet, "", "/api/doc", sg.ServeDoc)

	// Register round endpoints
	rG := roundGroup{
		log:         cfg.Log,
		blockCore:	block2.NewCore(cfg.Log, cfg.CouchClient),
		algodClient: cfg.AlgodClient,
	}
	app.Handle(http.MethodGet, version, "/algod/current-round", rG.getCurrentRoundFromAPI)
	app.Handle(http.MethodGet, version, "/algod/rounds/:num", rG.getRoundFromAPI)
	app.Handle(http.MethodGet, version, "/current-round", rG.getLatestSyncedRound)
	app.Handle(http.MethodGet, version, "/earliest-round", rG.getEarliestSyncedRound)
	app.Handle(http.MethodGet, version, "/round/:num", rG.getRound)
	app.Handle(http.MethodGet, version, "/rounds", rG.getRoundsPagination)

	// Register transaction endpoints
	tG := transactionGroup{
		log:         cfg.Log,
		transactionCore: transaction2.NewCore(cfg.Log, cfg.CouchClient),
		algodClient: cfg.AlgodClient,
	}
	app.Handle(http.MethodGet, version, "/current-txn", tG.getLatestSyncedTransaction)
	app.Handle(http.MethodGet, version, "/earliest-txn", tG.getEarliestSyncedTransaction)
	app.Handle(http.MethodGet, version, "/transaction/:num", tG.getTransaction)
	app.Handle(http.MethodGet, version, "/transactions", tG.getTransactionsPagination)
}
