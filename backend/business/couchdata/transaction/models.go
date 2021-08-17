package transaction

import "github.com/algorand/go-algorand-sdk/client/v2/common/models"

type Transaction struct {
	models.Transaction
	ID	string `json:"_id"`
	Rev string `json:"_rev,omitempty"`
}
