package types

import (
	"fmt"

	"github.com/ethereum/go-ethereum/rlp"
)

type Header struct {
	ParentHash       [32]byte // The hash of the parent block
	Miner            [20]byte // The address of the miner
	StateRoot        [32]byte // The root of the underlying state
	TransactionsRoot [32]byte // The root of the transactions in the block
	Number           uint64   // The number of the block
	Timestamp        uint64   // The unix timestamp of the block (in seconds)
	ExtraData        []byte   // Extra data of the block
}

func CreateHeader(parentHash [32]byte, miner [20]byte, stateRoot [32]byte, transactionsRoot [32]byte, number uint64, timestamp uint64, extraData []byte) Header {
	return Header{parentHash, miner, stateRoot, transactionsRoot, number, timestamp, extraData}
}

func EncodeHeader(h Header) ([]byte, error) {
	encodedHeader, err := rlp.EncodeToBytes(h)
	if err != nil {
		return nil, err
	}
	return encodedHeader, nil
}

func DecodeHeader(data []byte) (Header, error) {
	var h Header
	if err := rlp.DecodeBytes(data, &h); err != nil {
		return h, err
	}
	return h, nil
}

func (h *Header) Print() {
	fmt.Println("Header:")
	fmt.Println("Parent Hash: ", h.ParentHash)
	fmt.Println("Miner: ", h.Miner)
	fmt.Println("State Root: ", h.StateRoot)
	fmt.Println("Transactions Root: ", h.TransactionsRoot)
	fmt.Println("Number: ", h.Number)
	fmt.Println("Timestamp: ", h.Timestamp)
	fmt.Println("Extra Data: ", h.ExtraData)
}
