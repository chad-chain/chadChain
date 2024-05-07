package types

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	db "github.com/malay44/chadChain/core/storage"
	"github.com/malay44/chadChain/core/utils"
)

type Block struct {
	Header       Header        // The header of the block containing metadata
	Transactions []Transaction // List of transactions in the block
}

var LatestBlock Block

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
		err = db.Update([]byte("latestBlock"), b)(txn)
		if err != nil {
			fmt.Printf("error updating latest block hash: %v", err)
			if err == badger.ErrKeyNotFound {
				err = db.Insert([]byte("stateRootHash"), hash)(txn)
				if err != nil {
					return err
				}
			}
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error adding block to DB: %v", err)
	}

	return nil
}

// create block
func CreateBlock(header Header, transactions []Transaction) Block {
	block := new(Block)
	block.Header = header
	block.Transactions = transactions
	return *block
}

func EncodeBlock(block Block) ([]byte, error) {
	// Encode the header
	encodedHeader, err := EncodeHeader(block.Header)
	if err != nil {
		return nil, err
	}

	// Encode the transactions
	var encodedTransactions [][]byte
	for _, tx := range block.Transactions {
		encodedTx, err := EncodeTransaction(tx)
		if err != nil {
			return nil, err
		}
		encodedTransactions = append(encodedTransactions, encodedTx)
	}

	// Combine encoded header and transactions
	encodedBlock := [][]byte{encodedHeader}
	encodedBlock = append(encodedBlock, encodedTransactions...)

	// RLP encode the entire block
	encodedData, err := rlp.EncodeToBytes(encodedBlock)
	if err != nil {
		return nil, err
	}

	return encodedData, nil
}

func DecodeBlock(data []byte) (Block, error) {
	var decodedBlock [][]byte
	if err := rlp.DecodeBytes(data, &decodedBlock); err != nil {
		return Block{}, err
	}

	// Decode header
	decodedHeader, err := DecodeHeader(decodedBlock[0])
	if err != nil {
		return Block{}, err
	}

	// Decode transactions
	var decodedTransactions []Transaction
	for _, txData := range decodedBlock[1:] {
		decodedTx, err := DecodeTransaction(txData)
		if err != nil {
			return Block{}, err
		}
		decodedTransactions = append(decodedTransactions, decodedTx)
	}

	return Block{
		Header:       decodedHeader,
		Transactions: decodedTransactions,
	}, nil
}

func GetParentBlock() Block {
	// get parent block
	return Block{}
}

// Get Parent Block Height
func GetParentBlockHeight() uint64 {
	// get parent block height
	return 0
}

// Get ParentBlock StateRoot
func GetParentBlockStateRoot() [32]byte {
	// get parent block state root
	return [32]byte{}
}

// Get ParentBlock TransactionsRoot
func GetParentBlockTransactionsRoot() [32]byte {
	// get parent block transactions root
	return [32]byte{}
}

// Get ParentHash
func GetParentHash() [32]byte {
	// get parent hash
	return [32]byte{}
}
