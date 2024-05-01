package types

import (
	"log"

	"github.com/dgraph-io/badger/v4"
	"github.com/malay44/chadChain/core/storage"
	"github.com/vmihailenco/msgpack/v5"
)

type Account struct {
	Address [32]byte // The address of the account
	Nonce   uint64   // The nonce of the account
	Balance uint64   // The balance of the account
}

// State root hash
var StateRootHash [32]byte

func (ac *Account) CreateAccount(address [32]byte, nonce uint64, balance uint64) Account {
	return Account{address, nonce, balance}
}

// Get account from network and save to db
func (ac *Account) AddAccount() (string, string, error) {

	// create keys for account and its hash
	addrSlice := ac.Address[:]
	accKey := "account" + string(addrSlice)
	hashKey := "hash" + string(addrSlice)

	// marshal account info to byte array and hash it
	val, err := msgpack.Marshal(ac)
	if err != nil {
		panic(err)
	}
	hash := Keccak256(val)

	// save account info and its hash to db
	err = storage.BadgerDB.Update(func(tx *badger.Txn) error {
		err := storage.Insert([]byte(accKey), ac)(tx)
		if err != nil {
			return err
		}
		err = storage.Insert([]byte(hashKey), hash)(tx)
		if err != nil {
			return err
		}
		return nil
	})

	// check for errors
	if err != nil {
		log.Default().Println(err.Error())
		return "", "", err
	} else {
		log.Default().Println("Account saved to db with key", accKey, "and hash key", hashKey)
		return accKey, hashKey, nil
	}
}

// Get account from db
func GetAccount(address string) (Account, error) {
	// create keys for account and its hash
	accKey := "account"
	acc := Account{}
	err := storage.BadgerDB.View(func(tx *badger.Txn) error {
		err := storage.Get([]byte(accKey+address), &acc)(tx)
		return err
	})

	if err != nil {
		return acc, err
	}

	return acc, nil
}

// send account over network
func (ac *Account) SendAccount(Account Account) {
	// propagate in the network
}
