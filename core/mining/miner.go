package mining

import (
	"encoding/binary"
	"errors"
	"log"
	"time"

	t "github.com/malay44/chadChain/core/types"
)

// MineBlock mines the block by finding a hash that meets the difficulty target
func MineBlock(header t.Header, transactionpool *t.Transactionpool) (t.Block, error) {
	log.Default().Println("Mining Started")
	b := new(t.Block)
	b.Header = header
	b.Transactions = transactionpool.Get_all_transactions_and_clear()

	timeout := 2 * time.Second
	start := time.Now()
	time.Sleep(4 * time.Second)
	for {
		elapsed := time.Since(start) // Check if the timeout has been reached
		if elapsed >= timeout {
			tmpblk := new(t.Block)
			log.Default().Println("Timeout For Current Block Mining")
			transactionpool.AddTransactions(b.Transactions)
			return *tmpblk, errors.New("Timeout")
		}

		b.Header.Nonce++
		generated_hash := t.Keccak256(SerializeBlock(*b)) // Creating Hash

		if meetsDifficultyTarget([32]byte(generated_hash), header.Difficulty) {
			log.Default().Println("Block mined")
			return *b, nil
		}
	}
}

// SerializeHeader serializes the header and the transactions of the block
func SerializeBlock(block t.Block) []byte {
	var serializedBlock []byte

	serializedBlock = append(serializedBlock, SerializeHeader(block.Header)...)
	for _, transaction := range block.Transactions {
		serializedBlock = append(serializedBlock, SerializeTransaction(transaction)...)
	}
	return serializedBlock
}

// function to serialize the transaction
func SerializeTransaction(transaction t.Transaction) []byte {
	var serializedTransaction []byte

	serializedTransaction = append(serializedTransaction, transaction.To[:]...)
	serializedTransaction = append(serializedTransaction, uint64ToBytes(transaction.Value)...)
	serializedTransaction = append(serializedTransaction, uint64ToBytes(transaction.Nonce)...)
	serializedTransaction = append(serializedTransaction, transaction.V.Bytes()...) // Convert transaction.V to []byte
	serializedTransaction = append(serializedTransaction, transaction.R.Bytes()...)
	serializedTransaction = append(serializedTransaction, transaction.S.Bytes()...)

	return serializedTransaction
}

func SerializeHeader(header t.Header) []byte {
	var serializedHeader []byte

	serializedHeader = append(serializedHeader, header.ParentHash[:]...)
	serializedHeader = append(serializedHeader, header.Miner[:]...)
	serializedHeader = append(serializedHeader, header.StateRoot[:]...)
	serializedHeader = append(serializedHeader, header.TransactionsRoot[:]...)
	serializedHeader = append(serializedHeader, uint64ToBytes(header.Difficulty)...)
	serializedHeader = append(serializedHeader, uint64ToBytes(header.TotalDifficulty)...)
	serializedHeader = append(serializedHeader, uint64ToBytes(header.Number)...)
	serializedHeader = append(serializedHeader, uint64ToBytes(header.Timestamp)...)
	serializedHeader = append(serializedHeader, header.ExtraData...)
	serializedHeader = append(serializedHeader, uint64ToBytes(header.Nonce)...)
	return serializedHeader
}

func meetsDifficultyTarget(hash [32]byte, difficulty uint64) bool {

	maxTarget := ^uint64(0) / difficulty        // Calculate the maximum target value
	hashInt := binary.BigEndian.Uint64(hash[:]) // Convert the hash to a big-endian integer
	return hashInt <= maxTarget                 // Check if the hash meets the difficulty target
}

// uint64ToBytes converts a uint64 value to a byte slice
func uint64ToBytes(value uint64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, value)
	return bytes
}
