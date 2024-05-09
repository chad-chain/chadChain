package types

import (
	"fmt"

	db "github.com/chad-chain/chadChain/core/storage"
	"github.com/chad-chain/chadChain/core/utils"
	"github.com/dgraph-io/badger/v4"
	"github.com/ethereum/go-ethereum/crypto"
)

type Block struct {
	Header       Header        // The header of the block containing metadata
	Transactions []Transaction // List of transactions in the block
}

var LatestBlock Block
var LatestBlockHash [32]byte

// add single transaction to block
func (b *Block) AddTransactionToBlock(t Transaction) {
	b.Transactions = append(b.Transactions, t)
}

// add list of transaction to block
func (b *Block) AddTransactions(t []Transaction) {
	b.Transactions = append(b.Transactions, t...)
}

// remove single transaction from block
func (b *Block) RemoveTransactionFromBlock(t Transaction) {
	for i, transaction := range b.Transactions {
		if transaction == t {
			b.Transactions = append(b.Transactions[:i], b.Transactions[i+1:]...)
		}
	}
}

// Getting Validated Block from network
func (b *Block) AddBlockToChain() error {
	marshalledHeader, err := utils.EncodeData(b.Header, false)
	hash := crypto.Keccak256(marshalledHeader)

	if err != nil {
		return fmt.Errorf("error encoding block header: %v", err)
	}
	key := []byte("block" + string(hash))
	err = db.BadgerDB.Update(func(txn *badger.Txn) error {
		err := db.Insert(key, *b)(txn)
		if err != nil {
			return fmt.Errorf("error inserting block into db: %v", err)
		}
		// Update the latest block in the db
		err = db.Update([]byte("latestBlock"), *b)(txn)
		if err != nil {
			fmt.Printf("error updating latest block hash: %v", err)
			if err == badger.ErrKeyNotFound {
				err = db.Insert([]byte("latestBlock"), *b)(txn)
				if err != nil {
					return err
				}
			}
			return err
		}
		LatestBlock = *b
		copy(LatestBlockHash[:], hash)
		return nil
	})
	if err != nil {
		return fmt.Errorf("error adding block to DB: %v", err)
	}
	return nil
}

// create block
func CreateBlock(header *Header, transactions *[]Transaction) *Block {
	block := new(Block)
	block.Header = *header
	block.Transactions = *transactions
	return block
}
