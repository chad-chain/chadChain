package types

import (
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Transaction struct {
	To    common.Address // The address of the receiver
	Value uint64         // The value of the transaction
	Nonce uint64         // The nonce of the sender of the transaction
	V     *big.Int       // Signature value v of the transaction
	R     *big.Int       // Signature value r of the transaction
	S     *big.Int       // Signature value s of the transaction
}

type UnSignedTx struct {
	To    common.Address // The address of the receiver
	Value uint64         // The value of the transaction
	Nonce uint64         // The nonce of the sender of the transaction
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

func (t *Transaction) Print() {
	log.Default().Println("Transaction:")
	log.Default().Println("To:", t.To)
	log.Default().Println("Value:", t.Value)
	log.Default().Println("Nonce:", t.Nonce)
	log.Default().Println("V:", t.V)
	log.Default().Println("R:", t.R)
	log.Default().Println("S:", t.S)
}
