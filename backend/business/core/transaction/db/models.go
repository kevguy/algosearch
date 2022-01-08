package db

import "github.com/algorand/go-algorand-sdk/client/v2/common/models"

type NewTransaction struct {
	ID *string							`json:"_id"`
	models.Transaction
	DocType					string		`json:"doc_type"`
	AssociatedAccounts		[]string	`json:"associated_accounts"`
	AssociatedApplications	[]uint64	`json:"associated_applications"`
	AssociatedAssets		[]uint64	`json:"associated_assets"`
}

type Transaction struct {
	NewTransaction
	ID		string	`json:"_id,omitempty"`
	Rev		string	`json:"_rev,omitempty"`
}
