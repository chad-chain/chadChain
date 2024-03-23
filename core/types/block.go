package types

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
func (b *Block) CreateBlock(header Header, transactions []Transaction) Block {
	return Block{header, transactions}
}

func (b *Block) CalculateHash() [32]byte {
	// hash generation Finction
	return [32]byte{}
}

// mine block
func (b *Block) MineBlock(blck *Block) {

	for b.CalculateHash()[0] != byte(0xff>>uint8(b.Header.Difficulty)) {
		b.Header.Nonce++
	}
}
