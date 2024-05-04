package mining

import (
	"log"
	"time"

	t "github.com/malay44/chadChain/core/types"
)

func BuildBlock(chn chan t.Block, transactionpool *t.Transactionpool, Miner [20]byte) (t.Block, error) {
	log.Default().Println("Building Block Function")
	timestamp := time.Now().Unix()
	transactionroot, err := t.EncodeTransactions(transactionpool.Get_all_transactions_and_clear())

	if err != nil {
		log.Default().Println("Error in encoding transaction")
		emptyBlock := new(t.Block)
		emptyBlock.Header.ParentHash = t.GetParentHash()
		emptyBlock.Header.StateRoot = [32]byte(t.StateRootHash)
		emptyBlock.Header.TransactionsRoot = [32]byte(t.GetParentBlockTransactionsRoot())
		emptyBlock.Header.Number = t.GetParentBlockHeight() + 1
		emptyBlock.Header.Timestamp = uint64(timestamp)

		return *emptyBlock, err
	}

	transactionrootBytes := [32]byte{}
	copy(transactionrootBytes[:], transactionroot)

	header := new(t.Header)
	header.ParentHash = t.GetParentHash()
	header.Miner = Miner
	// header.StateRoot = [32]byte(t.StateRootHash)
	header.TransactionsRoot = transactionrootBytes
	header.Number = t.GetParentBlockHeight() + 1
	header.Timestamp = uint64(timestamp)

	b := t.CreateBlock(*header, transactionpool.Get_all_transactions_and_clear())
	log.Default().Println("Block Created", b)
	minedBlock, err := MineBlock(b)
	if err != nil {
		log.Default().Println("Error in encoding transaction")
		emptyBlock := new(t.Block)
		emptyBlock.Header.ParentHash = t.GetParentHash()
		emptyBlock.Header.StateRoot = [32]byte(t.StateRootHash)
		emptyBlock.Header.TransactionsRoot = [32]byte(t.GetParentBlockTransactionsRoot())
		emptyBlock.Header.Number = t.GetParentBlockHeight() + 1
		emptyBlock.Header.Timestamp = uint64(timestamp)

		return *emptyBlock, err
	}
	chn <- minedBlock
	return b, nil
}

func MineBlock(b t.Block) (t.Block, error) {
	log.Default().Println("Mining Started")
	encodded_block, err := t.EncodeBlock(b)
	if err != nil {
		log.Default().Println("Error in encoding block")
		return b, err
	}
	log.Default().Println("Block Hash: ", encodded_block)
	generated_hash := t.Keccak256(encodded_block)
	log.Default().Println("Generated Hash: ", generated_hash)
	log.Default().Println("Block mined", b)
	return b, nil
}
