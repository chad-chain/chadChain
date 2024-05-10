package main

import (
	"context"
	"log"
	"time"

	"github.com/chad-chain/chadChain/core/crypto"
	"github.com/chad-chain/chadChain/core/initialize"
	m "github.com/chad-chain/chadChain/core/mining"
	n "github.com/chad-chain/chadChain/core/network"
	db "github.com/chad-chain/chadChain/core/storage"
	t "github.com/chad-chain/chadChain/core/types"
	"github.com/chad-chain/chadChain/core/validator"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	n.CtxVar = context.Background()
	// n.PeerAddrs = []string{
	// 	// "/ip4/127.0.0.1/tcp/63795/p2p/12D3KooWBMNwiqwM1DhRXDaTXU2CYCdpKMDY5tNxfviq7ogFmFhW",
	// }
	go func() {
		n.Run()
	}()
	db.InitBadger()
	defer db.BadgerDB.Close()
	initialize.GlobalDBVar()
	initialize.Keys()
	initialize.InitFaucet()

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
	time.Sleep(5 * time.Second)
	expectedMiners := make(chan string)
	m.MiningInit(expectedMiners, &n.PeerAddrs, n.GetHostAddr()[1])
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

	signedTx, err := crypto.SignTransaction(&transaction, privateKey)

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

	if validator.ValidateHeader(&block) {
		log.Default().Println("Block is valid")
	} else {
		log.Default().Println("Block is invalid")
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
