package main

import (
	"context"
	"errors"
	"expvar" // Calls init function.
	"fmt"
	"github.com/kevguy/algosearch/backend/app/algosearch/blocksynchronizer"
	"github.com/kevguy/algosearch/backend/business/sys/auth"
	"github.com/kevguy/algosearch/backend/foundation/algod"
	"github.com/kevguy/algosearch/backend/foundation/couchdb"
	"github.com/kevguy/algosearch/backend/foundation/indexer"
	"github.com/kevguy/algosearch/backend/foundation/keystore"
	"github.com/kevguy/algosearch/backend/foundation/websocket"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	_indexer "github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"github.com/ardanlabs/conf/v2"
	"github.com/kevguy/algosearch/backend/app/algosearch/handlers"
	"github.com/kevguy/algosearch/backend/foundation/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
)

/*
Need to figure out timeouts for http service.
*/

// build is the git version of this program. It is set using build flags in the makefile.
var build = "develop"

func main() {

	// Construct the application logger.
	log, err := logger.New("ALGOSEARCH")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	// Perform the startup and shutdown sequence.
	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		log.Sync()
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {

	// =========================================================================
	// GOMAXPROCS

	// Set the correct number of threads for the service
	// based on what is available either by the machine or quotas.
	if _, err := maxprocs.Set(); err != nil {
		return fmt.Errorf("maxprocs: %w", err)
	}
	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// =========================================================================
	// Configuration

	cfg := struct {
		conf.Version
		Web struct {
			DeployProtocol	string		  `conf:"default:http,help:the protocol the deployment of this service will be using"`
			DeployHost		string		  `conf:"default:0.0.0.0:3000,help:the endpoint this service is deployed to"`
			APIHost         string        `conf:"default:0.0.0.0:3000"`
			DebugHost       string        `conf:"default:0.0.0.0:4000"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:120s"`
			ShutdownTimeout time.Duration `conf:"default:20s"`
			EnableSync		bool		  `conf:"default:false,help:specifies if the API should auto-sync new blocks"`
		}
		Auth struct {
			KeysFolder string `conf:"default:zarf/keys/"`
			ActiveKID  string `conf:"default:54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"`
		}
		CouchDB struct {
			Protocol   string `conf:"default:http"`
			User       string `conf:"default:admin"`
			Password   string `conf:"default:password,mask"`
			//Host       string `conf:"default:127.0.0.1:5984"`
			Host       string `conf:"default:89.39.110.254:5984"`
		}
		Algorand struct {
			//AlgodAddr		string `conf:"default:http://localhost:4001"`
			//AlgodToken		string `conf:"default:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"`
			//KmdAddr			string `conf:"default:http://localhost:7833"`
			//KmdToken		string `conf:"default:a"`
			//IndexerAddr 	string `conf:"default:http://localhost:8980"`
			//IndexerToken	string `conf:"default:empty"`
			AlgodAddr		string `conf:"default:http://89.39.110.254:4001"`
			AlgodToken		string `conf:"default:a2d2ac864300588718c6c05ff241a14fad99d30a19806356f3b9c8008559c4c1"`
			KmdAddr			string `conf:""`
			KmdToken		string `conf:""`
			IndexerAddr 	string `conf:""`
			IndexerToken	string `conf:"default:empty"`
		}
		Zipkin struct {
			ReporterURI string  `conf:"default:http://localhost:9411/api/v2/spans"`
			ServiceName string  `conf:"default:algosearch"`
			Probability float64 `conf:"default:0.05"`
		}
	}{
		Version: conf.Version{
			Build:  build,
			Desc: "copyright information here",
		},
	}

	const prefix = "ALGOSEARCH"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	// This is some special handling that the configuration library cannot
	// handle default value being an empty string
	// TODO: fix this
	if cfg.Algorand.IndexerToken == "empty" {
		cfg.Algorand.IndexerAddr = ""
	}

	// =========================================================================
	// App Starting

	log.Infow("starting service", "version", build)
	defer log.Infow("shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Infow("startup", "config", out)

	expvar.NewString("build").Set(build)

	// =========================================================================
	// Initialize authentication support

	log.Infow("startup", "status", "initializing authentication support")

	// Construct a key store based on the key files stored in
	// the specified directory.
	ks, err := keystore.NewFS(os.DirFS(cfg.Auth.KeysFolder))
	if err != nil {
		return fmt.Errorf("reading keys: %w", err)
	}

	auth, err := auth.New(cfg.Auth.ActiveKID, ks)
	if err != nil {
		return fmt.Errorf("constructing auth: %w", err)
	}


	// =========================================================================
	// Start Tracing Support

	log.Infow("startup", "status", "initializing OT/Zipkin tracing support")

	traceProvider, err := startTracing(
		cfg.Zipkin.ServiceName,
		cfg.Zipkin.ReporterURI,
		cfg.Zipkin.Probability,
	)
	if err != nil {
		return fmt.Errorf("starting tracing: %w", err)
	}
	defer traceProvider.Shutdown(context.Background())

	// =========================================================================
	// Start Algorand Algod Client

	log.Infow("startup", "status", "initializing algorand algod client support", "host", cfg.Algorand.AlgodAddr)

	algodClient, err := algod.Open(algod.Config{
		AlgodAddr: cfg.Algorand.AlgodAddr,
		AlgodToken: cfg.Algorand.AlgodToken,
		KmdAddr: cfg.Algorand.KmdAddr,
		KmdToken: cfg.Algorand.KmdToken,
	})
	if err != nil {
		return fmt.Errorf("connecting to algorand node: %w", err)
	}
	defer func() {
		log.Infow("shutdown", "status", "stopping algorand algod client support", "host", cfg.Algorand.AlgodAddr)
		//algodClient.Close()
	}()

	// =========================================================================
	// Start Algorand Indexer Client

	var indexerClient *_indexer.Client = nil
	if cfg.Algorand.IndexerAddr != "" {
		log.Infow("startup", "status", "initializing algorand indexer client support", "host", cfg.Algorand.AlgodAddr)

		indexerClient, err = indexer.Open(indexer.Config{
			IndexerAddr: cfg.Algorand.IndexerAddr,
			IndexerToken: cfg.Algorand.IndexerToken,
		})
		if err != nil {
			return fmt.Errorf("connecting to algorand indexer: %w", err)
		}
		defer func() {
			log.Infow("shutdown", "status", "stopping algorand indexer client support", "host", cfg.Algorand.IndexerAddr)
			//algodClient.Close()
		}()
	}

	// =========================================================================
	// Start CouchDB Client

	log.Infow("startup", "status", "initializing couchdb client support", "host", cfg.CouchDB.Host)

	couchConfig := couchdb.Config{
		Protocol: cfg.CouchDB.Protocol,
		User:     cfg.CouchDB.User,
		Password: cfg.CouchDB.Password,
		Host:     cfg.CouchDB.Host,
	}

	db, err := couchdb.Open(couchConfig)
	if err != nil {
		return fmt.Errorf("connecting to couchdb database: %w", err)
	}

	hub := websocket.NewHub()
	go hub.Run()

	// =========================================================================
	// Start Debug Service

	log.Infow("startup", "status", "debug v1 router started", "host", cfg.Web.DebugHost)

	// The Debug function returns a mux to listen and serve on for all the debug
	// related endpoints. This includes the standard library endpoints.

	// Construct the mux for the debug calls.
	debugMux := handlers.DebugMux(build, log, db, algodClient)

	// Start the service listening for debug requests.
	// Not concerned with shutting this down with load shedding.
	go func() {
		if err := http.ListenAndServe(cfg.Web.DebugHost, debugMux); err != nil {
			log.Errorw("shutdown", "status", "debug v1 router closed", "host", cfg.Web.DebugHost, "ERROR", err)
		}
	}()

	// =========================================================================
	// Start API Service

	log.Infow("startup", "status", "initializing V1 API support")

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Construct the mux for the API calls.
	apiMux := handlers.APIMux(handlers.APIMuxConfig{
		APIProtocol:   	cfg.Web.DeployProtocol,
		APIHost:       	cfg.Web.DeployHost,
		Shutdown:      	shutdown,
		Log:         	log,
		Auth:     		auth,
		AlgodClient: 	algodClient,
		IndexerClient: 	indexerClient,
		CouchClient: 	db,
		Hub: 			hub,
	})

	// Construct a server to service the requests against the mux.
	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      apiMux,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     zap.NewStdLog(log.Desugar()),
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)


	/*
	hub := websocket.NewHub()
	go hub.Run()
	http.HandleFunc("/wstest", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hey yo")
		fmt.Println("Hey yo")
		fmt.Println("Hey yo")
		fmt.Println("Hey yo")
		fmt.Println("Hey yo")
		fmt.Println("Hey yo")
		fmt.Println("Hey yo")
		fmt.Println("Hey yo")
		fmt.Println("Hey yo")
		websocket.ServeWs(hub, w, r)
	})
	 */

	// Start the service listening for api requests.
	go func() {
		log.Infow("startup", "status", "api router started", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()


	if cfg.Web.EnableSync {
		// Start the publisher to collect/publish metrics.
		blocksync, err := blocksynchronizer.New(log, 100*time.Millisecond, algodClient, couchConfig)
		if err != nil {
			return fmt.Errorf("starting publisher: %w", err)
		}
		defer blocksync.Stop()
	}

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shut down and shed load.
		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}

// =============================================================================

// startTracing configure open telemetry to be used with zipkin.
func startTracing(serviceName string, reporterURI string, probability float64) (*trace.TracerProvider, error) {

	// WARNING: The current settings are using defaults which may not be
	// compatible with your project. Please review the documentation for
	// opentelemetry.

	exporter, err := zipkin.New(
		reporterURI,
		// zipkin.WithLogger(zap.NewStdLog(log)),
	)
	if err != nil {
		return nil, fmt.Errorf("creating new exporter: %w", err)
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.TraceIDRatioBased(probability)),
		trace.WithBatcher(exporter,
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
			trace.WithBatchTimeout(trace.DefaultBatchTimeout),
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
		),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName),
				attribute.String("exporter", "zipkin"),
			),
		),
	)

	// I can only get this working properly using the singleton :(
	otel.SetTracerProvider(traceProvider)
	return traceProvider, nil
}
