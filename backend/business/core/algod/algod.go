// Package algod provides the core business API of handling
// everything algod related.
package algod

import (
	"encoding/base32"
	algodv2 "github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/types"
	"go.uber.org/zap"
	"log"
)

type Core struct {
	log *zap.SugaredLogger
	algodClient *algodv2.Client
}

func NewCore(log *zap.SugaredLogger, algodClient *algodv2.Client) Core {
	return Core{
		log: log,
		algodClient: algodClient,
	}
}

// Reference: https://forum.algorand.org/t/finding-blocks-proposer-with-new-api-v2/1778
func rawAddressAsAddress(rawAddress string) types.Address {
	address, err := base32.StdEncoding.DecodeString(rawAddress)
	if err != nil {
		log.Fatalln(err)
	}
	dgst := types.Digest{}
	copy(dgst[:], address[:])

	return types.Address(dgst)
}

// Reference: https://forum.algorand.org/t/finding-blocks-proposer-with-new-api-v2/1778
func byteArrAsAddress(byteArr []byte) types.Address {
	//address, err := base32.StdEncoding.DecodeString(rawAddress)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	dgst := types.Digest{}
	copy(dgst[:], byteArr[:])

	return types.Address(dgst)
}
