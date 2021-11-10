package algod

import (
	"context"
	"encoding/base64"
	algodv2 "github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/kevguy/algosearch/backend/business/core/block/db"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"strconv"
)

// GetRound retrieves a block from the Algod API based on the round number given
func GetRound(ctx context.Context, traceID string, log *zap.SugaredLogger, algodClient *algodv2.Client, roundNum uint64) (*db.NewBlock, error) {

	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "algod.GetRound")
	span.SetAttributes(attribute.Int64("round", int64(roundNum)))
	defer span.End()

	log.Infow("algod.GetRound", "traceid", traceID)

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

// GetCurrentRound retrieves retrieves the current block from the Algod API
func GetCurrentRound(ctx context.Context, traceID string, log *zap.SugaredLogger, algodClient *algodv2.Client) (*db.NewBlock, error) {

	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "algod.GetCurrentRound")
	defer span.End()

	log.Infow("algod.GetCurrentRound", "traceid", traceID)

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

// GetCurrentRoundNum retrieves the current round number from the Algod API
func GetCurrentRoundNum(ctx context.Context, algodClient *algodv2.Client) (uint64, error) {

	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "algod.GetCurrentRoundInRawBytes")
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
func GetRoundInRawBytes(ctx context.Context, algodClient *algodv2.Client, roundNum uint64) ([]byte, error) {

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
func ConvertBlockRawBytes(ctx context.Context, rawBlock []byte) (db.NewBlock, error) {

	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "algod.ConvertBlockRawBytes")
	//span.SetAttributes(attribute.String("query", q))
	defer span.End()

	var response models.BlockResponse
	err := msgpack.Decode(rawBlock, &response)
	if err != nil {
		return db.NewBlock{}, errors.Wrap(err, "parsing response")
	}

	//var blockInfo = block.NewBlock{response.Block, "", ""}
	var blockInfo = response.Block

	var rewards = models.BlockRewards{
		FeeSink:                 blockInfo.FeeSink.String(),
		RewardsCalculationRound: uint64(blockInfo.RewardsRecalculationRound),
		RewardsLevel:            blockInfo.RewardsLevel,
		RewardsPool:             blockInfo.RewardsPool.String(),
		RewardsRate:             blockInfo.RewardsRate,
		RewardsResidue:          blockInfo.RewardsResidue,
	}

	var upgradeState = models.BlockUpgradeState{
		CurrentProtocol:        blockInfo.CurrentProtocol,
		NextProtocol:           blockInfo.NextProtocol,
		NextProtocolApprovals:  blockInfo.NextProtocolApprovals,
		NextProtocolSwitchOn:   uint64(blockInfo.NextProtocolSwitchOn),
		NextProtocolVoteBefore: uint64(blockInfo.NextProtocolVoteBefore),
	}

	var upgradeVote = models.BlockUpgradeVote{
		UpgradeApprove: blockInfo.UpgradeApprove,
		UpgradeDelay:   uint64(blockInfo.UpgradeDelay),
		UpgradePropose: blockInfo.UpgradePropose,
	}

	var newBlock = db.NewBlock{
		Block:     models.Block{
			GenesisHash:       blockInfo.GenesisHash[:],
			GenesisId:         blockInfo.GenesisID,
			PreviousBlockHash: blockInfo.Branch[:],
			Rewards:           rewards,
			Round:             uint64(blockInfo.Round),
			Seed:              blockInfo.Seed[:],
			Timestamp:         uint64(blockInfo.TimeStamp),
			Transactions:      nil,
			TransactionsRoot:  blockInfo.TxnRoot[:],
			TxnCounter:        blockInfo.TxnCounter,
			UpgradeState:      upgradeState,
			UpgradeVote:       upgradeVote,
		},
		Proposer:  "",
		BlockHash: "",
	}

	// Process Genesis Hash
	//var genesisHash = [32]byte{}
	//copy(genesisHash[:], blockInfo.GenesisHash[:])
	var genesisHash = [32]byte(blockInfo.GenesisHash)
	//newBlock.GenesisHash = base64.StdEncoding.EncodeToString(genesisHash[:])
	newBlock.GenesisHash = genesisHash[:]

	// Process Previous Block Hash
	//var prevBlockHashStr = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(blockInfo.Branch[:])
	//newBlock.PrevBlockHash = "blk-" + prevBlockHashStr
	newBlock.PreviousBlockHash = blockInfo.Branch[:]

	// Process Seed
	//var seedStr = base64.StdEncoding.EncodeToString(blockInfo.Seed[:])
	newBlock.Seed = blockInfo.Seed[:]

	// Process Transactions
	newBlock.Transactions = []models.Transaction{}
	for _, txn := range blockInfo.Payset {
		newBlock.Transactions = append(newBlock.Transactions, ProcessTransactionInBlock(
			txn,
			blockInfo))
	}

	// Print Transactions Root
	// Don't use the String() from types.Address
	//newBlock.TransactionsRoot = base64.StdEncoding.EncodeToString(blockInfo.TxnRoot[:])
	newBlock.TransactionsRoot = blockInfo.TxnRoot[:]

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
