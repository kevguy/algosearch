package transaction

import "github.com/algorand/go-algorand-sdk/client/v2/common/models"

type NewTransaction struct {
	ID *string `json:"_id"`
	models.Transaction
	DocType string	`json:"doc_type"`
}

type Transaction struct {
	NewTransaction
	ID		string	`json:"_id,omitempty"`
	Rev		string	`json:"_rev,omitempty"`
}
