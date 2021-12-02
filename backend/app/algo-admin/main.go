// This program performs administrative tasks for the garage sale service.

package main

import (
	"fmt"
	"github.com/kevguy/algosearch/backend/app/algo-admin/commands"
	"github.com/kevguy/algosearch/backend/foundation/algod"
	"github.com/kevguy/algosearch/backend/foundation/couchdb"
	"github.com/kevguy/algosearch/backend/foundation/indexer"
	"github.com/kevguy/algosearch/backend/foundation/logger"
	"os"
	"strconv"

	"github.com/ardanlabs/conf/v2"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// build is the git version of this program. It is set using build flags in the makefile.
var build = "develop"

func main() {

	// Construct the application logger.
	log, err := logger.New("ADMIN")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	if err := run(log); err != nil {
		if errors.Cause(err) != commands.ErrHelp {
			log.Errorw("", zap.Error(err))
		}
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {

	// =========================================================================
	// Configuration

	var cfg struct {
		conf.Version
		Args conf.Args
		CouchDB struct {
			Protocol   string `conf:"default:http"`
			User       string `conf:"default:admin"`
			Password   string `conf:"default:password,mask"`
			//Host       string `conf:"default:127.0.0.1:5984"`
			Host		string `conf:"default:89.39.110.254:5984"`
		}
		Algorand struct {
			//AlgodAddr	string `conf:"default:http://localhost:4001"`
			//AlgodToken	string `conf:"default:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"`
			//KmdAddr		string `conf:"default:http://localhost:7833"`
			//KmdToken	string `conf:"default:a"`
			//IndexerAddr 	string `conf:"default:http://localhost:8980"`
			//IndexerToken	string `conf:"default:empty"`
			AlgodAddr		string `conf:"default:http://89.39.110.254:4001"`
			AlgodToken		string `conf:"default:a2d2ac864300588718c6c05ff241a14fad99d30a19806356f3b9c8008559c4c1"`
			KmdAddr			string `conf:""`
			KmdToken		string `conf:""`
			IndexerAddr 	string `conf:""`
			IndexerToken	string `conf:"default:empty"`
		}
	}
	cfg.Version.Build = build
	cfg.Version.Desc = "copyright information here"


	// This is some special handling that the configuration library cannot
	// handle default value being an empty string
	// TODO: fix this
	if cfg.Algorand.IndexerToken == "empty" {
		cfg.Algorand.IndexerAddr = ""
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

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Infow("startup", "config", out)

	// =========================================================================
	// Commands

	couchConfig := couchdb.Config{
		Protocol: cfg.CouchDB.Protocol,
		User:     cfg.CouchDB.User,
		Password: cfg.CouchDB.Password,
		Host:     cfg.CouchDB.Host,
	}

	algorandConfig := algod.Config{
		AlgodAddr: cfg.Algorand.AlgodAddr,
		AlgodToken: cfg.Algorand.AlgodToken,
		KmdAddr: cfg.Algorand.KmdAddr,
		KmdToken: cfg.Algorand.KmdToken,
	}

	indexerConfig := indexer.Config{
		IndexerAddr:  cfg.Algorand.IndexerAddr,
		IndexerToken: cfg.Algorand.IndexerToken,
	}

	traceID := "00000000-0000-0000-0000-000000000000"

	switch cfg.Args.Num(0) {
	case "pprint-round-algod":
		numStr := cfg.Args.Num(1)
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return fmt.Errorf("num arg format wrong: %w", err)
		}
		if err := commands.PrettyPrintBlockFromAlgodCmd(traceID, log, algorandConfig, uint64(num)); err != nil {
			return errors.Wrap(err, "pretty print block from algod")
		}

	case "pprint-round-indexer":
		numStr := cfg.Args.Num(1)
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return fmt.Errorf("num arg format wrong: %w", err)
		}
		if err := commands.PrettyPrintBlockFromIndexerCmd(traceID, log, indexerConfig, uint64(num)); err != nil {
			return fmt.Errorf("pretty print block from indexer: : %w", err)
		}

	case "compare-round-algod-indexer":
		numStr := cfg.Args.Num(1)
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return fmt.Errorf("num arg format wrong: %w", err)
		}
		if err := commands.CompareBlockBetweenAlgodAndIndexer(traceID, log, algorandConfig, indexerConfig, uint64(num)); err != nil {
			return fmt.Errorf("comparing block bytes: %w", err)
		}

	case "add-current-round":
		if err := commands.AddCurrentRoundCmd(traceID, log, algorandConfig, couchConfig); err != nil {
			return fmt.Errorf("add current round: %w", err)
		}

	case "get-current-round":
		if err := commands.GetCurrentRoundCmd(log, algorandConfig); err != nil {
			return fmt.Errorf("getting current round: %w", err)
		}

	case "add-round":
		numStr := cfg.Args.Num(1)
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return fmt.Errorf("num arg format wrong: %w", err)
		}
		if err := commands.AddRoundCmd(traceID, log, algorandConfig, couchConfig, uint64(num)); err != nil {
			return fmt.Errorf("getting current round: %w", err)
		}

	case "get-round":
		numStr := cfg.Args.Num(1)
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return fmt.Errorf("num arg format wrong: %w", err)
		}
		if err := commands.GetRoundCmd(log, algorandConfig, uint64(num)); err != nil {
			return fmt.Errorf("getting current round: %w", err)
		}

	case "get-round-from-db":
		blockHashStr := cfg.Args.Num(1)
		if err != nil {
			return fmt.Errorf("num arg format wrong: %w", err)
		}
		if err := commands.GetRoundInDbCmd(traceID, log, couchConfig, blockHashStr); err != nil {
			return fmt.Errorf("add round from db: %w", err)
		}

	case "get-round-from-db-by-num":
		numStr := cfg.Args.Num(1)
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return fmt.Errorf("num arg format wrong: %w", err)
		}
		if err := commands.GetRoundNumInDbCmd(traceID, log, couchConfig, uint64(num)); err != nil {
			return fmt.Errorf("num arg format wrong: %w", err)
		}

	case "get-last-synced-round-num":
		if err := commands.GetLastSyncedRoundCmd(traceID, log, couchConfig); err != nil {
			return fmt.Errorf("add current round: %w", err)
		}

	case "get-rounds-pagination":
		latestBlockNumStr := cfg.Args.Num(1)
		latestBlockNum, err := strconv.Atoi(latestBlockNumStr)
		if err != nil {
			return fmt.Errorf("latestBlockNum arg format wrong: %w", err)
		}
		noOfItemsStr := cfg.Args.Num(2)
		noOfItems, err := strconv.Atoi(noOfItemsStr)
		if err != nil {
			return fmt.Errorf("noOfItems arg format wrong: %w", err)
		}
		pageNoStr := cfg.Args.Num(3)
		pageNo, err := strconv.Atoi(pageNoStr)
		if err != nil {
			return fmt.Errorf("pageNo arg format wrong: %w", err)
		}
		order := cfg.Args.Num(4)
		if err := commands.GetRoundsPaginationCmd(traceID, log, couchConfig, int64(latestBlockNum), int64(noOfItems), int64(pageNo), order); err != nil {
			return fmt.Errorf("add round from db: %w", err)
		}

	case "get-txns-pagination":
		latestTxnId := cfg.Args.Num(1)
		noOfItemsStr := cfg.Args.Num(2)
		noOfItems, err := strconv.Atoi(noOfItemsStr)
		if err != nil {
			return fmt.Errorf("noOfItems arg format wrong: %w", err)
		}
		pageNoStr := cfg.Args.Num(3)
		pageNo, err := strconv.Atoi(pageNoStr)
		if err != nil {
			return fmt.Errorf("pageNo arg format wrong: %w", err)
		}
		order := cfg.Args.Num(4)
		if err := commands.GetTransactionsPaginationCmd(traceID, log, couchConfig, latestTxnId, int64(noOfItems), int64(pageNo), order); err != nil {
			return fmt.Errorf("add round from db: %w", err)
		}

	case "get-and-insert-blocks":
		startBlockStr := cfg.Args.Num(1)
		startBlock, err := strconv.Atoi(startBlockStr)
		endBlockStr := cfg.Args.Num(2)
		endBlock, err := strconv.Atoi(endBlockStr)

		if err != nil {
			return fmt.Errorf("num arg format wrong: %w", err)
		}
		if err := commands.GetAndInsertBlockCmd(log, algorandConfig, couchConfig, uint64(startBlock), uint64(endBlock)); err != nil {
			return fmt.Errorf("num arg format wrong: %w", err)
		}

	case "get-blocks-count-from-db":
		if err := commands.GetBlocksCountFromDbCmd(log, couchConfig); err != nil {
			return fmt.Errorf("get number of blocks from db: %w", err)
		}

	case "get-txn-info-from-db":
		if err := commands.GetTransactionInfoFromDbCmd(log, couchConfig); err != nil {
			return fmt.Errorf("get transactions info from db: %w", err)
		}

	case "get-txns-from-db":
		limitStr := cfg.Args.Num(1)
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return fmt.Errorf("limit arg format wrong: %w", err)
		}
		pageNoStr := cfg.Args.Num(2)
		pageNo, err := strconv.Atoi(pageNoStr)
		if err != nil {
			return fmt.Errorf("pageNo arg format wrong: %w", err)
		}
		order := cfg.Args.Num(3)
		if order != "asc" && order != "desc" {
			return fmt.Errorf("order arg format wrong")
		}

		if err := commands.GetTransactionsFromDbWithPaginationCmd(log, couchConfig, int64(pageNo), int64(limit), order); err != nil {
			return fmt.Errorf("get transactions data from db %w", err)
		}

	case "get-txns-by-acct-from-db":
		acctID := cfg.Args.Num(1)
		if acctID == "" {
			return fmt.Errorf("acctID should not be empty")
		}
		if err := commands.GetTransactionsByAcctFromDbCmd(log, couchConfig, acctID); err != nil {
			return fmt.Errorf("get transactions info by account %s from db: %w", acctID, err)
		}

	case "get-txns-by-acct-pagination-from-db":
		acctID := cfg.Args.Num(1)
		if acctID == "" {
			return fmt.Errorf("acctID should not be empty")
		}
		limitStr := cfg.Args.Num(2)
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return fmt.Errorf("limit arg format wrong: %w", err)
		}
		pageNoStr := cfg.Args.Num(3)
		pageNo, err := strconv.Atoi(pageNoStr)
		if err != nil {
			return fmt.Errorf("pageNo arg format wrong: %w", err)
		}
		order := cfg.Args.Num(4)
		if order != "asc" && order != "desc" {
			return fmt.Errorf("order arg format wrong")
		}

		if err := commands.GetTransactionsByAcctFromDbWithPaginationCmd(log, couchConfig, acctID, int64(pageNo), int64(limit), order); err != nil {
			return fmt.Errorf("get transactions data from db %w", err)
		}

	case "migrate":
		if err := commands.Migrate(couchConfig); err != nil {
			return fmt.Errorf("migrating database: %w", err)
		}

	default:
		fmt.Println("get-txn-info-from-db: get general information about the transactions from the database")
		fmt.Println("add-current-round: add the current round to the database")
		fmt.Println("add-round: add a round to the database")
		fmt.Println("get-current-round: get the current round and print it nicely")
		fmt.Println("get-round: get a round and print it nicely")
		fmt.Println("get-round-from-db: get a round from the database")
		fmt.Println("get-last-synced-round-num: get the round number of the last block synced to the database")
		fmt.Println("migrate: create the schema in the CouchDB database")
		fmt.Println("provide a command to get more help.")
		return commands.ErrHelp
	}

	return nil
}
