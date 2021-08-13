package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	app "github.com/kevguy/algosearch/backend/business/indexer"
	"github.com/kevguy/algosearch/backend/foundation/indexer"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

// PrettyPrintBlockFromIndexerCmd retrieves information about the block for the specified round and prints it out in JSON format.
func PrettyPrintBlockFromIndexerCmd(traceID string, log *zap.SugaredLogger, cfg indexer.Config, blockNum uint64) error {
	client, err := indexer.Open(cfg)
	if err != nil {
		return errors.Wrap(err, "connect to Algorand Node")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	block, err := app.GetRound(ctx, traceID, log, client, blockNum)
	//rawBlock, err := app.GetRoundInRawBytes(ctx, client, blockNum)
	if err != nil {
		return errors.Wrap(err, "getting current round from Indexer")
	}

	blockBytes, err := json.Marshal(block)
	if err != nil {
		return errors.Wrap(err, "marshaling JSON")
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, blockBytes, "", "\t")
	if err != nil {
		//fmt.Println("JSON parse error: ", err)
		return errors.Wrap(err, "JSON parse")
	}

	//fmt.Println("Block Info:", string(prettyJSON.Bytes()))
	fmt.Println("Block Info:", prettyJSON.String())

	return nil
}
