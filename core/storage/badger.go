package storage

import (
	"errors"
	"log"

	"github.com/dgraph-io/badger/v4"
	"github.com/vmihailenco/msgpack/v5"
)

var BadgerDB *badger.DB

func InitBadger() {
	db, err := badger.Open(badger.DefaultOptions("./database/block"))
	if err != nil {
		log.Fatal(err)
	}
	BadgerDB = db
	log.Default().Println("BadgerDB initialized")
}

func Insert(key []byte, value interface{}) func(*badger.Txn) error {
	return func(txn *badger.Txn) error {
		// check if the key already exists in the db
		_, err := txn.Get(key)
		if err == nil {
			return errors.New("key already exists")
		}

		val, err := msgpack.Marshal(value)
		if err != nil {
			return err
		}

		err = txn.Set(key, val)
		if err != nil {
			return err
		}

		return nil
	}
}

func Get(key []byte, entity interface{}) func(*badger.Txn) error {
	return func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			err := msgpack.Unmarshal(val, entity)
			return err
		})
		return err
	}
}
