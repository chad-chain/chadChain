package initialize

import (
	"fmt"
	"log"
	"os"

	"github.com/chad-chain/chadChain/core/crypto"
	db "github.com/chad-chain/chadChain/core/storage"
	t "github.com/chad-chain/chadChain/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
)

func GlobalDBVar() {
	err := db.BadgerDB.View(db.Get([]byte("stateRootHash"), &t.StateRootHash))
	if err != nil {
		log.Default().Printf("Error in init StateRootHash\n")
		log.Default().Println(err.Error())
	} else {
		log.Default().Printf("StateRootHash: %x\n", t.StateRootHash)
	}
	err = db.BadgerDB.View(db.Get([]byte("latestBlock"), &t.LatestBlock))
	if err != nil {
		log.Default().Printf("Error in init latestBlock: \n")
		log.Default().Println(err.Error())
	} else {
		log.Default().Printf("LatestBlock: %x\n", t.LatestBlock)
	}
}

func Keys() {
	godotenv.Load(".env")
	if os.Getenv("PRIV_HEX") == "" {
		// create a new private key
		_, hex, Addr, err := crypto.GenerateNewPrivateKey()
		if err != nil {
			log.Fatalln("Failed to generate private key")
			return
		}
		os.Setenv("PRIV_HEX", hex)
		os.Setenv("WALLET_ADDR", Addr.Hex())
		log.Default().Println("Private key and Address generated")
		log.Default().Println("Save the following in your env file")
		fmt.Println("---------------------------------------------------------------------------------------------------")
		fmt.Printf("PRIV_HEX=%s\n", hex)
		fmt.Printf("WALLET_ADDR=%s\n", Addr.Hex())
		fmt.Println("---------------------------------------------------------------------------------------------------")

		// load the private
	} else {
		crypto.PrivateKeyHex = os.Getenv("PRIV_HEX")
		FaucetPrivateKeyHex := os.Getenv("FAUCET_PRIV_HEX")
		FaucetPrivateKeyHex2 := os.Getenv("FAUCET_PRIV_HEX2")
		_, err := crypto.LoadPrivateKeyAndAddr(crypto.PrivateKeyHex)
		if err != nil {
			log.Fatalln("Failed to load private key")
			return
		}
		_, err = crypto.LoadFaucetPrivateKeyAndAddr(FaucetPrivateKeyHex2)
		if err != nil {
			log.Fatalln("Failed to load faucet private key")
			return
		}
		_, err = crypto.LoadFaucetPrivateKeyAndAddr(FaucetPrivateKeyHex)
		if err != nil {
			log.Fatalln("Failed to load faucet private key")
			return
		}

		log.Default().Println("Private key and Address loaded")
		fmt.Println("---------------------------------------------------------------------------------------------------")
		fmt.Println("Private key: ", crypto.PrivateKeyHex)
		fmt.Println("Wallet Address: ", crypto.MinerAddress.Hex())
		fmt.Println("---------------------------------------------------------------------------------------------------")
	}
}

func InitFaucet() {
	var toAddress [20]byte
	copy(toAddress[:], crypto.HexStringToBytes(os.Getenv("FAUCET_WALLET_ADDR")))
	fmt.Println("Value of address:", toAddress)
	acc := t.Account{
		Address: toAddress,
		Nonce:   0,
		Balance: 100,
	}
	_, _, err := acc.AddAccount()
	if err != nil {
		log.Default().Println("Failed to initialize faucet account")
		return
	}
	log.Default().Println("Faucet1 account initialized")
	toAddress2 := common.HexToAddress(os.Getenv("FAUCET_WALLET_ADDR2"))
	acc2 := t.Account{
		Address: [20]byte(toAddress2),
		Nonce:   0,
		Balance: 100,
	}
	_, _, err = acc2.AddAccount()
	if err != nil {
		log.Default().Println("Failed to initialize faucet account")
		return
	}
	log.Default().Println("Faucet2 account initialized")
}
