package types

import (
	"fmt"
	"math/big"
)

type Transaction struct {
	To    [20]byte // The address of the receiver
	Value uint64   // The value of the transaction
	Nonce uint64   // The nonce of the sender of the transaction
	V     *big.Int // Signature value v of the transaction
	R     *big.Int // Signature value r of the transaction
	S     *big.Int // Signature value s of the transaction
}

type UnSignedTx struct {
	To    [20]byte // The address of the receiver
	Value uint64   // The value of the transaction
	Nonce uint64   // The nonce of the sender of the transaction
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
	fmt.Println("Transaction:")
	fmt.Println("To:", t.To)
	fmt.Println("Value:", t.Value)
	fmt.Println("Nonce:", t.Nonce)
	fmt.Println("V:", t.V)
	fmt.Println("R:", t.R)
	fmt.Println("S:", t.S)
}
