package algod

import (
	"context"
	"encoding/base32"
	"encoding/base64"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/kevguy/algosearch/backend/business/couchdata/block"
	"github.com/kevguy/algosearch/backend/business/couchdata/transaction"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"strconv"
)

func GetRound(ctx context.Context, traceID string, log *zap.SugaredLogger, algodClient *algod.Client, roundNum uint64) (*block.NewBlock, error) {
	log.Infow("algorand.GetGetRound", "traceid", traceID)

	rawBlock, err := GetRoundInRawBytes(ctx, algodClient, roundNum)
	if err != nil {
		return nil, errors.Wrap(err, "unable to query for current round")
	}

	blockData, err := ConvertBlockRawBytes(ctx, rawBlock)
	if err != nil {
		return nil, errors.Wrap(err, "unable to convert raw block for current round")
	}
	return &blockData, nil
}

func GetCurrentRound(ctx context.Context, traceID string, log *zap.SugaredLogger, algodClient *algod.Client) (*block.NewBlock, error) {
	log.Infow("algorand.GetCurrentRound", "traceid", traceID)

	currNum, err := GetCurrentRoundNum(ctx, algodClient)
	if err != nil {
		return nil, errors.Wrap(err, "unable to query for current round num")
	}

	blockData, err := GetRound(ctx, traceID, log, algodClient, currNum)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get round data")
	}
	return blockData, nil
}

// GetCurrentRoundNum retrieves the current round number.
func GetCurrentRoundNum(ctx context.Context, algodClient *algod.Client) (uint64, error) {

	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "algorand.GetCurrentRoundInRawBytes")
	//span.SetAttributes(attribute.String("query", q))
	defer span.End()

	nodeStatus, pingError := algodClient.Status().Do(ctx)
	if pingError != nil {
		return 0, errors.Wrap(pingError, "Error getting node status")
	}

	//fmt.Println("Current Round: " + strconv.FormatUint(nodeStatus.LastRound, 10))
	//fmt.Printf("algod last round: %d\n", nodeStatus.LastRound)

	return nodeStatus.LastRound, nil
}

// GetRoundInRawBytes retrieves the specified round and returns result in byte format.
func GetRoundInRawBytes(ctx context.Context, algodClient *algod.Client, roundNum uint64) ([]byte, error) {

	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "algorand.GetRoundInRawBytes")
	span.SetAttributes(attribute.String("blockNum", strconv.FormatUint(roundNum, 10)))
	defer span.End()

	rawBlock, err := algodClient.BlockRaw(roundNum).Do(ctx)
	if err != nil {
		return []byte{}, errors.Wrap(err, "getting ground in raw bytes")
	}

	return rawBlock, err
}

// ConvertBlockRawBytes processes the block in byte format and returns it in a struct.
func ConvertBlockRawBytes(ctx context.Context, rawBlock []byte) (block.NewBlock, error) {

	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "algorand.ConvertBlockRawBytes")
	//span.SetAttributes(attribute.String("query", q))
	defer span.End()

	var response models.BlockResponse
	err := msgpack.Decode(rawBlock, &response)
	if err != nil {
		return block.NewBlock{}, errors.Wrap(err, "parsing response")
	}

	//var blockInfo = block.NewBlock{response.Block, "", ""}
	var blockInfo = response.Block

	var rewards = block.Rewards{
		FeeSink:			blockInfo.FeeSink.String(),
		RewardsCalRound:	uint64(blockInfo.RewardsRecalculationRound),
		RewardsLevel:		blockInfo.RewardsLevel,
		RewardsPool:		blockInfo.RewardsPool.String(),
		RewardsRate:		blockInfo.RewardsRate,
		RewardsResidue:		blockInfo.RewardsResidue,
	}

	var upgradeState = block.UpgradeState{
		CurrProtocol: 			blockInfo.CurrentProtocol,
		NextProtocol:			&blockInfo.NextProtocol,
		NextProtocolApprovals:	blockInfo.NextProtocolApprovals,
		NextProtocolSwitchOn:	uint64(blockInfo.NextProtocolSwitchOn),
		NextProtocolVoteBefore:	uint64(blockInfo.NextProtocolVoteBefore),
	}

	var upgradeVote = block.UpgradeVote{
		UpgradeApprove:	blockInfo.UpgradeApprove,
		UpgradeDelay: 	uint64(blockInfo.UpgradeDelay),
		UpgradePropose: &blockInfo.UpgradePropose,
	}

	var newBlock = block.NewBlock{
		GenesisHash:		"",
		GenesisID:			blockInfo.GenesisID,
		PrevBlockHash:		"",
		Rewards:			rewards,
		Round: 				uint64(blockInfo.Round),
		Seed: 				"",
		Timestamp:			uint64(blockInfo.TimeStamp),
		//Transactions		[]uint64		`json:"transactions"`
		TransactionsRoot:	"",
		TransactionCounter:	blockInfo.TxnCounter,
		UpgradeState:		upgradeState,
		UpgradeVote:		upgradeVote,
		Proposer:			"",
		BlockHash:			"",
	}

	// Process Genesis Hash
	//var genesisHash = [32]byte{}
	//copy(genesisHash[:], blockInfo.GenesisHash[:])
	var genesisHash = [32]byte(blockInfo.GenesisHash)
	newBlock.GenesisHash = base64.StdEncoding.EncodeToString(genesisHash[:])

	// Process Previous Block Hash
	var prevBlockHashStr = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(blockInfo.Branch[:])
	//newBlock.PrevBlockHash = "blk-" + prevBlockHashStr
	newBlock.PrevBlockHash = prevBlockHashStr

	// Process Seed
	var seedStr = base64.StdEncoding.EncodeToString(blockInfo.Seed[:])
	newBlock.Seed = seedStr

	// Process Transactions
	newBlock.Transactions = []transaction.Transaction{}
	for _, txn := range blockInfo.Payset {
		newBlock.Transactions = append(newBlock.Transactions, ProcessTransactionInBlock(txn))
	}

	// Print Transactions Root
	// Don't use the String() from types.Address
	newBlock.TransactionsRoot = base64.StdEncoding.EncodeToString(blockInfo.TxnRoot[:])

	var certInfo = *response.Cert
	var prop = certInfo["prop"].(map[interface{}]interface{})

	// Find the Proposer, that is the correct implementation
	var oprop = prop["oprop"].([]byte)
	oprop_ := byteArrAsAddress(oprop)
	newBlock.Proposer = oprop_.String()

	// Find the Block Hash
	var dig = prop["dig"].([]byte)
	newBlock.BlockHash = base64.StdEncoding.EncodeToString(dig)

	return newBlock, nil
}
