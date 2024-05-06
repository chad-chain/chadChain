package types

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
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

func GenerateNewPrivateKey() (*ecdsa.PrivateKey, string, common.Address, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, "", common.Address{}, err
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	privateKeyHex := hex.EncodeToString(crypto.FromECDSA(privateKey))
	return privateKey, privateKeyHex, address, nil
}

func SignTransaction(tx *UnSignedTx, privateKey *ecdsa.PrivateKey) (Transaction, error) {
	h := Hash(tx)
	sig, err := crypto.Sign(h[:], privateKey)
	if err != nil {
		return Transaction{}, err
	}
	R, S, V := decodeSignature(sig)
	// Create a signed transaction by V, R, S values
	signedTx := Transaction{
		To:    tx.To,
		Value: tx.Value,
		Nonce: tx.Nonce,
		V:     V,
		R:     R,
		S:     S,
	}
	return signedTx, nil
}

func VerifyTxSign(t *Transaction) (common.Address, error) {
	UnSignedTx := UnSignedTx{
		To:    t.To,
		Value: t.Value,
		Nonce: t.Nonce,
	}
	sender, err := recoverPlain(Hash(&UnSignedTx), t.R, t.S, t.V, true)
	if err != nil {
		return common.Address{}, err
	}
	return sender, nil
}

// BytesToHexString converts a byte slice to a hex string
func BytesToHexString(b []byte) string {
	return hex.EncodeToString(b)
}

// sealHash returns the hash of a block prior to it being sealed.
func sealHash(header *Header) (hash common.Hash) {
	hasher := sha3.NewLegacyKeccak256()
	encodeHeader(hasher, header)
	hasher.Sum(hash[:0])

	return hash
}

func SignHeader(header *Header, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	sig, err := crypto.Sign(sealHash(header).Bytes(), privateKey)
	if err != nil {
		return nil, err
	}
	return sig, nil
}

func VerifyHeader(block *Block) bool {
	pubKey, err := crypto.Ecrecover(sealHash(&block.Header).Bytes(), block.Header.ExtraData)
	if err != nil {
		log.Default().Println("Failed to recover public key:", err)
		return false
	}

	var signer common.Address

	copy(signer[:], crypto.Keccak256(pubKey[1:])[12:])
	return signer == block.Header.Miner
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

// recoverPlain recovers the address which has signed the given data using the v, r, and s values
func recoverPlain(sigHash common.Hash, R, S, Vb *big.Int, homestead bool) (common.Address, error) {
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
	pub, err := crypto.Ecrecover(sigHash[:], sig)
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

// This is not transaction hash. This is only used for generating signatures
func Hash(tx *UnSignedTx) common.Hash {
	return rlpHash([]interface{}{
		tx.To,
		tx.Value,
		tx.Nonce,
	})
}

// HashSigned returns the tx hash
func HashSigned(tx *Transaction) common.Hash {
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

func encodeHeader(w io.Writer, header *Header) {
	// Remove the ExtraData field from the header before encoding
	enc := []interface{}{
		header.ParentHash,
		header.Miner,
		header.StateRoot,
		header.TransactionsRoot,
		header.Number,
		header.Timestamp,
	}

	if err := rlp.Encode(w, enc); err != nil {
		panic("can't encode: " + err.Error())
	}
}

// hasherPool holds LegacyKeccak256 hashers for rlpHash.
var hasherPool = sync.Pool{
	New: func() interface{} { return sha3.NewLegacyKeccak256() },
}
