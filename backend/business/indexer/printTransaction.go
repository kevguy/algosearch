package indexer

import (
	"encoding/base64"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/types"
)

func printPaymentTx(txn models.TransactionPayment, padding int) {

	var paddingStr = ""
	for i := 0; i < padding; i++ {
		paddingStr += "\t"
	}

	fmt.Println(paddingStr + "\t- Payment Transaction")
	fmt.Println(paddingStr + "\t- Receiver: " + txn.Receiver)
	fmt.Printf(paddingStr + "\t- Amount: %d\n", txn.Amount)
	fmt.Println(paddingStr + "\t- Close Remainder To: " + txn.CloseRemainderTo)
}

func printKeyRegistrationTx(txn models.TransactionKeyreg, padding int) {

	var paddingStr = ""
	for i := 0; i < padding; i++ {
		paddingStr += "\t"
	}

	fmt.Print(paddingStr + "\t- Vote Participation Key: " + base64.StdEncoding.EncodeToString(txn.VoteParticipationKey[:]))
	fmt.Print(paddingStr + "\t- Selection Key: " + base64.StdEncoding.EncodeToString(txn.SelectionParticipationKey[:]))
	fmt.Printf(paddingStr + "\t- Vote First Valid: %d\n", txn.VoteFirstValid)
	fmt.Printf(paddingStr + "\t- Vote Last Valid: %d\n", txn.VoteLastValid)
	fmt.Printf(paddingStr + "\t- Vote Key Dilution: %d\n", txn.VoteKeyDilution)
	fmt.Printf(paddingStr + "\t- Non Participation: %t\n", txn.NonParticipation)
}

func printAssetTransferTx(txn models.TransactionAssetTransfer, padding int) {

	var paddingStr = ""
	for i := 0; i < padding; i++ {
		paddingStr += "\t"
	}

	fmt.Printf(paddingStr + "\t- XferAsset: %d\n", txn.AssetId)
	fmt.Printf(paddingStr + "\t- Asset Amount: %d\n", txn.Amount)
	fmt.Println(paddingStr + "\t- Asset Sender:" + txn.Sender)
	fmt.Println(paddingStr + "\t- Asset Receiver:" + txn.Receiver)
	fmt.Println(paddingStr + "\t- Asset Close To:" + txn.CloseTo)
	// Note that I can't find this is block msgpack data
	fmt.Printf(paddingStr + "\t- Asset Close Amount: %d\n", txn.CloseAmount)
}

func printAssetConfigTx(txn models.TransactionAssetConfig, padding int) {

	var paddingStr = ""
	for i := 0; i < padding; i++ {
		paddingStr += "\t"
	}

	fmt.Printf(paddingStr + "\t- Config Asset Index: %d\n", txn.AssetId)
	fmt.Println(paddingStr + "\t- Asset Params:")
	fmt.Printf(paddingStr + "\t\t- Total: %d\n", txn.Params.Total)
	fmt.Printf(paddingStr + "\t\t- Decimals: %d\n", txn.Params.Decimals)
	fmt.Printf(paddingStr + "\t\t- Default Frozen: %t\n", txn.Params.DefaultFrozen)
	fmt.Println(paddingStr + "\t\t- Unit Name: " + txn.Params.UnitName)
	fmt.Println(paddingStr + "\t\t- Asset Name: " + txn.Params.Name)
	fmt.Println(paddingStr + "\t\t- URL: " + txn.Params.Url)
	fmt.Println(paddingStr + "\t\t- MetadataHash: " + base64.StdEncoding.EncodeToString(txn.Params.MetadataHash[:]))
	fmt.Println(paddingStr + "\t\t- Manager: " + txn.Params.Manager)
	fmt.Println(paddingStr + "\t\t- Reserve: " + txn.Params.Reserve)
	fmt.Println(paddingStr + "\t\t- Freeze: " + txn.Params.Freeze)
	fmt.Println(paddingStr + "\t\t- Clawback: " + txn.Params.Clawback)
}

func printAssetFreezeTx(txn models.TransactionAssetFreeze, padding int) {

	var paddingStr = ""
	for i := 0; i < padding; i++ {
		paddingStr += "\t"
	}

	fmt.Println(paddingStr + "\t- Freeze Account:" + txn.Address)
	fmt.Printf(paddingStr + "\t- Freeze Asset: %d\n", txn.AssetId)
	fmt.Printf(paddingStr + "\t- Asset Frozen: %t\n", txn.NewFreezeStatus)
}

func PrintTransactionInBlock(txn models.Transaction, padding int) {

	var paddingStr = ""
	for i := 0; i < padding; i++ {
		paddingStr += "\t"
	}

	var padding2Str = ""
	if len(txn.GenesisId) == 0 {
		fmt.Printf(paddingStr + "- Has Genesis ID: false\n")
	} else {
		fmt.Printf(paddingStr + "- Has Genesis ID: true\n")
	}
	if len(txn.GenesisHash) == 0 {
		fmt.Printf(paddingStr + "- Has Genesis Hash: false\n")
	} else {
		fmt.Printf(paddingStr + "- Has Genesis Hash: true\n")
	}
	fmt.Println(paddingStr + "- SignedTxnWithAD")

	// - SignedTxn
	padding2Str = "\t"
	fmt.Println(paddingStr + "- SignedTxn")

	fmt.Println(paddingStr + padding2Str + "- Signature: " + base64.StdEncoding.EncodeToString(txn.Signature.Sig[:]))
	// TODO: Process it
	fmt.Println(paddingStr + padding2Str + "- Sub Siguatures (Msig MultisigSig): Not processing it")
	// TODO: Process it
	fmt.Println(paddingStr + padding2Str + "- Logic Signature (LSig LogicSig): Not processing it")
	fmt.Println(paddingStr + padding2Str + "- Auth Addreess: " + txn.AuthAddr)

	// -- Transaction Info
	padding2Str = "\t\t"
	fmt.Println(paddingStr + "\t- Transaction Info:")

	// --- Type
	fmt.Println(paddingStr + padding2Str + "- Type: " + txn.Type)

	// --- Header
	padding2Str = "\t\t\t"
	fmt.Println(paddingStr + "\t\t- Header")
	fmt.Println(paddingStr + padding2Str + "- Sender: " + txn.Sender)
	fmt.Printf(paddingStr + padding2Str + "- Fee: %d\n", txn.Fee)
	fmt.Printf(paddingStr + padding2Str + "- First Valid Round: %d\n", txn.FirstValid)
	fmt.Printf(paddingStr + padding2Str + "- Last Valid Round: %d\n", txn.LastValid)
	fmt.Println(paddingStr + padding2Str + "- Note: " + base64.StdEncoding.EncodeToString(txn.Note))
	fmt.Println(paddingStr + padding2Str + "- Genesis ID: " + txn.GenesisId)
	var genesisHashStr = base64.StdEncoding.EncodeToString(txn.GenesisHash[:])
	fmt.Println(paddingStr + padding2Str + "- Genesis Hash: " + genesisHashStr)
	var groupStr = base64.StdEncoding.EncodeToString(txn.Group[:])
	fmt.Println(paddingStr + padding2Str + "- Group: " + groupStr)
	fmt.Println(paddingStr + padding2Str + "- Lease (Not Verified): " + base64.StdEncoding.EncodeToString(txn.Lease[:]))
	fmt.Println(paddingStr + padding2Str + "- RekeyTo: " + txn.RekeyTo)

	fmt.Println(paddingStr + padding2Str + "- Transaction Specific Info for:" + txn.Type)
	switch txn.Type {
	case string(types.PaymentTx):
		printPaymentTx(txn.PaymentTransaction, 6)
		break
	case string(types.KeyRegistrationTx):
		printKeyRegistrationTx(txn.KeyregTransaction, 6)
		break
	case string(types.AssetConfigTx):
		printAssetConfigTx(txn.AssetConfigTransaction, 6)
		break
	case string(types.AssetTransferTx):
		printAssetTransferTx(txn.AssetTransferTransaction, 6)
		break
	case string(types.AssetFreezeTx):
		printAssetFreezeTx(txn.AssetFreezeTransaction, 6)
		break
	case string(types.ApplicationCallTx):
		// TODO: finish this
		break
	}

	// - ApplyData
	padding2Str = "\t\t"
	fmt.Println(paddingStr + "\t- ApplyData")
	fmt.Printf(paddingStr + padding2Str + "- Closing Amount: %d\n", txn.ClosingAmount)
	//fmt.Printf(paddingStr + padding2Str + "- Asset Closing Amount: %d\n", txn.AssetClosingAmount)
	fmt.Printf(paddingStr + padding2Str + "- Asset Closing Amount: Unavailable in data\n")
	fmt.Printf(paddingStr + padding2Str + "- Sender Rewards: %d\n", txn.SenderRewards)
	fmt.Printf(paddingStr + padding2Str + "- Receiver Rewards: %d\n", txn.ReceiverRewards)
	fmt.Printf(paddingStr + padding2Str + "- Close Rewards: %d\n", txn.CloseRewards)
	fmt.Println(paddingStr + padding2Str + "- EvalDelta is ignored")
}
