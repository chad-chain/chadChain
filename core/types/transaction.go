package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/rlp"
)

type Transaction struct {
	To    [20]byte // The address of the receiver
	Value uint64   // The value of the transaction
	Nonce uint64   // The nonce of the sender of the transaction
	V     *big.Int // Signature value v of the transaction
	R     *big.Int // Signature value r of the transaction
	S     *big.Int // Signature value s of the transaction
}

func CreateTransaction(to [20]byte, value uint64, nonce uint64, v *big.Int, r *big.Int, s *big.Int) Transaction {
	return Transaction{to, value, nonce, v, r, s}
}

// send transaction over network
func (t *Transaction) SendTransaction(tr Transaction) {
	// propogate in the network
}

// Get transaction from network
func (t *Transaction) AddTransaction(tr Transaction) {
	// get transaction from network
	// add to transaction pool
}

func EncodeTransaction(tx Transaction) ([]byte, error) {
	encodedTx, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}
	return encodedTx, nil
}

func DecodeTransaction(data []byte) (Transaction, error) {
	var tx Transaction
	if err := rlp.DecodeBytes(data, &tx); err != nil {
		return tx, err
	}
	return tx, nil
}

// Validate transaction
func (t *Transaction) ValidateTransaction(tr Transaction) bool {
	// validate transaction
	return true
}
