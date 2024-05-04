package main

import (
	"log"
	"time"

	m "github.com/malay44/chadChain/core/mining"
	n "github.com/malay44/chadChain/core/network"
	t "github.com/malay44/chadChain/core/types"
)

func main() {
	// n.Http()
	// n.PeerAddrs = []string{
	// 	// "/ip4/127.0.0.1/tcp/64561/p2p/12D3KooWPot5PSrTg6KQA5VChBzTs6GSoNgfnPzXtkWdKQ8wFAxQ",
	// }
	n.PeerAddrs = []string{
		"12D3KooWPot5PSrTg6K",
		"12D3KooWPot5PSrTg6K",
		"12D3KooWPot5PSrTg6K"}

	// n.CtxVar = context.Background()
	// n.Run()
	// db.InitBadger()
	// defer db.BadgerDB.Close()
	// initGlobalVar()

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

	// ch := make(chan t.Block)
	chn := make(chan t.Block)
	timerCh := make(chan string)

	go Timer(timerCh, n.PeerAddrs)
	log.Default().Println("Both Chanells Created")

	transactionpool := t.Transactionpool{}
	for {
		select {
		case miner := <-timerCh:
			log.Default().Println("Miner selected", miner)
			// var minerinByte [20]byte
			// copy(minerinByte[:], []byte(miner))

			go m.BuildBlock(chn, &transactionpool, [20]byte{})
		case blk := <-chn:
			log.Default().Println("Mined Block: ", blk)
			log.Default().Println(blk)
		}
	}
}

func Timer(timerCh chan string, miners []string) {
	// find miner inde
	index := 1
	numberOfMiners := len(miners)
	interval := 2 * numberOfMiners
	time.Sleep(time.Duration(index) * time.Second)
	timerCh <- miners[index]

	for {
		log.Default().Println("Timer loop 2 sec")
		time.Sleep(time.Duration(interval) * time.Second)
		timerCh <- miners[index]
	}
}

// func test(ch chan t.Block, chn chan t.Block) {
// 	for {
// 		select {
// 		case blk := <-ch:
// 			log.Default().Println("Block received")
// 			log.Default().Println(blk)
// 			go m.MineBlock(chn, blk)
// 		case blkm := <-chn:
// 			log.Default().Println("Mined block received")
// 			log.Default().Println(blkm)
// 			chn <- blkm
// 			return
// 		}
// 	}
// }

// func initGlobalVar() {
// 	err := db.BadgerDB.View(db.Get([]byte("stateRootHash"), &t.StateRootHash))
// 	if err != nil {
// 		log.Default().Printf("StateRootHash not found\n")
// 		log.Default().Println(err.Error())
// 		return
// 	}
// 	log.Default().Printf("StateRootHash: %x\n", t.StateRootHash)
// }

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
