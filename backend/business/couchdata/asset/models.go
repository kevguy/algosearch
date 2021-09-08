package asset

import "github.com/algorand/go-algorand-sdk/client/v2/common/models"

type NewAsset struct {
	ID *string `json:"_id"`
	models.Asset
	DocType string	`json:"doc_type"`
}

type Asset struct {
	NewAsset
	ID		string	`json:"_id,omitempty"`
	Rev		string	`json:"_rev,omitempty"`
}
