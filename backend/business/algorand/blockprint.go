package algorand

import (
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"strconv"
)

// PrintBlockInfoFromRawBytes processes the raw block bytes and prints all the block information out
func PrintBlockInfoFromRawBytes(rawBlock []byte) error {

	var response models.BlockResponse
	err := msgpack.Decode(rawBlock, &response)
	if err != nil {
		fmt.Printf("error parsing response: %s\n", err)
		return err
	}

	var blockInfo = response.Block

	fmt.Println("========================================================")
	fmt.Println("Block Information:")

	// Print Genesis Hash
	//var genesisHash = [32]byte{}
	//copy(genesisHash[:], blockInfo.GenesisHash[:])
	var genesisHash = [32]byte(blockInfo.GenesisHash)
	var genesisHashStr = base64.StdEncoding.EncodeToString(genesisHash[:])
	fmt.Println("\t- Genesis Hash: " + genesisHashStr)

	// Print Genesis ID
	fmt.Println("\t- Genesis ID: " + blockInfo.GenesisID)

	// Print Previous Block Hash
	var prevBlockHashStr = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(blockInfo.Branch[:])
	fmt.Println("\t- Previous Block Hash: blk-" + prevBlockHashStr)

	// Print Reward
	fmt.Println("\t- Rewards:")

	// Print Fee Sink
	fmt.Println("\t\t- FeeSink: " + blockInfo.FeeSink.String())

	// Print Reward Calculation Round
	fmt.Println("\t\t- Reward Calculation Round: " + strconv.FormatUint(uint64(blockInfo.RewardsRecalculationRound), 10))

	// Print Rewards Level
	fmt.Printf("\t\t- Reward Level: %d\n", blockInfo.RewardsLevel)
	fmt.Println("\t\t- Reward Level: " + strconv.FormatUint(blockInfo.RewardsLevel, 10))

	// Print Rewards Pool
	fmt.Println("\t\t- Rewards Pool: " + blockInfo.RewardsPool.String())

	// Print Rewards Rate
	fmt.Printf("\t\t- Rewards Rate: %d\n", blockInfo.RewardsRate)
	fmt.Println("\t\t- Rewards Rate: " + strconv.FormatUint(blockInfo.RewardsRate, 10))

	// Print Rewards Residue
	fmt.Printf("\t\t- Rewards Residue: %d\n", blockInfo.RewardsResidue)
	fmt.Println("\t\t- Rewards Residue: " + strconv.FormatUint(blockInfo.RewardsResidue, 10))

	// Print Round
	fmt.Printf("\t- Round: %d\n", blockInfo.Round)
	fmt.Println("\t- Round: " + strconv.FormatUint(uint64(blockInfo.Round), 10))

	// Print Seed
	var seedStr = base64.StdEncoding.EncodeToString(blockInfo.Seed[:])
	fmt.Println("\t- Seed: " + seedStr)

	// Print Timestamp
	fmt.Printf("\t- Timestamp: %d\n", blockInfo.TimeStamp)

	// Print Transaction
	fmt.Println("\t- Transactions:" + string(len(blockInfo.Payset)))
	for idx, txn := range blockInfo.Payset {
		fmt.Printf("\t\t- Transaction %d\n", idx)
		PrintTransactionInBlock(txn, 3)
	}

	// Print Transactions Root
	// Don't use the String() from types.Address
	fmt.Printf("\t- Transactions Root: %s\n", base64.StdEncoding.EncodeToString(blockInfo.TxnRoot[:]))

	// Print Transactions Counter
	fmt.Println("\t- Transactions Counter: " + strconv.FormatUint(blockInfo.TxnCounter, 10))

	// Print Upgrade State
	fmt.Println("\t- Upgrade State")

	// Print Current Protocol
	fmt.Println("\t\t- Current Protocol: " + blockInfo.CurrentProtocol)

	// Print Next Protocol
	fmt.Println("\t\t- Next Protocol: " + blockInfo.NextProtocol)

	// Print Next Protocol Approvals
	fmt.Println("\t\t- Next Protocol Approvals: " + strconv.FormatUint(blockInfo.NextProtocolApprovals, 10))

	// Print Next Protocol Switch On
	fmt.Println("\t\t- Next Protocol Switch On: " + strconv.FormatUint(uint64(blockInfo.NextProtocolSwitchOn), 10))

	// Print Next Protocol Vote Before
	fmt.Println("\t\t- Next Protocol Vote Before: " + strconv.FormatUint(uint64(blockInfo.NextProtocolVoteBefore), 10))

	// Print Upgrade Vote
	fmt.Println("\t- Upgrade Vote")

	// Print Upgrade Propose
	fmt.Println("\t\t- Upgrade Propose: " + blockInfo.UpgradePropose)

	// Print Upgrade Approve
	fmt.Println("\t\t- Upgrade Approve: " + strconv.FormatBool(blockInfo.UpgradeApprove))

	// Print Upgrade Delay
	fmt.Println("\t\t- Upgrade Delay: " + strconv.FormatUint(uint64(blockInfo.UpgradeDelay), 10))

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
