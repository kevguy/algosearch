package block

import (
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

//// Rewards represents the rewards of the block
//type Rewards struct {
//	FeeSink			string `json:"fee-sink"`
//	RewardsCalRound	uint64 `json:"rewards-calculation-round"`
//	RewardsLevel	uint64 `json:"rewards-level"`
//	RewardsPool		string `json:"rewards-pool"`
//	RewardsRate		uint64 `json:"rewards-rate"`
//	RewardsResidue	uint64 `json:"rewards-residue"`
//}
//
//// UpgradeState represents the upgrade state of the block
//type UpgradeState struct {
//	CurrProtocol			string `json:"current-protocol"`
//	NextProtocol			*string `json:"next-protocol,omitempty"`
//	NextProtocolApprovals	uint64	`json:"next-protocol-approvals"`
//	NextProtocolSwitchOn	uint64	`json:"next-protocol-switch-on"`
//	NextProtocolVoteBefore	uint64	`json:"next-protocol-vote-before"`
//}
//
//// UpgradeVote represents the upgrade vote of the block
//type UpgradeVote struct {
//	UpgradeApprove	bool	`json:"upgrade-approve"`
//	UpgradeDelay	uint64	`json:"upgrade-delay"`
//	UpgradePropose	*string	`json:"upgrade-propose,omitempty"`
//}
//
//type SignedTransaction struct {
//	Signature string `json:"signature"`
//}
//
//// TransactionInBlock represents the data structure of a transaction that
//// is store in the block.
//type TransactionInBlock struct {
//	HasGenesisID  	bool `json:"has-genesis-id"`
//	HasGenesisHash	bool `json:"has-genesis-hash"`
//}

type NewBlock struct {
	models.Block
	Proposer           string                    `json:"proposer"`
	BlockHash          string                    `json:"block-hash"`
}

// NewBlock represents the data structure for constructing a new block
//type NewBlock struct {
//	GenesisHash        string                    `json:"genesis-hash"`
//	GenesisID          string                    `json:"genesis-id"`
//	PrevBlockHash      string                    `json:"previous-block-hash"`
//	Rewards            Rewards                   `json:"rewards"`
//	Round              uint64                    `json:"round"`
//	Seed               string                    `json:"seed"`
//	Timestamp          uint64                    `json:"timestamp"`
//	Transactions       []transaction.Transaction `json:"transactions"`
//	TransactionsRoot   string                    `json:"transactions-root"`
//	TransactionCounter uint64                    `json:"txn-counter"`
//	UpgradeState       UpgradeState              `json:"upgrade-state"`
//	UpgradeVote        UpgradeVote               `json:"upgrade-vote"`
//	Proposer           string                    `json:"proposer"`
//	BlockHash          string                    `json:"block-hash"`
//}

// Block represents the data structure of a block document.
type Block struct {
	NewBlock
	ID	string `json:"_id"`
	Rev string `json:"_rev,omitempty"`
}

