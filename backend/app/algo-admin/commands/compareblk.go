package commands

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	algod2 "github.com/kevguy/algosearch/backend/business/core/algod"
	indexerApp "github.com/kevguy/algosearch/backend/business/core/indexer"
	"github.com/kevguy/algosearch/backend/foundation/algod"
	"github.com/kevguy/algosearch/backend/foundation/indexer"
	"github.com/nsf/jsondiff"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

// CompareBlockBetweenAlgodAndIndexer retrieves information about the block for the specified round and prints it out in JSON format.
func CompareBlockBetweenAlgodAndIndexer(traceID string, log *zap.SugaredLogger, algodCfg algod.Config, indexerCfg indexer.Config, blockNum uint64) error {
	algodClient, err := algod.Open(algodCfg)
	if err != nil {
		return errors.Wrap(err, "connect to Algorand Node")
	}

	algodCore := algod2.NewCore(log, algodClient)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	algodBlock, err := algodCore.GetRound(ctx, traceID, blockNum)
	//rawBlock, err := app.GetRoundInRawBytes(ctx, client, blockNum)
	if err != nil {
		return errors.Wrap(err, "getting current round from Algorand Node")
	}

	algodBlockBytes, err := json.Marshal(algodBlock)
	if err != nil {
		return errors.Wrap(err, "marshaling JSON")
	}

	indexerClient, err := indexer.Open(indexerCfg)
	if err != nil {
		return errors.Wrap(err, "connect to Indexer")
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexerCore := indexerApp.NewCore(log, indexerClient)

	indexerBlock, err := indexerCore.GetRound(ctx, traceID, blockNum)
	//rawBlock, err := app.GetRoundInRawBytes(ctx, client, blockNum)
	if err != nil {
		return errors.Wrap(err, "getting current round from Indexer")
	}

	fmt.Println("Algorand")
	fmt.Println("Genesish Hash")
	fmt.Println(algodBlock.Transactions[0].GenesisHash)
	fmt.Println(base64.StdEncoding.EncodeToString(algodBlock.Transactions[0].GenesisHash))
	fmt.Println("Indexer")
	fmt.Println("Genesish Hash")
	fmt.Println(indexerBlock.Transactions[0].GenesisHash)
	fmt.Println(base64.StdEncoding.EncodeToString(indexerBlock.Transactions[0].GenesisHash))

	fmt.Println(indexerBlock.Transactions[0].Id)

	indexerBlockBytes, err := json.Marshal(indexerBlock)
	if err != nil {
		return errors.Wrap(err, "marshaling JSON")
	}

	diffOpts := jsondiff.DefaultConsoleOptions()
	//res, diff := jsondiff.Compare(algodBlockBytes, indexerBlockBytes, &diffOpts)
	res, diff := jsondiff.Compare(indexerBlockBytes, algodBlockBytes, &diffOpts)

	if res != jsondiff.FullMatch {
		fmt.Println("Come on come on")
		fmt.Printf("%s\n", diff)
		return errors.Wrapf(nil, "the expected result is not equal to what we have: %s", diff)
	}

	return nil
}
