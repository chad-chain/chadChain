package types

import (
	"golang.org/x/crypto/sha3"
)

func Keccak256(data ...[]byte) []byte {
	hashState := sha3.NewLegacyKeccak256()
	for _, input := range data {
		hashState.Write(input)
	}

	return hashState.Sum(nil)
}
