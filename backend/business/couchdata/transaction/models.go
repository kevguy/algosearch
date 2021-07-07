package transaction

// Transaction is how a signed transaction is encoded in a block.
type Transaction struct {
	SuggestedParams SuggestedParams `json:"suggested_params"`

	SignedTxnWithAD

	HasGenesisID   bool `json:"hgi"`
	HasGenesisHash bool `json:"hgh"`
}

// TransactionRecord represents the data structure of a transaction document.
type TransactionRecord struct {
	Transaction
	ID	string `json:"_id"`
	Rev string `json:"_rev,omitempty"`
}

// SignedTxnWithAD is a (decoded) SignedTxn with associated ApplyData
type SignedTxnWithAD struct {
	SignedTxn
	ApplyData
}

// ApplyData contains information about the transaction's execution.
type ApplyData struct {

	// Closing amount for transaction.
	ClosingAmount uint64 `json:"ca"`

	// Closing amount for asset transaction.
	AssetClosingAmount uint64 `codec:"aca"`

	// Rewards applied to the Sender, Receiver, and CloseRemainderTo accounts.
	SenderRewards   uint64 `json:"rs"`
	ReceiverRewards uint64 `json:"rr"`
	CloseRewards    uint64 `json:"rc"`
	// TODO: handle this
	//EvalDelta       EvalDelta  `json:"dt"`
}

// SignedTxn wraps a transaction and a signature. The encoding of this struct
// is suitable to broadcast on the network
type SignedTxn struct {

	Sig      string	`json:"sig"`
	// TODO: handle this
	//Msig     string `json:"msig"`
	// TODO: handle this
	//Lsig     string `json:"lsig"`
	Txn      InternalTransactionData `json:"txn"`
	AuthAddr string                  `json:"sgnr"`
}

type InternalTransactionData struct {
	// Type of transaction
	Type string `string:"type"`

	// Common fields for all types of transactions
	Header

	// Fields for different types of transactions
	*KeyregTxnFields
	*PaymentTxnFields
	*AssetConfigTxnFields
	*AssetTransferTxnFields
	*AssetFreezeTxnFields
	// TODO" Handle this
	//ApplicationFields
}

// SuggestedParams wraps the transaction parameters common to all transactions,
// typically received from the SuggestedParams endpoint of algod.
// This struct itself is not sent over the wire to or from algod: see models.TransactionParams.
type SuggestedParams struct {
	// Fee is the suggested transaction fee
	// Fee is in units of micro-Algos per byte.
	// Fee may fall to zero but a group of N atomic transactions must
	// still have a fee of at least N*MinTxnFee for the current network protocol.
	Fee uint64 `json:"fee"`

	// Genesis ID
	GenesisID string `json:"genesis-id"`

	// Genesis hash
	GenesisHash string `json:"genesis-hash"`

	// FirstRoundValid is the first protocol round on which the txn is valid
	FirstRoundValid uint64 `json:"first-round"`

	// LastRoundValid is the final protocol round on which the txn may be committed
	LastRoundValid uint64 `json:"last-round"`

	Type string `json:"type"`

	// ConsensusVersion indicates the consensus protocol version
	// as of LastRound.
	//ConsensusVersion string `json:"consensus-version"`

	// FlatFee indicates whether the passed fee is per-byte or per-transaction
	// If true, txn fee may fall below the MinTxnFee for the current network protocol.
	//FlatFee bool `json:"flat-fee"`

	// The minimum transaction fee (not per byte) required for the
	// txn to validate for the current network protocol.
	//MinFee uint64  `json:"min-fee"`
}

// Header captures the fields common to every transaction type.
type Header struct {

	Sender      string    `json:"snd"`
	Fee         uint64 `json:"fee"`
	FirstValid  uint64      `json:"fv"`
	LastValid   uint64      `json:"lv"`
	Note        string     `json:"note"`
	GenesisID   string     `json:"gen"`
	GenesisHash string     `json:"gh"`

	// Group specifies that this transaction is part of a
	// transaction group (and, if so, specifies the hash
	// of a TxGroup).
	Group string `json:"grp"`

	// Lease enforces mutual exclusion of transactions.  If this field is
	// nonzero, then once the transaction is confirmed, it acquires the
	// lease identified by the (Sender, Lease) pair of the transaction until
	// the LastValid round passes.  While this transaction possesses the
	// lease, no other transaction specifying this lease can be confirmed.
	Lease string `json:"lx"`

	// RekeyTo, if nonzero, sets the sender's SpendingKey to the given address
	// If the RekeyTo address is the sender's actual address, the SpendingKey is set to zero
	// This allows "re-keying" a long-lived account -- rotating the signing key, changing
	// membership of a multisig account, etc.
	RekeyTo string `json:"rekey"`
}

// KeyregTxnFields captures the fields used for key registration transactions.
type KeyregTxnFields struct {
	VotePK           string `json:"votekey"`
	SelectionPK      string  `json:"selkey"`
	VoteFirst        uint64  `json:"votefst"`
	VoteLast         uint64  `json:"votelst"`
	VoteKeyDilution  uint64 `json:"votekd"`
	Nonparticipation bool   `json:"nonpart"`
}

// PaymentTxnFields captures the fields used by payment transactions.
type PaymentTxnFields struct {
	Receiver string    `json:"rcv"`
	Amount   uint64 `json:"amt"`

	// When CloseRemainderTo is set, it indicates that the
	// transaction is requesting that the account should be
	// closed, and all remaining funds be transferred to this
	// address.
	CloseRemainderTo string `json:"close"`
}

// AssetParams describes the parameters of an asset.
type AssetParams struct {
	// Total specifies the total number of units of this asset
	// created.
	Total uint64 `json:"t"`

	// Decimals specifies the number of digits to display after the decimal
	// place when displaying this asset. A value of 0 represents an asset
	// that is not divisible, a value of 1 represents an asset divisible
	// into tenths, and so on. This value must be between 0 and 19
	// (inclusive).
	Decimals uint32 `json:"dc"`

	// DefaultFrozen specifies whether slots for this asset
	// in user accounts are frozen by default or not.
	DefaultFrozen bool `json:"df"`

	// UnitName specifies a hint for the name of a unit of
	// this asset.
	UnitName string `json:"un"`

	// AssetName specifies a hint for the name of the asset.
	AssetName string `json:"an"`

	// URL specifies a URL where more information about the asset can be
	// retrieved
	URL string `json:"au"`

	// MetadataHash specifies a commitment to some unspecified asset
	// metadata. The format of this metadata is up to the application.
	MetadataHash string `json:"am"`

	// Manager specifies an account that is allowed to change the
	// non-zero addresses in this AssetParams.
	Manager string `json:"m"`

	// Reserve specifies an account whose holdings of this asset
	// should be reported as "not minted".
	Reserve string `json:"r"`

	// Freeze specifies an account that is allowed to change the
	// frozen state of holdings of this asset.
	Freeze string `json:"f"`

	// Clawback specifies an account that is allowed to take units
	// of this asset from any account.
	Clawback string `json:"c"`
}

// AssetConfigTxnFields captures the fields used for asset
// allocation, re-configuration, and destruction.
type AssetConfigTxnFields struct {
	// ConfigAsset is the asset being configured or destroyed.
	// A zero value means allocation.
	ConfigAsset uint64 `json:"caid"`

	// AssetParams are the parameters for the asset being
	// created or re-configured.  A zero value means destruction.
	AssetParams AssetParams `json:"apar"`
}

// AssetTransferTxnFields captures the fields used for asset transfers.
type AssetTransferTxnFields struct {
	XferAsset uint64 `json:"xaid"`

	// AssetAmount is the amount of asset to transfer.
	// A zero amount transferred to self allocates that asset
	// in the account's Assets map.
	AssetAmount uint64 `json:"aamt"`

	// AssetSender is the sender of the transfer.  If this is not
	// a zero value, the real transaction sender must be the Clawback
	// address from the AssetParams.  If this is the zero value,
	// the asset is sent from the transaction's Sender.
	AssetSender string `json:"asnd"`

	// AssetReceiver is the recipient of the transfer.
	AssetReceiver string `json:"arcv"`

	// AssetCloseTo indicates that the asset should be removed
	// from the account's Assets map, and specifies where the remaining
	// asset holdings should be transferred.  It's always valid to transfer
	// remaining asset holdings to the creator account.
	AssetCloseTo string `json:"aclose"`
}

// AssetFreezeTxnFields captures the fields used for freezing asset slots.
type AssetFreezeTxnFields struct {

	// FreezeAccount is the address of the account whose asset
	// slot is being frozen or un-frozen.
	FreezeAccount string `json:"fadd"`

	// FreezeAsset is the asset ID being frozen or un-frozen.
	FreezeAsset uint64 `json:"faid"`

	// AssetFrozen is the new frozen value.
	AssetFrozen bool `json:"afrz"`
}
