package algod

import (
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

func removeDuplicateUint64Values(intSlice []uint64) []uint64 {
	keys := make(map[uint64]bool)
	var list []uint64

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func removeDuplicateStrValues(strSlice []string) []string {
	keys := make(map[string]bool)
	var list []string

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func ExtractAssetIdsFromTxn(txn models.Transaction) []uint64 {
	var list []uint64

	if txn.CreatedAssetIndex != 0 {
		list = append(list, txn.CreatedAssetIndex)
	}

	// process ApplicationTransaction: TransactionApplication
	if txn.ApplicationTransaction.ForeignAssets != nil &&
		len(txn.ApplicationTransaction.ForeignAssets) != 0 {
		list = append(list, txn.ApplicationTransaction.ForeignAssets...)
	}

	// process AssetConfigTransaction: TransactionAssetConfig
	if txn.AssetConfigTransaction.AssetId != 0 {
		list = append(list, txn.AssetConfigTransaction.AssetId)
	}

	// process AssetFreezeTransaction: TransactionAssetFreeze
	if txn.AssetFreezeTransaction.AssetId != 0 {
		list = append(list, txn.AssetFreezeTransaction.AssetId)
	}

	// process AssetTransferTransaction: TransactionAssetTransfer
	if txn.AssetTransferTransaction.AssetId != 0 {
		list = append(list, txn.AssetTransferTransaction.AssetId)
	}

	return removeDuplicateUint64Values(list)
}

func ExtractAccountAddrsFromTxn(txn models.Transaction) []string {
	var list []string

	if txn.Sender != "" {
		list = append(list, txn.Sender)
	}
	if txn.AuthAddr != "" {
		list = append(list, txn.AuthAddr)
	}

	// process ApplicationTransaction: TransactionApplication
	if txn.ApplicationTransaction.Accounts != nil &&
		len(txn.ApplicationTransaction.Accounts) != 0 {
		list = append(list, txn.ApplicationTransaction.Accounts...)
	}

	// process AssetConfigTransaction: TransactionAssetConfig
	if txn.AssetConfigTransaction.AssetId != 0 {
		if txn.AssetConfigTransaction.Params.Clawback != "" {
			list = append(list, txn.AssetConfigTransaction.Params.Clawback)
		}
		if txn.AssetConfigTransaction.Params.Creator != "" {
			list = append(list, txn.AssetConfigTransaction.Params.Creator)
		}
		if txn.AssetConfigTransaction.Params.Freeze != "" {
			list = append(list, txn.AssetConfigTransaction.Params.Freeze)
		}
		if txn.AssetConfigTransaction.Params.Manager != "" {
			list = append(list, txn.AssetConfigTransaction.Params.Manager)
		}
		if txn.AssetConfigTransaction.Params.Reserve != "" {
			list = append(list, txn.AssetConfigTransaction.Params.Reserve)
		}
	}

	// process AssetFreezeTransaction: TransactionAssetFreeze
	if txn.AssetFreezeTransaction.Address != "" {
		list = append(list, txn.AssetFreezeTransaction.Address)
	}

	// process AssetTransferTransaction: TransactionAssetTransfer
	if txn.AssetTransferTransaction.AssetId != 0 {
		// TODO: Is this necessary?? Ask Algorand people
		if txn.AssetTransferTransaction.CloseTo != "" {
			list = append(list, txn.AssetTransferTransaction.CloseTo)
		}
		if txn.AssetTransferTransaction.Receiver != "" {
			list = append(list, txn.AssetTransferTransaction.Receiver)
		}
		if txn.AssetTransferTransaction.Sender != "" {
			list = append(list, txn.AssetTransferTransaction.Sender)
		}
	}

	// process PaymentTransaction: TransactionPayment
	if txn.PaymentTransaction.Receiver != "" {
		list = append(list, txn.PaymentTransaction.Receiver)
	}

	return removeDuplicateStrValues(list)
}


func ExtractApplicationIdsFromTxn(txn models.Transaction) []uint64 {
	var list []uint64

	if txn.CreatedApplicationIndex != 0 {
		list = append(list, txn.CreatedApplicationIndex)
	}

	// process ApplicationTransaction: TransactionApplication
	if txn.ApplicationTransaction.ApplicationId != 0  {
		list = append(list, txn.ApplicationTransaction.ApplicationId)
	}

	// TODO: Is this necessary?? Ask Algorand people
	if txn.ApplicationTransaction.ExtraProgramPages != 0  {
		list = append(list, txn.ApplicationTransaction.ExtraProgramPages)
	}

	if txn.ApplicationTransaction.ForeignApps != nil &&
		len(txn.ApplicationTransaction.ForeignApps) != 0 {
		list = append(list, txn.ApplicationTransaction.ForeignApps...)
	}

	return removeDuplicateUint64Values(list)
}
