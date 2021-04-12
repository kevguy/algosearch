package db

import (
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// NewBlock represents the data structure for constructing a new block data
type NewBlock struct {
	models.Block
	Proposer        string  `json:"proposer"`
	BlockHash       string  `json:"block-hash"`
}

type NewBlockDoc struct {
	NewBlock
	DocType			string	`json:"doc_type"`
}

// Block represents the data structure of a block document.
type Block struct {
	NewBlockDoc
	ID	string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
}

