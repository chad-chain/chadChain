package types

type Header struct {
	ParentHash       [32]byte // The hash of the parent block
	Miner            [20]byte // The address of the miner
	StateRoot        [32]byte // The root of the underlying state
	TransactionsRoot [32]byte // The root of the transactions in the block
	Difficulty       uint64   // The difficulty of the block
	TotalDifficulty  uint64   // The total difficulty of the chain till this block
	Number           uint64   // The number of the block
	Timestamp        uint64   // The unix timestamp of the block (in seconds)
	ExtraData        []byte   // Extra data of the block
	Nonce            uint64   // The nonce of the block
}

// create header
func (h *Header) CreateHeader(parentHash [32]byte, miner [20]byte, stateRoot [32]byte, transactionsRoot [32]byte, difficulty uint64, totalDifficulty uint64, number uint64, timestamp uint64, extraData []byte, nonce uint64) Header {
	return Header{parentHash, miner, stateRoot, transactionsRoot, difficulty, totalDifficulty, number, timestamp, extraData, nonce}
}
