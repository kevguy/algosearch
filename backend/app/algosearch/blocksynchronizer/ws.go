package blocksynchronizer

import (
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type WsMessage struct {
	Block 				models.Block 	`json:"block"`
	TransactionList 	[]string 		`json:"transaction_ids"`
	AccountList 		[]string  		`json:"account_ids"`
	AssetList 			[]uint64 		`json:"asset_ids"`
	AppList 			[]uint64 		`json:"app_ids"`
	AvgBlockTxnSpeed	float64			`json:"avg_block_txn_speed"`
}
