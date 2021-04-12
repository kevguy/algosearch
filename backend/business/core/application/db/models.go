package db

import "github.com/algorand/go-algorand-sdk/client/v2/common/models"

type NewApplication struct {
	ID *string `json:"_id"`
	models.Application
	DocType string	`json:"doc_type"`
}

type Application struct {
	NewApplication
	ID		string	`json:"_id,omitempty"`
	Rev		string	`json:"_rev,omitempty"`
}
