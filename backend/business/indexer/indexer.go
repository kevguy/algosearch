package indexer





//func Fuck(ctx context.Context, indexerClient *indexerv2.Client, roundNum uint64) (models.Block, error) {
//}

//func ConvertBlockJSON(ctx context.Context, jsonBlock models.Block) (block.NewBlock, error) {
//
//	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "indexer.ConvertBlockJSON")
//	defer span.End()
//
//	var genesisHashStr = base64.StdEncoding.EncodeToString(jsonBlock.GenesisHash[:])
//
//	var newBLock = block.NewBlock{
//		GenesisHash:        genesisHashStr,
//		GenesisID:          jsonBlock.GenesisId,
//		PrevBlockHash:      "",
//		Rewards:            block.Rewards{},
//		Round:              0,
//		Seed:               "",
//		Timestamp:          0,
//		Transactions:       nil,
//		TransactionsRoot:   "",
//		TransactionCounter: 0,
//		UpgradeState:       block.UpgradeState{},
//		UpgradeVote:        block.UpgradeVote{},
//		Proposer:           "",
//		BlockHash:          "",
//	}
//}
