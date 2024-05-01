package main

import (
	"context"
	// "log"

	n "github.com/malay44/chadChain/core/network"
	// db "github.com/malay44/chadChain/core/storage"
	// t "github.com/malay44/chadChain/core/types"
)

func main() {
	// n.Http()
	ctx := context.Background()
	n.Run(ctx, []string{
		"/ip4/127.0.0.1/tcp/61062/p2p/12D3KooWBFe8MmXXtghkiXp94kxt87Egta2v4Uqov8yvjKBAVKEe",
	})
	// db.InitBadger()
	// defer db.BadgerDB.Close()

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
	// if err != nil {
	// 	println("View Error", err.Error())
	// 	log.Fatal(err)
	// }
	// log.Default().Println("Block retrieved")
	// log.Default().Println(retrievedBlock)

	// log.Default().Println("Hello, world!")

}
