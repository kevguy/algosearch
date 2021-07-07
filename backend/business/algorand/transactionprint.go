package algorand

import (
	"encoding/base64"
	"fmt"
	"github.com/algorand/go-algorand-sdk/types"
)

func printPaymentTx(txn types.Transaction, padding int) {

	var paddingStr = ""
	for i := 0; i < padding; i++ {
		paddingStr += "\t"
	}

	fmt.Println(paddingStr + "\t- Payment Transaction")
	fmt.Println(paddingStr + "\t- Receiver: " + txn.Receiver.String())
	fmt.Printf(paddingStr + "\t- Amount: %d\n", uint64(txn.Amount))
	fmt.Println(paddingStr + "\t- Close Remainder To: " + txn.CloseRemainderTo.String())
}

func printKeyRegistrationTx(txn types.Transaction, padding int) {

	var paddingStr = ""
	for i := 0; i < padding; i++ {
		paddingStr += "\t"
	}

	fmt.Print(paddingStr + "\t- Vote Key: " + base64.StdEncoding.EncodeToString(txn.VotePK[:]))
	fmt.Print(paddingStr + "\t- Selection Key: " + base64.StdEncoding.EncodeToString(txn.SelectionPK[:]))
	fmt.Printf(paddingStr + "\t- Vote First: %d\n", uint64(txn.VoteFirst))
	fmt.Printf(paddingStr + "\t- Vote Last: %d\n", uint64(txn.VoteLast))
	fmt.Printf(paddingStr + "\t- Vote Key Dilution: %d\n", txn.VoteKeyDilution)
	fmt.Printf(paddingStr + "\t- Non Participation: %t\n", txn.Nonparticipation)
}

func printAssetTransferTx(txn types.Transaction, padding int) {

	var paddingStr = ""
	for i := 0; i < padding; i++ {
		paddingStr += "\t"
	}

	fmt.Printf(paddingStr + "\t- XferAsset: %d\n", txn.XferAsset)
	fmt.Printf(paddingStr + "\t- Asset Amount: %d\n", txn.AssetAmount)
	fmt.Println(paddingStr + "\t- Asset Sender:" + txn.AssetSender.String())
	fmt.Println(paddingStr + "\t- Asset Receiver:" + txn.AssetSender.String())
	fmt.Println(paddingStr + "\t- Asset Close To:" + txn.AssetCloseTo.String())
}

func printAssetConfigTx(txn types.Transaction, padding int) {

	var paddingStr = ""
	for i := 0; i < padding; i++ {
		paddingStr += "\t"
	}

	fmt.Printf(paddingStr + "\t- Config Asset Index: %d\n", txn.ConfigAsset)
	fmt.Println(paddingStr + "\t- Asset Params:")
	fmt.Printf(paddingStr + "\t\t- Total: %d\n", txn.AssetParams.Total)
	fmt.Printf(paddingStr + "\t\t- Decimals: %d\n", txn.AssetParams.Decimals)
	fmt.Printf(paddingStr + "\t\t- Default Frozen: %t\n", txn.AssetParams.DefaultFrozen)
	fmt.Println(paddingStr + "\t\t- Unit Name: " + txn.AssetParams.UnitName)
	fmt.Println(paddingStr + "\t\t- Asset Name: " + txn.AssetParams.AssetName)
	fmt.Println(paddingStr + "\t\t- URL: " + txn.AssetParams.URL)
	fmt.Println(paddingStr + "\t\t- MetadataHash: " + base64.StdEncoding.EncodeToString(txn.AssetParams.MetadataHash[:]))
	fmt.Println(paddingStr + "\t\t- Manager: " + txn.AssetParams.Manager.String())
	fmt.Println(paddingStr + "\t\t- Reserve: " + txn.AssetParams.Reserve.String())
	fmt.Println(paddingStr + "\t\t- Freeze: " + txn.AssetParams.Freeze.String())
	fmt.Println(paddingStr + "\t\t- Clawback: " + txn.AssetParams.Clawback.String())
}

func printAssetFreezeTx(txn types.Transaction, padding int) {

	var paddingStr = ""
	for i := 0; i < padding; i++ {
		paddingStr += "\t"
	}

	fmt.Println(paddingStr + "\t- Freeze Account:" + txn.FreezeAccount.String())
	fmt.Printf(paddingStr + "\t- Freeze Asset: %d\n", txn.FreezeAsset)
	fmt.Printf(paddingStr + "\t- Asset Frozen: %t\n", txn.AssetFrozen)
}

func PrintTransactionInBlock(txn types.SignedTxnInBlock, padding int) {

	var paddingStr = ""
	for i := 0; i < padding; i++ {
		paddingStr += "\t"
	}

	var padding2Str = ""
	fmt.Printf(paddingStr + "- Has Genesis ID: %t\n", txn.HasGenesisID)
	fmt.Printf(paddingStr + "- Has Genesis Hash: %t\n", txn.HasGenesisHash)
	fmt.Println(paddingStr + "- SignedTxnWithAD")

	// - SignedTxn
	padding2Str = "\t"
	fmt.Println(paddingStr + "- SignedTxn")

	fmt.Println(paddingStr + padding2Str + "- Signature: " + base64.StdEncoding.EncodeToString(txn.Sig[:]))
	fmt.Println(paddingStr + padding2Str + "- Sub Siguatures (Msig MultisigSig): Not processing it")
	fmt.Println(paddingStr + padding2Str + "- Logic Signature (LSig LogicSig): Not processing it")
	fmt.Println(paddingStr + padding2Str + "- Auth Addreess: " + txn.AuthAddr.String())

	// -- Transaction Info
	padding2Str = "\t\t"
	fmt.Println(paddingStr + "\t- Transaction Info:")

	// --- Type
	fmt.Println(paddingStr + padding2Str + "- Type: " + string(txn.Txn.Type))

	// --- Header
	padding2Str = "\t\t\t"
	fmt.Println(paddingStr + "\t\t- Header")
	fmt.Println(paddingStr + padding2Str + "- Sender: " + txn.Txn.Sender.String())
	fmt.Printf(paddingStr + padding2Str + "- Fee: %d\n", uint64(txn.Txn.Fee))
	fmt.Printf(paddingStr + padding2Str + "- First Valid Round: %d\n", uint64(txn.Txn.FirstValid))
	fmt.Printf(paddingStr + padding2Str + "- Last Valid Round: %d\n", uint64(txn.Txn.LastValid))
	fmt.Println(paddingStr + padding2Str + "- Note: " + base64.StdEncoding.EncodeToString(txn.Txn.Note))
	fmt.Println(paddingStr + padding2Str + "- Genesis ID: " + txn.Txn.GenesisID)
	var genesisHash = [32]byte(txn.Txn.GenesisHash)
	var genesisHashStr = base64.StdEncoding.EncodeToString(genesisHash[:])
	fmt.Println(paddingStr + padding2Str + "- Genesis Hash: " + genesisHashStr)
	var group = [32]byte(txn.Txn.Group)
	var groupStr = base64.StdEncoding.EncodeToString(group[:])
	fmt.Println(paddingStr + padding2Str + "- Group: " + groupStr)
	fmt.Println(paddingStr + padding2Str + "- Lease (Not Verified): " + base64.StdEncoding.EncodeToString(txn.Txn.Lease[:]))
	fmt.Println(paddingStr + padding2Str + "- RekeyTo: " + txn.Txn.RekeyTo.String())

	fmt.Println(paddingStr + padding2Str + "- Transaction Specific Info for:" + string(txn.Txn.Type))
	switch txn.Txn.Type {
	case types.PaymentTx:
		printPaymentTx(txn.Txn, 6)
		break
	case types.KeyRegistrationTx:
		printKeyRegistrationTx(txn.Txn, 6)
		break
	case types.AssetConfigTx:
		printAssetConfigTx(txn.Txn, 6)
		break
	case types.AssetTransferTx:
		printAssetTransferTx(txn.Txn, 6)
		break
	case types.AssetFreezeTx:
		printAssetFreezeTx(txn.Txn, 6)
		break
	case types.ApplicationCallTx:
		// TODO: finished this
		break
	}

	// - ApplyData
	padding2Str = "\t\t"
	fmt.Println(paddingStr + "\t- ApplyData")
	fmt.Printf(paddingStr + padding2Str + "- Closing Amount: %d\n", uint64(txn.ClosingAmount))
	fmt.Printf(paddingStr + padding2Str + "- Asset Closing Amount: %d\n", txn.AssetClosingAmount)
	fmt.Printf(paddingStr + padding2Str + "- Sender Rewards: %d\n", uint64(txn.SenderRewards))
	fmt.Printf(paddingStr + padding2Str + "- Receiver Rewards: %d\n", uint64(txn.ReceiverRewards))
	fmt.Printf(paddingStr + padding2Str + "- Close Rewards: %d\n", uint64(txn.CloseRewards))
	fmt.Println(paddingStr + padding2Str + "- EvalDelta is ignored")
}
