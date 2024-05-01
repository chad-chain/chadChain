package main

import (
	"log"

	// n "github.com/malay44/chadChain/core/network"
	db "github.com/malay44/chadChain/core/storage"
	t "github.com/malay44/chadChain/core/types"
)

func main() {
	// n.Http()
	db.InitBadger()
	defer db.BadgerDB.Close()
	// n.Rpc()
	bbc := t.Block{}
	bbch := t.Header{}
	bbct := t.Transaction{}
	bbch.CreateHeader([32]byte{}, [20]byte{}, [32]byte{}, [32]byte{}, 0, 0, 0, 0, []byte{}, 0)
	bbct.CreateTransaction([20]byte{}, 0, 0, nil, nil, nil)

	bbc.CreateBlock(bbch, []t.Transaction{bbct})
	// testDBFunc(bbc)

	// Test Account creation
	acc := t.Account{}
	accAddr := [32]byte{}
	copy(accAddr[:], "0x123456789012345678901234567890")
	accNonce := uint64(0)
	accBalance := uint64(0)
	acc = acc.CreateAccount(accAddr, accNonce, accBalance)
	acc.AddAccount()

	// Test Account retrieval
	// accAddr := "0x1234567890123456789012345678901234567890"
	// acc, err := t.GetAccount(accAddr)
	// if err != nil {
	// 	log.Default().Println(err.Error())
	// }
	// log.Default().Println(acc)

	log.Default().Println("Hello, world!")
}

// func testDBFunc(block t.Block) {
// 	err := db.BadgerDB.Update(db.Insert([]byte("block"), block))
// 	if err != nil {
// 		println("Update Error", err.Error())
// 		log.Fatal(err)
// 	}

// 	retrievedBlock := t.Block{}
// 	err = db.BadgerDB.View(db.Get([]byte("block"), &retrievedBlock))
// 	if err != nil {
// 		println("View Error", err.Error())
// 		log.Fatal(err)
// 	}
// 	log.Default().Println("Block retrieved")
// 	log.Default().Println(retrievedBlock)
// }
