package indexer

import (
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"strconv"
)

func PrintBlockInfoFromJsonBlock(jsonBlock models.Block) error {
	fmt.Println("========================================================")
	fmt.Println("Block Information:")

	// Print Genesis Hash
	var genesisHashStr = base64.StdEncoding.EncodeToString(jsonBlock.GenesisHash[:])
	fmt.Println("\t- Genesis Hash: " + genesisHashStr)

	// Print Genesis ID
	fmt.Println("\t- Genesis ID: " + jsonBlock.GenesisId)

	// Print Previous Block Hash
	var prevBlockHashStr = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(jsonBlock.PreviousBlockHash[:])
	fmt.Println("\t- Previous Block Hash: blk-" + prevBlockHashStr)

	// Print Reward
	fmt.Println("\t- Rewards:")

	// Print Fee Sink
	fmt.Println("\t\t- FeeSink: " + jsonBlock.Rewards.FeeSink)

	// Print Reward Calculation Round
	fmt.Println("\t\t- Reward Calculation Round: " + strconv.FormatUint(jsonBlock.Rewards.RewardsCalculationRound, 10))

	// Print Rewards Level
	fmt.Printf("\t\t- Reward Level: %d\n", jsonBlock.Rewards.RewardsLevel)
	fmt.Println("\t\t- Reward Level: " + strconv.FormatUint(jsonBlock.Rewards.RewardsLevel, 10))

	// Print Rewards Pool
	fmt.Println("\t\t- Rewards Pool: " + jsonBlock.Rewards.RewardsPool)

	// Print Rewards Rate
	fmt.Printf("\t\t- Rewards Rate: %d\n", jsonBlock.Rewards.RewardsRate)
	fmt.Println("\t\t- Rewards Rate: " + strconv.FormatUint(jsonBlock.Rewards.RewardsRate, 10))

	// Print Rewards Residue
	fmt.Printf("\t\t- Rewards Residue: %d\n", jsonBlock.Rewards.RewardsResidue)
	fmt.Println("\t\t- Rewards Residue: " + strconv.FormatUint(jsonBlock.Rewards.RewardsResidue, 10))

	// Print Round
	fmt.Printf("\t- Round: %d\n", jsonBlock.Round)
	fmt.Println("\t- Round: " + strconv.FormatUint(jsonBlock.Round, 10))

	// Print Seed
	var seedStr = base64.StdEncoding.EncodeToString(jsonBlock.Seed[:])
	fmt.Println("\t- Seed: " + seedStr)

	// Print Timestamp
	fmt.Printf("\t- Timestamp: %d\n", jsonBlock.Timestamp)

	// Print Transaction
	fmt.Println("\t- Transactions:" + string(len(jsonBlock.Transactions)))
	for idx, txn := range jsonBlock.Transactions {
		fmt.Printf("\t\t- Transaction %d\n", idx)
		PrintTransactionInBlock(txn, 3)
	}

	// Print Transactions Root
	// Don't use the String() from types.Address
	fmt.Printf("\t- Transactions Root: %s\n", base64.StdEncoding.EncodeToString(jsonBlock.TransactionsRoot[:]))

	// Print Transactions Counter
	fmt.Println("\t- Transactions Counter: " + strconv.FormatUint(jsonBlock.TxnCounter, 10))

	// Print Upgrade State
	fmt.Println("\t- Upgrade State")

	// Print Current Protocol
	fmt.Println("\t\t- Current Protocol: " + jsonBlock.UpgradeState.CurrentProtocol)

	// Print Next Protocol
	fmt.Println("\t\t- Next Protocol: " + jsonBlock.UpgradeState.NextProtocol)

	// Print Next Protocol Approvals
	fmt.Println("\t\t- Next Protocol Approvals: " + strconv.FormatUint(jsonBlock.UpgradeState.NextProtocolApprovals, 10))

	// Print Next Protocol Switch On
	fmt.Println("\t\t- Next Protocol Switch On: " + strconv.FormatUint(jsonBlock.UpgradeState.NextProtocolSwitchOn, 10))

	// Print Next Protocol Vote Before
	fmt.Println("\t\t- Next Protocol Vote Before: " + strconv.FormatUint(jsonBlock.UpgradeState.NextProtocolVoteBefore, 10))

	// Print Upgrade Vote
	fmt.Println("\t- Upgrade Vote")

	// Print Upgrade Propose
	fmt.Println("\t\t- Upgrade Propose: " + jsonBlock.UpgradeVote.UpgradePropose)

	// Print Upgrade Approve
	fmt.Println("\t\t- Upgrade Approve: " + strconv.FormatBool(jsonBlock.UpgradeVote.UpgradeApprove))

	// Print Upgrade Delay
	fmt.Println("\t\t- Upgrade Delay: " + strconv.FormatUint(uint64(jsonBlock.UpgradeVote.UpgradeDelay), 10))

	//const proposer = algosdk.encodeAddress(blk["cert"]["prop"]["oprop"]);
	//const blockHash = Buffer.from(blk["cert"]["prop"]["dig"]).toString("base64");
	//fmt.Println("Cert Info")
	//fmt.Println(response.Cert)

	var certInfo = *response.Cert
	var prop = certInfo["prop"].(map[interface{}]interface{})

	// Find the Proposer, that is the correct implementation
	var oprop = prop["oprop"].([]byte)
	oprop_ := byteArrAsAddress(oprop)
	fmt.Println("\t- Proposer: " + oprop_.String())

	// Find the Block Hash
	var dig = prop["dig"].([]byte)
	fmt.Println("\t- Block Hash: " + base64.StdEncoding.EncodeToString(dig))

	fmt.Println("========================================================")

	return nil
}
