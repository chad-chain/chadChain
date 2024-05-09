package mining

import (
	"log"
	"time"

	t "github.com/malay44/chadChain/core/types"
	rlp "github.com/malay44/chadChain/core/utils"
)

func BuildBlock(chn chan t.Block, transactionPool *t.TransactionPool, Miner string) (t.Block, error) {
	log.Default().Println("Building Block Function")
	timestamp := time.Now().Unix()
	transactionroot, err := rlp.EncodeData(transactionPool.Get_all_transactions_and_clear(), false)

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

	transactionRootBytes := [32]byte{}
	copy(transactionRootBytes[:], transactionroot)
	miner, err := rlp.EncodeData(Miner, false)
	minerBytes := [20]byte{}
	copy(minerBytes[:], miner)
	header := new(t.Header)
	header.ParentHash = t.GetParentHash()
	header.Miner = minerBytes
	// header.StateRoot = [32]byte(t.StateRootHash)
	header.TransactionsRoot = transactionRootBytes
	header.Number = t.GetParentBlockHeight() + 1
	header.Timestamp = uint64(timestamp)

	b := t.CreateBlock(*header, transactionPool.Get_all_transactions_and_clear())
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
		emptyBlock.Header.Miner = minerBytes
		return *emptyBlock, err
	}
	chn <- minedBlock
	return b, nil
}

func MineBlock(b t.Block) (t.Block, error) {
	log.Default().Println("Mining Started")
	encodded_block, err := rlp.EncodeData(b, false)
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
