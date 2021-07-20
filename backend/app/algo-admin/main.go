// This program performs administrative tasks for the garage sale service.

package main

import (
	"fmt"
	"github.com/kevguy/algosearch/backend/app/algo-admin/commands"
	"github.com/kevguy/algosearch/backend/foundation/algorand"
	"github.com/kevguy/algosearch/backend/foundation/couchdb"
	"github.com/kevguy/algosearch/backend/foundation/logger"
	"os"
	"strconv"

	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// build is the git version of this program. It is set using build flags in the makefile.
var build = "develop"

func main() {

	// Construct the application logger.
	log := logger.New("ADMIN")
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
			Host       string `conf:"default:127.0.0.1:5984"`
		}
		Algorand struct {
			AlgodAddr	string `conf:"default:http://localhost:4001"`
			AlgodToken	string `conf:"default:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"`
			KmdAddr		string `conf:"default:http://localhost:7833"`
			KmdToken	string `conf:"default:a"`
		}
	}
	cfg.Version.SVN = build
	cfg.Version.Desc = "copyright information here"

	const prefix = "SALES"
	if err := conf.Parse(os.Args[1:], prefix, &cfg); err != nil {
		switch err {
		case conf.ErrHelpWanted:
			usage, err := conf.Usage(prefix, &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		case conf.ErrVersionWanted:
			version, err := conf.VersionString(prefix, &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config version")
			}
			fmt.Println(version)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
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

	algorandConfig := algorand.Config{
		AlgodAddr: cfg.Algorand.AlgodAddr,
		AlgodToken: cfg.Algorand.AlgodToken,
		KmdAddr: cfg.Algorand.KmdAddr,
		KmdToken: cfg.Algorand.KmdToken,
	}

	traceID := "00000000-0000-0000-0000-000000000000"

	switch cfg.Args.Num(0) {
	case "add-current-round":
		if err := commands.AddCurrentRoundCmd(traceID, log, algorandConfig, couchConfig); err != nil {
			return errors.Wrap(err, "add current round")
		}

	case "get-current-round":
		if err := commands.GetCurrentRoundCmd(algorandConfig); err != nil {
			return errors.Wrap(err, "getting current round")
		}

	case "add-round":
		numStr := cfg.Args.Num(1)
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return errors.Wrap(err, "num arg format wrong")
		}
		if err := commands.AddRoundCmd(traceID, log, algorandConfig, couchConfig, uint64(num)); err != nil {
			return errors.Wrap(err, "getting current round")
		}

	case "get-round":
		numStr := cfg.Args.Num(1)
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return errors.Wrap(err, "num arg format wrong")
		}
		if err := commands.GetRoundCmd(algorandConfig, uint64(num)); err != nil {
			return errors.Wrap(err, "getting current round")
		}

	case "get-round-from-db":
		numStr := cfg.Args.Num(1)
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return errors.Wrap(err, "num arg format wrong")
		}
		if err := commands.GetRoundInDbCmd(traceID, log, couchConfig, uint64(num)); err != nil {
			return errors.Wrap(err, "add round from db")
		}

	case "get-last-synced-round-num":
		if err := commands.GetLastSyncedRoundCmd(traceID, log, couchConfig); err != nil {
			return errors.Wrap(err, "add current round")
		}

	case "migrate":
		if err := commands.Migrate(couchConfig); err != nil {
			return errors.Wrap(err, "migrating database")
		}

	default:
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