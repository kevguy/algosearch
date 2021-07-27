package algod

import (
	"encoding/base64"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/kevguy/algosearch/backend/business/couchdata/transaction"
)

func ProcessPaymentTx(txn types.Transaction) transaction.PaymentTxnFields {
	return transaction.PaymentTxnFields{
		Receiver: txn.Receiver.String(),
		Amount: uint64(txn.Amount),
		CloseRemainderTo: txn.CloseRemainderTo.String(),
	}
}

func ProcessKeyRegistrationTx(txn types.Transaction) transaction.KeyregTxnFields {
	return transaction.KeyregTxnFields{
		VotePK: base64.StdEncoding.EncodeToString(txn.VotePK[:]),
		SelectionPK: base64.StdEncoding.EncodeToString(txn.SelectionPK[:]),
		VoteFirst: uint64(txn.VoteFirst),
		VoteLast: uint64(txn.VoteLast),
		VoteKeyDilution: txn.VoteKeyDilution,
		Nonparticipation: txn.Nonparticipation,
	}
}

func ProcessAssetTransferTx(txn types.Transaction) transaction.AssetTransferTxnFields {
	return transaction.AssetTransferTxnFields{
		XferAsset: uint64(txn.XferAsset),
		AssetAmount: txn.AssetAmount,
		AssetSender: txn.AssetSender.String(),
		AssetReceiver: txn.AssetReceiver.String(),
		AssetCloseTo: txn.AssetCloseTo.String(),
	}
}

func ProcessAssetConfigTx(txn types.Transaction) transaction.AssetConfigTxnFields {
	return transaction.AssetConfigTxnFields{
		ConfigAsset: uint64(txn.ConfigAsset),
		AssetParams: transaction.AssetParams{
			Total: txn.AssetParams.Total,
			Decimals: txn.AssetParams.Decimals,
			DefaultFrozen: txn.AssetParams.DefaultFrozen,
			UnitName: txn.AssetParams.UnitName,
			URL: txn.AssetParams.URL,
			MetadataHash: base64.StdEncoding.EncodeToString(txn.AssetParams.MetadataHash[:]),
			Manager: txn.AssetParams.Manager.String(),
			Reserve: txn.AssetParams.Reserve.String(),
			Freeze: txn.AssetParams.Freeze.String(),
			Clawback: txn.AssetParams.Clawback.String(),
		},
	}
}

func ProcessAssetFreezeTx(txn types.Transaction) transaction.AssetFreezeTxnFields {
	return transaction.AssetFreezeTxnFields{
		FreezeAccount: txn.FreezeAccount.String(),
		FreezeAsset: uint64(txn.FreezeAsset),
		AssetFrozen: txn.AssetFrozen,
	}
}

func ProcessHeader(data types.Header) transaction.Header {

	var genesisHash = [32]byte(data.GenesisHash)
	var genesisHashStr = base64.StdEncoding.EncodeToString(genesisHash[:])

	var group = [32]byte(data.Group)
	var groupStr = base64.StdEncoding.EncodeToString(group[:])

	return transaction.Header{
		Sender: data.Sender.String(),
		Fee: uint64(data.Fee),
		FirstValid: uint64(data.FirstValid),
		LastValid: uint64(data.LastValid),
		Note: base64.StdEncoding.EncodeToString(data.Note),
		GenesisID: data.GenesisID,
		GenesisHash: genesisHashStr,
		Group: groupStr,
		Lease: base64.StdEncoding.EncodeToString(data.Lease[:]),
		RekeyTo: data.RekeyTo.String(),
	}
}

func ProcessInternalTransactionData(data types.Transaction) transaction.InternalTransactionData {
	var processedData = transaction.InternalTransactionData{
		Type:   string(data.Type),
		Header: ProcessHeader(data.Header),
	}

	switch data.Type {
	case types.PaymentTx:
		var output = ProcessPaymentTx(data)
		processedData.PaymentTxnFields = &output
		break
	case types.KeyRegistrationTx:
		var output = ProcessKeyRegistrationTx(data)
		processedData.KeyregTxnFields = &output
		break
	case types.AssetConfigTx:
		var output = ProcessAssetConfigTx(data)
		processedData.AssetConfigTxnFields = &output
		break
	case types.AssetTransferTx:
		var output = ProcessAssetTransferTx(data)
		processedData.AssetTransferTxnFields = &output
		break
	case types.AssetFreezeTx:
		var output = ProcessAssetFreezeTx(data)
		processedData.AssetFreezeTxnFields = &output
		break
	case types.ApplicationCallTx:
		// TODO: finished this
		break
	}

	return  processedData
}

func ProcessSignedTxn(data types.SignedTxn) transaction.SignedTxn {
	return transaction.SignedTxn{
		Sig:      base64.StdEncoding.EncodeToString(data.Sig[:]),
		AuthAddr: data.AuthAddr.String(),
		Txn:      ProcessInternalTransactionData(data.Txn),
	}
}

func ProcessApplyData(data types.ApplyData) transaction.ApplyData {
	return transaction.ApplyData{
		ClosingAmount: uint64(data.ClosingAmount),
		AssetClosingAmount: data.AssetClosingAmount,
		SenderRewards: uint64(data.SenderRewards),
		ReceiverRewards: uint64(data.ReceiverRewards),
		CloseRewards: uint64(data.CloseRewards),
	}
}

func ProcessSignedTxnWithAD(data types.SignedTxnWithAD) transaction.SignedTxnWithAD {
	return transaction.SignedTxnWithAD{
		SignedTxn: ProcessSignedTxn(data.SignedTxn),
		ApplyData: ProcessApplyData(data.ApplyData),
	}
}

func ProcessTransactionInBlock(txn types.SignedTxnInBlock) transaction.Transaction {

	// Process Genesis Hash
	var genesisHash = [32]byte(txn.Txn.GenesisHash)
	var genesisHashStr = base64.StdEncoding.EncodeToString(genesisHash[:])

	var suggestedParams = transaction.SuggestedParams{
		Fee: uint64(txn.Txn.Fee),
		GenesisID: txn.Txn.GenesisID,
		GenesisHash: genesisHashStr,
		FirstRoundValid: uint64(txn.Txn.FirstValid),
		LastRoundValid: uint64(txn.Txn.LastValid),
		Type: string(txn.Txn.Type),
	}

	var transaction = transaction.Transaction{
		HasGenesisID:    txn.HasGenesisID,
		HasGenesisHash:  txn.HasGenesisHash,
		SuggestedParams: suggestedParams,
		SignedTxnWithAD: ProcessSignedTxnWithAD(txn.SignedTxnWithAD),
	}

	return transaction
}
