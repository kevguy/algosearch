package algod

import (
	"encoding/base64"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/types"
	"sort"
	"unicode"
	"unicode/utf8"
)

//func ProcessPaymentTx(txn types.Transaction) transaction.PaymentTxnFields {
//	return transaction.PaymentTxnFields{
//		Receiver: txn.Receiver.String(),
//		Amount: uint64(txn.Amount),
//		CloseRemainderTo: txn.CloseRemainderTo.String(),
//	}
//}
//
//func ProcessKeyRegistrationTx(txn types.Transaction) transaction.KeyregTxnFields {
//	return transaction.KeyregTxnFields{
//		VotePK: base64.StdEncoding.EncodeToString(txn.VotePK[:]),
//		SelectionPK: base64.StdEncoding.EncodeToString(txn.SelectionPK[:]),
//		VoteFirst: uint64(txn.VoteFirst),
//		VoteLast: uint64(txn.VoteLast),
//		VoteKeyDilution: txn.VoteKeyDilution,
//		Nonparticipation: txn.Nonparticipation,
//	}
//}
//
//func ProcessAssetTransferTx(txn types.Transaction) transaction.AssetTransferTxnFields {
//	return transaction.AssetTransferTxnFields{
//		XferAsset: uint64(txn.XferAsset),
//		AssetAmount: txn.AssetAmount,
//		AssetSender: txn.AssetSender.String(),
//		AssetReceiver: txn.AssetReceiver.String(),
//		AssetCloseTo: txn.AssetCloseTo.String(),
//	}
//}
//
//func ProcessAssetConfigTx(txn types.Transaction) transaction.AssetConfigTxnFields {
//	return transaction.AssetConfigTxnFields{
//		ConfigAsset: uint64(txn.ConfigAsset),
//		AssetParams: transaction.AssetParams{
//			Total: txn.AssetParams.Total,
//			Decimals: txn.AssetParams.Decimals,
//			DefaultFrozen: txn.AssetParams.DefaultFrozen,
//			UnitName: txn.AssetParams.UnitName,
//			URL: txn.AssetParams.URL,
//			MetadataHash: base64.StdEncoding.EncodeToString(txn.AssetParams.MetadataHash[:]),
//			Manager: txn.AssetParams.Manager.String(),
//			Reserve: txn.AssetParams.Reserve.String(),
//			Freeze: txn.AssetParams.Freeze.String(),
//			Clawback: txn.AssetParams.Clawback.String(),
//		},
//	}
//}
//
//func ProcessAssetFreezeTx(txn types.Transaction) transaction.AssetFreezeTxnFields {
//	return transaction.AssetFreezeTxnFields{
//		FreezeAccount: txn.FreezeAccount.String(),
//		FreezeAsset: uint64(txn.FreezeAsset),
//		AssetFrozen: txn.AssetFrozen,
//	}
//}
//
//func ProcessHeader(data types.Header) transaction.Header {
//
//	var genesisHash = [32]byte(data.GenesisHash)
//	var genesisHashStr = base64.StdEncoding.EncodeToString(genesisHash[:])
//
//	var group = [32]byte(data.Group)
//	var groupStr = base64.StdEncoding.EncodeToString(group[:])
//
//	return transaction.Header{
//		Sender: data.Sender.String(),
//		Fee: uint64(data.Fee),
//		FirstValid: uint64(data.FirstValid),
//		LastValid: uint64(data.LastValid),
//		Note: base64.StdEncoding.EncodeToString(data.Note),
//		GenesisID: data.GenesisID,
//		GenesisHash: genesisHashStr,
//		Group: groupStr,
//		Lease: base64.StdEncoding.EncodeToString(data.Lease[:]),
//		RekeyTo: data.RekeyTo.String(),
//	}
//}
//
//func ProcessInternalTransactionData(data types.Transaction) transaction.InternalTransactionData {
//	var processedData = transaction.InternalTransactionData{
//		Type:   string(data.Type),
//		Header: ProcessHeader(data.Header),
//	}
//
//	switch data.Type {
//	case types.PaymentTx:
//		var output = ProcessPaymentTx(data)
//		processedData.PaymentTxnFields = &output
//		break
//	case types.KeyRegistrationTx:
//		var output = ProcessKeyRegistrationTx(data)
//		processedData.KeyregTxnFields = &output
//		break
//	case types.AssetConfigTx:
//		var output = ProcessAssetConfigTx(data)
//		processedData.AssetConfigTxnFields = &output
//		break
//	case types.AssetTransferTx:
//		var output = ProcessAssetTransferTx(data)
//		processedData.AssetTransferTxnFields = &output
//		break
//	case types.AssetFreezeTx:
//		var output = ProcessAssetFreezeTx(data)
//		processedData.AssetFreezeTxnFields = &output
//		break
//	case types.ApplicationCallTx:
//		// TODO: finished this
//		break
//	}
//
//	return  processedData
//}
//
//func ProcessSignedTxn(data types.SignedTxn) transaction.SignedTxn {
//	return transaction.SignedTxn{
//		Sig:      base64.StdEncoding.EncodeToString(data.Sig[:]),
//		AuthAddr: data.AuthAddr.String(),
//		Txn:      ProcessInternalTransactionData(data.Txn),
//	}
//}
//
//func ProcessApplyData(data types.ApplyData) transaction.ApplyData {
//	return transaction.ApplyData{
//		ClosingAmount: uint64(data.ClosingAmount),
//		AssetClosingAmount: data.AssetClosingAmount,
//		SenderRewards: uint64(data.SenderRewards),
//		ReceiverRewards: uint64(data.ReceiverRewards),
//		CloseRewards: uint64(data.CloseRewards),
//	}
//}
//
//func ProcessSignedTxnWithAD(data types.SignedTxnWithAD) transaction.SignedTxnWithAD {
//	return transaction.SignedTxnWithAD{
//		SignedTxn: ProcessSignedTxn(data.SignedTxn),
//		ApplyData: ProcessApplyData(data.ApplyData),
//	}
//}
//
//func ProcessTransactionInBlock(txn types.SignedTxnInBlock) transaction.Transaction {
//
//	// Process Genesis Hash
//	var genesisHash = [32]byte(txn.Txn.GenesisHash)
//	var genesisHashStr = base64.StdEncoding.EncodeToString(genesisHash[:])
//
//	var suggestedParams = transaction.SuggestedParams{
//		Fee: uint64(txn.Txn.Fee),
//		GenesisID: txn.Txn.GenesisID,
//		GenesisHash: genesisHashStr,
//		FirstRoundValid: uint64(txn.Txn.FirstValid),
//		LastRoundValid: uint64(txn.Txn.LastValid),
//		Type: string(txn.Txn.Type),
//	}
//
//	var transaction = transaction.Transaction{
//		HasGenesisID:    txn.HasGenesisID,
//		HasGenesisHash:  txn.HasGenesisHash,
//		SuggestedParams: suggestedParams,
//		SignedTxnWithAD: ProcessSignedTxnWithAD(txn.SignedTxnWithAD),
//	}
//
//	return transaction
//}

////////////////////////////////////////////////////////////

// PrintableUTF8OrEmpty checks to see if the entire string is a UTF8 printable string.
// If this is the case, the string is returned as is. Otherwise, the empty string is returned.
// https://github.com/algorand/indexer/blob/5ad47734a19f0ff319c7ae852053f45bfc226475/util/util.go#L13
func PrintableUTF8OrEmpty(in string) string {
	// iterate throughout all the characters in the string to see if they are all printable.
	// when range iterating on go strings, go decode each element as a utf8 rune.
	for _, c := range in {
		// is this a printable character, or invalid rune ?
		if c == utf8.RuneError || !unicode.IsPrint(c) {
			return ""
		}
	}
	return in
}

func extractPaymentTx(txn types.SignedTxnWithAD) models.TransactionPayment {
	return models.TransactionPayment{
		Amount:           uint64(txn.Txn.Amount),
		CloseAmount:      uint64(txn.ClosingAmount),
		CloseRemainderTo: txn.Txn.CloseRemainderTo.String(),
		Receiver:         txn.Txn.Receiver.String(),
	}
}

func extractKeyRegistrationTx(txn types.SignedTxnWithAD) models.TransactionKeyreg {
	return models.TransactionKeyreg{
		NonParticipation:          txn.Txn.Nonparticipation,
		SelectionParticipationKey: txn.Txn.SelectionPK[:],
		VoteFirstValid:            uint64(txn.Txn.VoteFirst),
		VoteKeyDilution:           txn.Txn.VoteKeyDilution,
		VoteLastValid:             uint64(txn.Txn.VoteLast),
		VoteParticipationKey:      txn.Txn.VotePK[:],
	}
}

func extractAssetConfigTx(txn types.SignedTxnWithAD) models.TransactionAssetConfig {
	assetParams := models.AssetParams{
		Clawback:      txn.Txn.AssetParams.Clawback.String(),
		Creator:       txn.Txn.Sender.String(),
		Decimals:      uint64(txn.Txn.AssetParams.Decimals),
		DefaultFrozen: txn.Txn.AssetParams.DefaultFrozen,
		Freeze:        txn.Txn.AssetParams.Freeze.String(),
		Manager:       txn.Txn.AssetParams.Manager.String(),
		MetadataHash:  txn.Txn.AssetParams.MetadataHash[:],
		Name:          PrintableUTF8OrEmpty(txn.Txn.AssetParams.AssetName),
		Reserve:       txn.Txn.AssetParams.Reserve.String(),
		Total:         txn.Txn.AssetParams.Total,
		UnitName:      PrintableUTF8OrEmpty(txn.Txn.AssetParams.UnitName),
		Url:           PrintableUTF8OrEmpty(txn.Txn.AssetParams.URL),
	}

	return models.TransactionAssetConfig{
		AssetId: uint64(txn.Txn.ConfigAsset),
		Params: assetParams,
	}
}

func extractAssetTransferTx(txn types.SignedTxnWithAD) models.TransactionAssetTransfer {
	return models.TransactionAssetTransfer{
		Amount:      txn.Txn.AssetAmount,
		AssetId:     uint64(txn.Txn.XferAsset),
		CloseAmount: txn.AssetClosingAmount,
		CloseTo:     txn.Txn.AssetCloseTo.String(),
		Receiver:    txn.Txn.AssetReceiver.String(),
		Sender:      txn.Txn.AssetSender.String(),
	}
}

func extractAssetFreezeTx(txn types.SignedTxnWithAD) models.TransactionAssetFreeze {
	return models.TransactionAssetFreeze{
		Address:         txn.Txn.FreezeAccount.String(),
		AssetId:         uint64(txn.Txn.FreezeAsset),
		NewFreezeStatus: txn.Txn.AssetFrozen,
	}
}

// https://github.com/algorand/indexer/blob/5ad47734a19f0ff319c7ae852053f45bfc226475/api/converter_utils.go#L213
func onCompletionToTransactionOnCompletion(oc types.OnCompletion) string {
	switch oc {
	case types.NoOpOC:
		return "noop"
	case types.OptInOC:
		return "optin"
	case types.CloseOutOC:
		return "closeout"
	case types.ClearStateOC:
		return "clear"
	case types.UpdateApplicationOC:
		return "update"
	case types.DeleteApplicationOC:
		return "delete"
	}
	return "unknown"
}

func extractApplicationTx(txn types.SignedTxnWithAD) models.TransactionApplication {

	//args := make([][]byte, 0)
	//for _, v := range txn.Txn.ApplicationArgs {
	//	args = append(args, v)
	//}

	accts := make([]string, 0)
	for _, v := range txn.Txn.Accounts {
		accts = append(accts, v.String())
	}

	apps := make([]uint64, 0)
	for _, v := range txn.Txn.ForeignApps {
		apps = append(apps, uint64(v))
	}

	assets := make([]uint64, 0)
	for _, v := range txn.Txn.ForeignAssets {
		assets = append(apps, uint64(v))
	}

	a := models.TransactionApplication{
		Accounts:          accts,
		ApplicationArgs:   txn.Txn.ApplicationArgs,
		ApplicationId:     uint64(txn.Txn.ApplicationID),
		ApprovalProgram:   txn.Txn.ApprovalProgram,
		ClearStateProgram: txn.Txn.ClearStateProgram,
		ExtraProgramPages: uint64(txn.Txn.ExtraProgramPages),
		ForeignApps:       apps,
		ForeignAssets:     assets,
		GlobalStateSchema: models.StateSchema{},
		LocalStateSchema:  models.StateSchema{},
		OnCompletion:      onCompletionToTransactionOnCompletion(txn.Txn.OnCompletion),
	}
	return a
}

// TODO: Replace with lsig.Blank() when that gets merged into go-algorand-sdk
// https://github.com/algorand/indexer/blob/6e4d737f2e4e49088b436a234caee6681435053d/api/converter_utils.go#L177
func isBlank(lsig types.LogicSig) bool {
	if lsig.Args != nil {
		return false
	}
	if len(lsig.Logic) != 0 {
		return false
	}
	if !lsig.Msig.Blank() {
		return false
	}
	if lsig.Sig != (types.Signature{}) {
		return false
	}
	return true
}

////////////////////////////////////////////////////
// Helpers to convert to and from generated types //
////////////////////////////////////////////////////

// https://github.com/algorand/indexer/blob/5ad47734a19f0ff319c7ae852053f45bfc226475/api/converter_utils.go#L146
func sigToTransactionSig(sig types.Signature) []byte {
	if sig == (types.Signature{}) {
		return nil
	}

	tsig := sig[:]
	return tsig
}

// https://github.com/algorand/indexer/blob/5ad47734a19f0ff319c7ae852053f45bfc226475/api/converter_utils.go#L155
func msigToTransactionMsig(msig types.MultisigSig) *models.TransactionSignatureMultisig {
	if msig.Blank() {
		return nil
	}

	subsigs := make([]models.TransactionSignatureMultisigSubsignature, 0)
	for _, subsig := range msig.Subsigs {
		signature := sigToTransactionSig(subsig.Sig)
		subsigs = append(subsigs, models.TransactionSignatureMultisigSubsignature{
			//PublicKey: bytePtr(subsig.Key[:]),
			PublicKey: subsig.Key,
			Signature: signature,
		})
	}

	ret := models.TransactionSignatureMultisig{
		Subsignature: subsigs,
		Threshold:    uint64(msig.Threshold),
		Version:      uint64(msig.Version),
	}
	return &ret
}

// https://github.com/algorand/indexer/blob/master/api/converter_utils.go#L193
func lsigToTransactionLsig(lsig types.LogicSig) *models.TransactionSignatureLogicsig {
	fmt.Println("fuck you ")
	if isBlank(lsig) {
		return nil
	}

	//args := make([]string, 0)
	//for _, arg := range lsig.Args {
	//	args = append(args, base64.StdEncoding.EncodeToString(arg))
	//}

	fmt.Println("fuck you too")
	txnMSig := msigToTransactionMsig(lsig.Msig)
	if txnMSig == nil {
		txnMSig = &models.TransactionSignatureMultisig{
			Subsignature: nil,
			Threshold:    0,
			Version:      0,
		}
	}
	ret := models.TransactionSignatureLogicsig{
		Args:              lsig.Args,
		Logic:             lsig.Logic,
		MultisigSignature: *txnMSig,
		Signature:         sigToTransactionSig(lsig.Sig),
	}
	return &ret
}

// https://github.com/algorand/indexer/blob/6e4d737f2e4e49088b436a234caee6681435053d/api/converter_utils.go#L233
// The state delta bits need to be sorted for testing. Maybe it would be
// for end users too, people always seem to notice results changing.
func stateDeltaToStateDelta(d types.StateDelta) []models.EvalDeltaKeyValue {
	if len(d) == 0 {
		return nil
	}
	var delta []models.EvalDeltaKeyValue
	keys := make([]string, 0)
	for k := range d {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := d[k]
		delta = append(delta, models.EvalDeltaKeyValue{
			Key: base64.StdEncoding.EncodeToString([]byte(k)),
			Value: models.EvalDelta{
				Action: uint64(v.Action),
				//Bytes:  strPtr(base64.StdEncoding.EncodeToString(v.Bytes)),
				Bytes:  v.Bytes,
				Uint:   v.Uint,
			},
		})
	}
	return delta
}

// https://github.com/algorand/indexer/blob/6e4d737f2e4e49088b436a234caee6681435053d/api/handlers.go
// https://github.com/algorand/indexer/blob/fac3b03349d108457abc27c083bd44052590c487/importer/importer.go
// https://github.com/algorand/indexer/blob/6e4d737f2e4e49088b436a234caee6681435053d/api/converter_utils.go
func ProcessTransactionInBlock(txn types.SignedTxnInBlock, blockInfo types.Block) models.Transaction {


	var genesisHash = [32]byte(txn.Txn.GenesisHash)
	var genesisHashStr = base64.StdEncoding.EncodeToString(genesisHash[:])
	fmt.Println("Indexer here")
	fmt.Println("- Genesis Hash: " + genesisHashStr)
	fmt.Println("- ID: ")
	fmt.Println("Bye")

	var payment *models.TransactionPayment
	var keyreg *models.TransactionKeyreg
	var assetConfig *models.TransactionAssetConfig
	var assetFreeze *models.TransactionAssetFreeze
	var assetTransfer *models.TransactionAssetTransfer
	var application *models.TransactionApplication

	//var group = [32]byte(txn.Txn.Group)
	//var groupStr = base64.StdEncoding.EncodeToString(group[:])

	switch txn.Txn.Type {
	case types.PaymentTx:
		paymentTx := extractPaymentTx(txn.SignedTxnWithAD)
		payment = &paymentTx
	case types.KeyRegistrationTx:
		keyRegTx := extractKeyRegistrationTx(txn.SignedTxnWithAD)
		keyreg = &keyRegTx
	case types.AssetConfigTx:
		assetConfigTx := extractAssetConfigTx(txn.SignedTxnWithAD)
		assetConfig = & assetConfigTx
	case types.AssetTransferTx:
		assetTransferTx := extractAssetTransferTx(txn.SignedTxnWithAD)
		assetTransfer = &assetTransferTx
	case types.AssetFreezeTx:
		assetFreezeTx := extractAssetFreezeTx(txn.SignedTxnWithAD)
		assetFreeze = &assetFreezeTx
	case types.ApplicationCallTx:
		applicationTx := extractApplicationTx(txn.SignedTxnWithAD)
		application = &applicationTx
	}

	sig := models.TransactionSignature{}

	logicSig := lsigToTransactionLsig(txn.SignedTxnWithAD.SignedTxn.Lsig)
	multiSig := msigToTransactionMsig(txn.SignedTxnWithAD.SignedTxn.Msig)
	sigsig := sigToTransactionSig(txn.SignedTxnWithAD.Sig)

	if logicSig != nil {
		sig.Logicsig = *logicSig
	}
	if multiSig != nil {
		sig.Multisig = *multiSig
	}
	if sigsig != nil {
		sig.Sig = sigsig
	}

	//sig := models.TransactionSignature{
	//	Logicsig: *logicSig,
	//	Multisig: *multiSig,
	//	Sig:      *sigsig,
	//}

	var localStateDelta []models.AccountStateDelta
	type tuple struct {
		key		uint64
		address	types.Address
	}
	if len(txn.ApplyData.EvalDelta.LocalDeltas) > 0 {
		keys := make([]tuple, 0)
		for k := range txn.ApplyData.EvalDelta.LocalDeltas {
			if k == 0 {
				keys = append(keys, tuple{
					key:     0,
					address: txn.Txn.Sender,
				})
			} else {
				addr := types.Address{}
				copy(addr[:], txn.Txn.Accounts[k-1][:])
				keys = append(keys, tuple{
					key:     k,
					address: addr,
				})
			}
		}
		sort.Slice(keys, func(i, j int) bool { return keys[i].key < keys[j].key })
		d := make([]models.AccountStateDelta, 0)
		for _, k := range keys {
			v := txn.ApplyData.EvalDelta.LocalDeltas[k.key]
			delta := stateDeltaToStateDelta(v)
			if delta != nil {
				d = append(d, models.AccountStateDelta{
					Address: k.address.String(),
					Delta:   delta,
				})
			}
		}
		localStateDelta = d
	}

	// TODO:
	txn.Txn.GenesisHash = blockInfo.GenesisHash
	txn.Txn.GenesisID = blockInfo.GenesisID

	var transaction = models.Transaction{
		//ApplicationTransaction:   *application,
		//AssetConfigTransaction:   *assetConfig,
		//AssetFreezeTransaction:   *assetFreeze,
		//AssetTransferTransaction: *assetTransfer,
		//PaymentTransaction:       *payment,
		//KeyregTransaction:        *keyreg,

		ClosingAmount:            uint64(txn.ClosingAmount),
		ConfirmedRound:           uint64(blockInfo.Round),
		// TODO: ask Algorand people on how to support this
		// I can't find it anywhere.
		IntraRoundOffset:         0.,
		// TODO: ask Algorand to verify if it's really this one
		RoundTime:                uint64(blockInfo.TimeStamp),
		Fee:                      uint64(txn.Txn.Fee),
		FirstValid:               uint64(txn.Txn.FirstValid),
		// TODO: this is because ... Kevin! Write down what you got from Jason Paulos
		//GenesisHash:              txn.Txn.GenesisHash[:],
		GenesisHash:			  blockInfo.GenesisHash[:],
		// TODO: this is because ... Kevin! Write down what you got from Jason Paulos
		//GenesisId:                txn.Txn.GenesisID,
		GenesisId: 				  blockInfo.GenesisID,
		Group:					  txn.Txn.Group[:],
		LastValid:                uint64(txn.Txn.LastValid),
		// TODO:
		Lease:                    txn.Txn.Lease[:],
		Note:                     txn.Txn.Note[:],
		Sender:                   txn.Txn.Sender.String(),
		ReceiverRewards:          uint64(txn.ReceiverRewards),
		CloseRewards:             uint64(txn.CloseRewards),
		SenderRewards:            uint64(txn.SenderRewards),
		Type:                     string(txn.Txn.Type),
		Signature:                sig,
		Id:   					  crypto.TransactionIDString(txn.Txn),
		// TODO
		RekeyTo:          txn.Txn.RekeyTo.String(),
		GlobalStateDelta: stateDeltaToStateDelta(txn.EvalDelta.GlobalDelta),
		LocalStateDelta:  localStateDelta,
		AuthAddr:         txn.AuthAddr.String(),
	}

	switch txn.Txn.Type {
	case types.PaymentTx:
		transaction.PaymentTransaction = *payment
	case types.KeyRegistrationTx:
		transaction.KeyregTransaction = *keyreg
	case types.AssetConfigTx:
		transaction.AssetConfigTransaction = *assetConfig
	case types.AssetTransferTx:
		transaction.AssetTransferTransaction = *assetTransfer
	case types.AssetFreezeTx:
		transaction.AssetFreezeTransaction = *assetFreeze
	case types.ApplicationCallTx:
		transaction.ApplicationTransaction = *application
	}

	if txn.Txn.Type == types.AssetConfigTx {
		if assetConfig != nil && assetConfig.AssetId == 0 {
			transaction.CreatedAssetIndex = uint64(txn.Txn.AssetConfigTxnFields.ConfigAsset)
		}
	}

	if txn.Txn.Type == types.ApplicationCallTx {
		if application != nil && application.ApplicationId == 0 {
			transaction.CreatedApplicationIndex = uint64(txn.Txn.ApplicationCallTxnFields.ApplicationID)
		}
	}

	return transaction
}
