package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	n "github.com/malay44/chadChain/core/network"
	db "github.com/malay44/chadChain/core/storage"
	t "github.com/malay44/chadChain/core/types"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	// n.Http()
	n.PeerAddrs = []string{
		// "/ip4/127.0.0.1/tcp/64561/p2p/12D3KooWPot5PSrTg6KQA5VChBzTs6GSoNgfnPzXtkWdKQ8wFAxQ",
	}

	n.CtxVar = context.Background()
	n.Run()
	db.InitBadger()
	defer db.BadgerDB.Close()
	initGlobalVar()

	// bbc := t.Block{}
	// bbch := t.Header{}
	// bbct := t.Transaction{}
	// bbca := t.Account{}
	// bbch.CreateHeader([32]byte{}, [20]byte{}, [32]byte{}, [32]byte{}, 0, 0, 0, 0, []byte{}, 0)
	// bbct.CreateTransaction([20]byte{}, 0, 0, nil, nil, nil)
	// bbca.CreateAccount([20]byte{}, 0, 0)

	// bbc.CreateBlock(bbch, []t.Transaction{bbct})
	// // log.Default().Println("Block created")
	// // log.Default().Println(bbc)

	// err := db.BadgerDB.Update(db.Insert([]byte("block"), bbc))
	// if err != nil {
	// 	println("Update Error", err.Error())
	// 	log.Fatal(err)
	// }

	// retrievedBlock := t.Block{}
	// err = db.BadgerDB.View(db.Get([]byte("block"), &retrievedBlock))

	log.Default().Println("Hello, world!")
}

func initGlobalVar() {
	err := db.BadgerDB.View(db.Get([]byte("stateRootHash"), &t.StateRootHash))
	if err != nil {
		log.Default().Printf("StateRootHash not found\n")
		log.Default().Println(err.Error())
		return
	}
	log.Default().Printf("StateRootHash: %x\n", t.StateRootHash)
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
