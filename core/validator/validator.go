package validator

import (
	"log"

	"github.com/chad-chain/chadChain/core/crypto"
	t "github.com/chad-chain/chadChain/core/types"
	"github.com/chad-chain/chadChain/core/utils"
	crypto2 "github.com/ethereum/go-ethereum/crypto"
)

func ValidateTransaction(tr *t.Transaction) bool {
	sender, err := crypto.VerifyTxSign(tr)
	if err != nil {
		log.Default().Println("Failed to recover sender:", err)
		return false
	}

	acc, err := t.GetAccount(sender)
	if err != nil {
		log.Default().Println("Failed to get account:", err)
		return false
	}

	if acc.Balance < tr.Value {
		log.Default().Println("Insufficient balance")
		log.Default().Println("Tx:", tr)
		log.Default().Println("Sender:", crypto.BytesToHexString(sender[:]))
		return false
	}

	return true
}

func ValidateHeader(b *t.Block) bool {
	if !crypto.VerifyHeader(b) {
		log.Default().Println("Failed to verify header signature")
		return false
	}
	return true
}

// function to verify the received block. It verifies the block header and the transactions
func VerifyBlock(block *t.Block) bool {
	// verify the header
	if !crypto.VerifyHeader(block) {
		log.Default().Println("Failed to verify header signature")
		return false
	}

	// verify the transactions root
	// encode the transactions
	transactionRLP, err := utils.EncodeData(block.Transactions, false)
	if err != nil {
		log.Default().Println("Error encoding the transactions")
		return false
	}
	if block.Header.TransactionsRoot != [32]byte(crypto2.Keccak256(transactionRLP)) {
		log.Default().Println("Transactions root doesn't match")
		return false
	}

	// validate all the transactions in the block
	for _, tx := range block.Transactions {
		if !ValidateTransaction(&tx) {
			log.Default().Println("Invalid transaction")
			return false
		}
	}
	return true
}
