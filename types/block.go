package types

type Block struct {
	prevHash     [32]byte
	timestamp    uint64
	transactions []Transaction
	hash         []byte // TODO: Make a universal Hash Type with methods
}

func (b *Block) String() string {
	res := "Block{\n"
	res += "  prevHash: " + string(b.prevHash[:]) + "\n"
	res += "  timestamp: " + string(rune(b.timestamp)) + "\n"
	res += "  transactions: [\n"
	for _, tx := range b.transactions {
		res += "    " + tx.String() + "\n"
	}
	res += "  ]\n"
	res += "  hash: " + string(b.hash) + "\n"
	res += "}"
	return res
}

func (b *Block) AddTransaction(t Transaction) {
	b.transactions = append(b.transactions, t)
	// TODO: update hash
}
