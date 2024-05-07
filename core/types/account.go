package types

import (
	"fmt"
	"log"

	"github.com/dgraph-io/badger/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	s "github.com/malay44/chadChain/core/storage"
	rlp "github.com/malay44/chadChain/core/utils"
)

type Account struct {
	Address common.Address // The address of the account
	Nonce   uint64         // The nonce of the account
	Balance uint64         // The balance of the account
}

// State root hash
var StateRootHash []byte

func (ac *Account) CreateAccount(address common.Address, nonce uint64, balance uint64) Account {
	return Account{address, nonce, balance}
}

// Get account from network and save to db
func (ac *Account) AddAccount() (string, string, error) {

	// create keys for account and its hash
	addrSlice := ac.Address[:]
	accKey := "account" + string(addrSlice)
	hashKey := "hash" + string(addrSlice)

	// marshal account info to byte array and hash it
	val, err := rlp.EncodeData(ac, false)
	if err != nil {
		return "", "", err
	}
	hash := crypto.Keccak256(val)
	marshaledHash, err := rlp.EncodeData(hash, false)
	if err != nil {
		return "", "", err
	}

	// save account info and its hash to db
	err = s.BadgerDB.Update(func(tx *badger.Txn) error {
		err := s.Insert([]byte(accKey), ac)(tx)
		if err != nil {
			return err
		}
		err = s.Insert([]byte(hashKey), marshaledHash)(tx)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return "", "", err
	}

	StateRootHash, err = ComputeRootHash()
	if err != nil {
		return "", "", err
	}

	// save state root hash to db
	err = s.BadgerDB.Update(func(tx *badger.Txn) error {
		err := s.Update([]byte("stateRootHash"), StateRootHash)(tx)
		if err != nil {
			fmt.Println("error updating state root hash: ", err.Error())
			if err == badger.ErrKeyNotFound {
				err = s.Insert([]byte("stateRootHash"), StateRootHash)(tx)
				if err != nil {
					return err
				}
			}
			return err
		}
		return nil
	})

	// check for errors
	if err != nil {
		return "", "", err
	}

	// log.Default().Println("Account saved to db with \nkey", accKey, "\nhash key", hashKey, "\nstate root hash", StateRootHash)
	log.Default().Printf("Account saved to db with \nkey %s \nhash key %s \nstate root hash %x\n", accKey, hashKey, StateRootHash)
	return accKey, hashKey, nil
}

// Get account from db
func GetAccount(address common.Address) (Account, error) {
	// create keys for account and its hash
	accKey := append([]byte("account"), address.Bytes()...)
	acc := Account{}
	err := s.BadgerDB.View(func(tx *badger.Txn) error {
		err := s.Get([]byte(accKey), &acc)(tx)
		return err
	})

	if err != nil {
		return acc, err
	}

	return acc, nil
}

// Compute the root hash
func ComputeRootHash() ([]byte, error) {
	prefix := "hash"
	var hashes [][]byte
	var accHash []byte

	err := s.BadgerDB.View(s.Traverse([]byte(prefix), func() (s.CheckFunc, s.CreateFunc, s.HandleFunc) {
		checkFunc := func(_ []byte) bool {
			return true
		}
		createFunc := func() interface{} {
			return &accHash
		}
		handleFunc := func() error {
			hashes = append(hashes, accHash)
			return nil
		}
		return checkFunc, createFunc, handleFunc
	}))

	if err != nil {
		return []byte{}, err
	}

	return crypto.Keccak256(hashes...), nil

}

// send account over network
func (ac *Account) SendAccount(Account Account) {
	// propagate in the network
}
