package types

import (
	"github.com/ethereum/go-ethereum/rlp"
)

type Block struct {
	Header       Header        // The header of the block containing metadata
	Transactions []Transaction // List of transactions in the block
}

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
func (b *Block) AddBlockToChain(blck Block) {
	// sava block into chain
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
