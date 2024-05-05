package main

import (
	"log"

	t "github.com/malay44/chadChain/core/types"
)

func ValidateTransaction(tr *t.Transaction) bool {
	sender, err := t.VerifySign(tr)
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
		log.Default().Println("Sender:", t.BytesToHexString(sender[:]))
		return false
	}

	return true
}

func ValidateBlock(b t.Block) bool {
	// validate block
	return true
}
