package types

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

func Keccak256(data ...[]byte) []byte {
	hashState := sha3.NewLegacyKeccak256()
	for _, input := range data {
		hashState.Write(input)
	}

	return hashState.Sum(nil)
}

type Tx struct {
	To    [20]byte
	Value uint64
	Nonce uint64
}

type SignedTx struct {
	To      [20]byte
	Value   uint64
	Nonce   uint64
	V, R, S *big.Int // signature values
}

// signing over header and verifying it

func SignTransaction(tx *Tx, privateKey *ecdsa.PrivateKey) (SignedTx, error) {
	h := Hash(tx)
	sig, err := crypto.Sign(h[:], privateKey)
	if err != nil {
		panic(err)
	}
	R, S, V := decodeSignature(sig)
	// Create a signed transaction by V, R, S values
	signedTx := SignedTx{
		To:    tx.To,
		Value: tx.Value,
		Nonce: tx.Nonce,
		V:     V,
		R:     R,
		S:     S,
	}
	return signedTx, nil
}

// recoverPlain recovers the address which has signed the given data using the v, r, and s values
func recoverPlain(sighash common.Hash, R, S, Vb *big.Int, homestead bool) (common.Address, error) {
	if Vb.BitLen() > 8 {
		// return common.Address{}, ErrInvalidSig
		panic("invalid signature")
	}
	V := byte(Vb.Uint64() - 27)
	if !crypto.ValidateSignatureValues(V, R, S, homestead) {
		// return common.Address{}, ErrInvalidSig
		panic("invalid signature")
	}
	// encode the signature in uncompressed format
	r, s := R.Bytes(), S.Bytes()
	sig := make([]byte, crypto.SignatureLength)
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[64] = V
	// recover the public key from the signature
	pub, err := crypto.Ecrecover(sighash[:], sig)
	if err != nil {
		return common.Address{}, err
	}
	if len(pub) == 0 || pub[0] != 4 {
		return common.Address{}, errors.New("invalid public key")
	}
	var addr common.Address
	copy(addr[:], crypto.Keccak256(pub[1:])[12:])
	return addr, nil
}

// decodeSignature decodes the signature into v, r, and s values
func decodeSignature(sig []byte) (r, s, v *big.Int) {

	if len(sig) != crypto.SignatureLength {
		panic(fmt.Sprintf("wrong size for signature: got %d, want %d", len(sig), crypto.SignatureLength))
	}
	r = new(big.Int).SetBytes(sig[:32])
	s = new(big.Int).SetBytes(sig[32:64])
	v = new(big.Int).SetBytes([]byte{sig[64] + 27})
	return r, s, v
}

// This is not transaction hash. This is only used for generating signatures
func Hash(tx *Tx) common.Hash {
	return rlpHash([]interface{}{
		tx.To,
		tx.Value,
		tx.Nonce,
	})
}

// HashSigned returns the tx hash
func HashSigned(tx *SignedTx) common.Hash {
	return rlpHash(tx)
}

// rlpHash encodes x and hashes the encoded bytes.
func rlpHash(x interface{}) (h common.Hash) {
	sha := hasherPool.Get().(crypto.KeccakState)
	defer hasherPool.Put(sha)
	sha.Reset()
	rlp.Encode(sha, x)
	sha.Read(h[:])

	return h
}

// SealHash returns the hash of a block prior to it being sealed.
func SealHash(header *Header) (hash common.Hash) {
	hasher := sha3.NewLegacyKeccak256()
	encodeHeader(hasher, header)
	hasher.Sum(hash[:0])

	return hash
}

func encodeHeader(w io.Writer, header *Header) {
	enc := []interface{}{
		header.Number,
		header.StateRoot,
	}

	if err := rlp.Encode(w, enc); err != nil {
		panic("can't encode: " + err.Error())
	}
}

// hasherPool holds LegacyKeccak256 hashers for rlpHash.
var hasherPool = sync.Pool{
	New: func() interface{} { return sha3.NewLegacyKeccak256() },
}
