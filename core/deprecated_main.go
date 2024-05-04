package deprecated_main

import (
	"fmt"
	"log"
	"math/big"

	t "github.com/malay44/chadChain/core/types"
)

func _main() {
	// n.Http()
	// n.Rpc()

	// The block timing of the network is 2s
	// A thred runs at every 2 Sec and checks whether the block is mined or not based on the block number

	// If primary node is not able to mine the block, then the secondary node will mine the block

	// bbc := t.Block{}
	// bbch := t.Header{}
	// bbct := t.Transaction{}
	// bbca := t.Account{}
	// bbch.CreateHeader([32]byte{}, [20]byte{}, [32]byte{}, [32]byte{}, 0, 0, 0, 0, []byte{}, 0)
	// bbct.CreateTransaction([20]byte{}, 0, 0, nil, nil, nil)
	// bbca.CreateAccount([20]byte{}, 0, 0)

	// fmt.Println(bbch)
	// fmt.Println(bbct)
	// fmt.Println(bbca)

	// bbc.CreateBlock(bbch, []t.Transaction{bbct})
	// fmt.Println(bbc)

	fmt.Println("Hello, world!")
	TestMineBlock()

}

func TestMineBlock() {
	// Create a dummy header for testing
	// header := t.Header{
	// 	ParentHash:       [32]byte{},
	// 	Miner:            [20]byte{},
	// 	StateRoot:        [32]byte{},
	// 	TransactionsRoot: [32]byte{},
	// 	Difficulty:       8,
	// 	TotalDifficulty:  8,
	// 	Number:           875463,
	// 	Timestamp:        uint64(time.Now().Unix()), // Set current Unix timestamp
	// 	ExtraData:        []byte{},
	// 	Nonce:            76,
	// }

	// Create a dummy transaction pool for testing

	dummyTransaction := t.Transaction{
		To:    [20]byte{},
		Value: 2354,
		Nonce: 67,
		V:     new(big.Int).SetInt64(6743),
		R:     new(big.Int).SetInt64(0234557),
		S:     new(big.Int).SetInt64(79652),
	}
	dummyTransaction1 := t.Transaction{
		To:    [20]byte{},
		Value: 2354,
		Nonce: 67,
		V:     new(big.Int).SetInt64(6743),
		R:     new(big.Int).SetInt64(0234557),
		S:     new(big.Int).SetInt64(79652),
	}

	transactionPool := &t.TransactionPool{
		Transactions: []t.Transaction{dummyTransaction, dummyTransaction1},
	}
	log.Default().Println("transactionPool size : ", len(transactionPool.Transactions))

	// Call the MineBlock function
	// minedBlock, err := m.MineBlock(header, transactionPool)

	// if err != nil {
	// 	log.Default().Println("Mining failed: ", err)
	// }
	// log.Default().Println("transactionPool size : ", len(transactionPool.Transactions))

	// dummyTransaction2 := t.Transaction{
	// 	To:    [20]byte{},
	// 	Value: 2354,
	// 	Nonce: 67,
	// 	V:     new(big.Int).SetInt64(67633),
	// 	R:     new(big.Int).SetInt64(023557),
	// 	S:     new(big.Int).SetInt64(7563652),
	// }

	// transactionPool.Transactions = append(transactionPool.Transactions, dummyTransaction2)
	// transactionPool.Transactions = append(transactionPool.Transactions, dummyTransaction1)
	// transactionPool.Transactions = append(transactionPool.Transactions, dummyTransaction)

	// transactionPool.RemoveCommonTransactions(minedBlock.Transactions)

	// Add more assertions as needed
}
