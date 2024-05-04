package storage

import (
	"fmt"
	"log"

	"github.com/dgraph-io/badger/v4"
	rlp "github.com/malay44/chadChain/core/utils"
)

// CheckFunc is called during key iteration through the badger DB in order to
// check whether we should process the given key-value pair. It can be used to
// avoid loading the value if its not of interest, as well as storing the key
// for the current iteration step.
type CheckFunc func(key []byte) bool

// CreateFunc returns a pointer to an initialized entity that we can potentially
// decode the next value into during a badger DB iteration.
type CreateFunc func() interface{}

// HandleFunc is a function that starts the processing of the current key-value
// pair during a badger iteration. It should be called after the key was checked
// and the entity was decoded.
// No errors are expected during normal operation. Any errors will halt the iteration.
type HandleFunc func() error

// IterationFunc is a function provided to our low-level iteration function that
// allows us to pass badger efficiencies across badger boundaries. By calling it
// for each iteration step, we can inject a function to check the key, a
// function to create the decode target and a function to process the current
// key-value pair. This a consumer of the API to decode when to skip the loading
// of values, the initialization of entities and the processing.
type IterationFunc func() (CheckFunc, CreateFunc, HandleFunc)

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
			return fmt.Errorf("badger.go/Insert: key already exists")
		}

		val, err := rlp.EncodeData(value, false)
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

func Update(key []byte, value interface{}) func(*badger.Txn) error {
	return func(txn *badger.Txn) error {
		_, err := txn.Get(key)
		if err != nil {
			return fmt.Errorf("badger.go/Update: key not found")
		}

		val, err := rlp.EncodeData(value, false)
		if err != nil {
			return fmt.Errorf("badger.go/Update: could not marshal value: %v", err)
		}

		err = txn.Set(key, val)
		if err != nil {
			return fmt.Errorf("badger.go/Update: could not set value: %v", err)
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
			err := rlp.DecodeData(val, entity)
			return err
		})
		return err
	}
}

// traverse iterates over a range of keys defined by a prefix.
//
// The prefix must be shared by all keys in the iteration.
//
// On each iteration, it will call the iteration function to initialize
// functions specific to processing the given key-value pair.
func Traverse(prefix []byte, iteration IterationFunc) func(*badger.Txn) error {
	return func(tx *badger.Txn) error {
		if len(prefix) == 0 {
			return fmt.Errorf("prefix must not be empty")
		}

		opts := badger.DefaultIteratorOptions
		// NOTE: this is an optimization only, it does not enforce that all
		// results in the iteration have this prefix.
		opts.Prefix = prefix

		it := tx.NewIterator(opts)
		defer it.Close()

		// this is where we actually enforce that all results have the prefix
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {

			item := it.Item()

			// initialize processing functions for iteration
			check, create, handle := iteration()

			// check if we should process the item at all
			key := item.Key()
			ok := check(key)
			if !ok {
				continue
			}

			// process the actual item
			err := item.Value(func(val []byte) error {

				// decode into the entity
				entity := create()
				err := rlp.DecodeData(val, entity)
				if err != nil {
					return fmt.Errorf("could not unmarshal entity: %w", err)
				}

				// process the entity
				err = handle()
				if err != nil {
					return fmt.Errorf("could not handle entity: %w", err)
				}

				return nil
			})
			if err != nil {
				return fmt.Errorf("could not process value: %w", err)
			}
		}

		return nil
	}
}
