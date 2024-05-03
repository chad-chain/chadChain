package types

import (
	"reflect"
	"sync"
)

type Transactionpool struct {
	Transactions []Transaction
	mutex        sync.Mutex
}

func NewTransactionpool() *Transactionpool {
	return &Transactionpool{}
}

// Method to add a transaction to the transaction pool
func (tp *Transactionpool) Add_transaction_to_transactionpool(tr Transaction) *Transactionpool {

	tp.mutex.Lock()
	defer tp.mutex.Unlock() // Ensure unlocking even if there's a panic
	tp.Transactions = append(tp.Transactions, tr)
	return tp // Return the updated transaction pool
}

// Method to get all transactions from the transaction pool and clear it
func (tp *Transactionpool) Get_all_transactions_and_clear() []Transaction {

	tp.mutex.Lock()
	defer tp.mutex.Unlock()
	if len(tp.Transactions) == 0 {
		return nil
	}
	transactions := make([]Transaction, len(tp.Transactions)) // Copy the transactions to a new slice
	copy(transactions, tp.Transactions)
	tp.Transactions = nil // Clear the original transaction pool
	return transactions
}

func (tp *Transactionpool) Get_all_transactions() []Transaction {

	tp.mutex.Lock()
	defer tp.mutex.Unlock()
	if len(tp.Transactions) == 0 {
		return nil
	}
	transactions := make([]Transaction, len(tp.Transactions)) // Copy the transactions to a new slice
	copy(transactions, tp.Transactions)
	return transactions
}

// add list of transaction to transaction pool
func (tp *Transactionpool) AddTransactions(transactions []Transaction) *Transactionpool {
	tp.mutex.Lock()
	defer tp.mutex.Unlock()
	tp.Transactions = append(tp.Transactions, transactions...)
	return tp
}

// clear transactionpool teansaction and add list of transactions
func (tp *Transactionpool) ClearAndAddTransactions(transactions []Transaction) *Transactionpool {
	tp.mutex.Lock()
	defer tp.mutex.Unlock()
	tp.Transactions = nil
	tp.Transactions = append(tp.Transactions, transactions...)
	return tp
}

// remove transaction from transaction pool which are in the valid block
func (tp *Transactionpool) RemoveCommonTransactions(transactions []Transaction) *Transactionpool {

	transactionsFromPool := tp.Get_all_transactions()
	var updatedTransactions []Transaction

	flag := false
	for _, i := range transactionsFromPool {
		flag = false
		for _, j := range transactions {
			if reflect.DeepEqual(i, j) {
				flag = true
			}
		}
		if !flag {
			updatedTransactions = append(updatedTransactions, i)
		}
	}
	tp.ClearAndAddTransactions(updatedTransactions)
	return tp
}
