package algod

import (
	"encoding/base32"
	"github.com/algorand/go-algorand-sdk/types"
	"log"
)

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
