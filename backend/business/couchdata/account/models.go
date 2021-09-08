package account

import "github.com/algorand/go-algorand-sdk/client/v2/common/models"

type NewAccount struct {
	ID *string `json:"_id"`
	models.Account
	DocType string `json:"doc_type"`
}

type Account struct {
	NewAccount
	ID		string	`json:"_id,omitempty"`
	Rev		string	`json:"_rev,omitempty"`
}
