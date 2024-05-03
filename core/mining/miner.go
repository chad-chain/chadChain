package mining

import (
	"log"

	t "github.com/malay44/chadChain/core/types"
)

func MineBlock(header t.Header, transactionpool *t.Transactionpool) (t.Block, error) {
	log.Default().Println("Mining Started")

	b := t.CreateBlock(header, transactionpool.Get_all_transactions_and_clear())
	encodded_block, err := t.EncodeBlock(b)
	if err != nil {
		log.Default().Println("Error in encoding block")
		return b, err
	}
	log.Default().Println("Block Hash: ", encodded_block)
	generated_hash := t.Keccak256(encodded_block)
	log.Default().Println("Generated Hash: ", generated_hash)
	log.Default().Println("Block mined")

	return b, nil
}
