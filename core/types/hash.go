package types

import (
	"golang.org/x/crypto/sha3"
)

// Keccak256 is the Keccak-256 hashing method
type Keccak256 struct{}

// New creates a new Keccak-256 hashing method
func NewKeccak256() *Keccak256 {
	return &Keccak256{}
}

// Hash generates a Keccak-256 hash from a byte array
func (h *Keccak256) Hash(data []byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(data)
	return hash.Sum(nil)
}
