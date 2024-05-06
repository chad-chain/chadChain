package main

import (
	"context"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"
	m "github.com/malay44/chadChain/core/mining"
	n "github.com/malay44/chadChain/core/network"
	db "github.com/malay44/chadChain/core/storage"
	t "github.com/malay44/chadChain/core/types"
	"github.com/malay44/chadChain/core/validator"
)

func main() {
	n.CtxVar = context.Background()
	n.PeerAddrs = []string{
		"/ip4/192.168.1.8/tcp/3000/p2p/12D3KooWEDdhybEFMXhN1kzH5iaCZvaBfAGHXqjo83AQ1dkE3Yt5",
	}
	go func() {
		n.Rpc()
	}()
	go func() {
		n.Run()
	}()
	// db.InitBadger()
	// defer db.BadgerDB.Close()
	// initialize.GlobalDBVar()

	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	// n.PeerAddrs = []string{
	// 	"12D3KooWPot5PSrTg6K",
	// 	"12D3KooWPot5PSrTg6K",
	// 	"12D3KooWPot5PSrTg6K"}

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

	miningInit()
	select {}
}

func TestTransactionSig() {
	transaction := t.UnSignedTx{
		To:    common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8"),
		Value: 0,
		Nonce: 0,
	}

	privateKey, _, accAddr, _ := t.GenerateNewPrivateKey()

	log.Default().Println("Private Key:", privateKey)

	log.Default().Println("Account Address:", accAddr)

	signedTx, err := t.SignTransaction(&transaction, privateKey)

	if err != nil {
		log.Default().Println("Failed to sign transaction:", err)
	}

	log.Default().Println("Signed Transaction:", signedTx)

	log.Default().Println("Transaction:", transaction)

	sender, err := t.VerifyTxSign(&signedTx)

	if err != nil {
		log.Default().Println("Failed to recover sender:", err)
	}

	log.Default().Println("Sender:", sender)

	if sender == accAddr {
		log.Default().Println("Transaction is valid")
	} else {
		log.Default().Println("Transaction is invalid")
	}
}

func TestBlockSig() {
	privateKey, _, accAddr, _ := t.GenerateNewPrivateKey()

	header := t.Header{
		ParentHash:       [32]byte{},
		StateRoot:        [32]byte{},
		TransactionsRoot: [32]byte{},
		Timestamp:        uint64(time.Now().Unix()),
		Number:           0,
		Miner:            accAddr,
		ExtraData:        []byte{},
	}

	sig, err := t.SignHeader(&header, privateKey)
	header.ExtraData = sig

	block := t.Block{
		Header:       header,
		Transactions: []t.Transaction{},
	}

	if err != nil {
		log.Default().Println("Failed to sign header:", err)
	}

	log.Default().Println("Signature:", sig)

	if validator.ValidateBlock(&block) {
		log.Default().Println("Block is valid")
	} else {
		log.Default().Println("Block is invalid")
	}

}

func miningInit() {

	// ch := make(chan t.Block)
	chn := make(chan t.Block)
	timerCh := make(chan string)

	go Timer(timerCh, n.PeerAddrs)
	log.Default().Println("Both Chanells Created")

	transactionPool := t.TransactionPool{}
	for {
		select {
		case miner := <-timerCh:
			log.Default().Println("Miner selected", miner)
			// var minerinByte [20]byte
			// copy(minerinByte[:], []byte(miner))

			go m.BuildBlock(chn, &transactionPool, [20]byte{})
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

func testDBFunc(block t.Block) {
	err := db.BadgerDB.Update(db.Insert([]byte("block"), block))
	if err != nil {
		println("Update Error", err.Error())
		log.Fatal(err)
	}

	retrievedBlock := t.Block{}
	err = db.BadgerDB.View(db.Get([]byte("block"), &retrievedBlock))
	if err != nil {
		println("View Error", err.Error())
		log.Fatal(err)
	}
	log.Default().Println("Block retrieved")
	log.Default().Println(retrievedBlock)
}
