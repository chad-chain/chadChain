package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/malay44/chadChain/core/crypto"
	m "github.com/malay44/chadChain/core/mining"
	n "github.com/malay44/chadChain/core/network"
	db "github.com/malay44/chadChain/core/storage"
	t "github.com/malay44/chadChain/core/types"
	"github.com/malay44/chadChain/core/validator"
)

func main() {
	n.CtxVar = context.Background()
	// n.PeerAddrs = []string{
	// 	// "/ip4/127.0.0.1/tcp/63795/p2p/12D3KooWBMNwiqwM1DhRXDaTXU2CYCdpKMDY5tNxfviq7ogFmFhW",
	// }
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
	// 	"/ip4/192.168.1.4/tcp/3000/p2p/12D3KooWEDdhybEFMXhN1kzH5iaCZvaBfAGHXqjo83AQ1dkDxBBB",
	// 	"/ip4/192.168.1.8/tcp/3000/p2p/12D3KooWEDdhybEFMXhN1kzH5iaCZvaBfAGHXqjo83AQ1dkE3Yt5",
	// }
	// n.PeerAddrs = []string{
	// 	"12D3KooWPot5PSrTg6K",
	// 	"12D3KooWPot5PSrTg6",
	// 	"12D3KooWPot5PSrT6K"}

	// n.CtxVar = context.Background()
	// n.Run()

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

	select {}
	log.Default().Println("Hello, world!")
	expectedMiners := make(chan string)
	MiningInit(expectedMiners)
}

func TestTransactionSig() {
	transaction := t.UnSignedTx{
		To:    common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8"),
		Value: 0,
		Nonce: 0,
	}

	privateKey, _, accAddr, _ := crypto.GenerateNewPrivateKey()

	log.Default().Println("Private Key:", privateKey)

	log.Default().Println("Account Address:", accAddr)

	signedTx, err := crypto.SignTransaction(&transaction)

	if err != nil {
		log.Default().Println("Failed to sign transaction:", err)
	}

	log.Default().Println("Signed Transaction:", signedTx)

	log.Default().Println("Transaction:", transaction)

	sender, err := crypto.VerifyTxSign(&signedTx)

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
	_, _, accAddr, _ := crypto.GenerateNewPrivateKey()

	header := t.Header{
		ParentHash:       [32]byte{},
		StateRoot:        [32]byte{},
		TransactionsRoot: [32]byte{},
		Timestamp:        uint64(time.Now().Unix()),
		Number:           0,
		Miner:            accAddr,
		ExtraData:        []byte{},
	}

	sig, err := crypto.SignHeader(&header)
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

func MiningInit(expectedMiner chan string) { // add transactionpool as argument

	// ch := make(chan t.Block)
	chn := make(chan t.Block)
	timerCh := make(chan string)

	go Timer(timerCh, n.PeerAddrs)
	log.Default().Println("Both Chanells Created")

	transactionPool := t.TransactionPool{} // temporary transaction pool
	for {
		select {
		case miner := <-timerCh: // string value of miner
			log.Default().Println("Miner selected", miner)

			// write in a global veriable or in expectedMiner channel

			log.Default().Println("Miner selected", miner)
			if strings.Compare(miner, "12D3KooWPot5PSrTg6K") == 0 {
				go m.MineBlock(chn, &transactionPool)
			}
		case blk := <-chn: // getting mined block
			log.Default().Println("Mined Block: ", blk)
			log.Default().Println(blk)
		}
	}
}

func Timer(timerCh chan string, miners []string) {
	// find miner index
	log.Default().Println("Timer started")
	index := 0
	numberOfMiners := len(miners)
	time.Sleep(time.Duration(0) * time.Second) // set time according to last block miner
	timerCh <- miners[index]
	index = (index + 1) % numberOfMiners

	for {
		time.Sleep(time.Duration(2) * time.Second)
		timerCh <- miners[index]
		index = (index + 1) % numberOfMiners
		log.Default().Println(miners[index])
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
