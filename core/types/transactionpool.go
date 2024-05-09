package types

import (
	"reflect"
	"sync"
)

type TransactionPool struct {
	Transactions []Transaction
	mutex        sync.Mutex
}

var TransactionPoolVar = NewTransactionPool()

func NewTransactionPool() *TransactionPool {
	return &TransactionPool{}
}

// Method to add a transaction to the transaction pool
func (tp *TransactionPool) Add_transaction_to_transactionPool(tr Transaction) *TransactionPool {

	tp.mutex.Lock()
	defer tp.mutex.Unlock() // Ensure unlocking even if there's a panic
	tp.Transactions = append(tp.Transactions, tr)
	return tp // Return the updated transaction pool
}

// Method to get all transactions from the transaction pool and clear it
func (tp *TransactionPool) Get_all_transactions_and_clear() []Transaction {

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

func (tp *TransactionPool) Get_all_transactions() []Transaction {

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
func (tp *TransactionPool) AddTransactions(transactions []Transaction) *TransactionPool {
	tp.mutex.Lock()
	defer tp.mutex.Unlock()
	tp.Transactions = append(tp.Transactions, transactions...)
	return tp
}

// clear transactionPool teansaction and add list of transactions
func (tp *TransactionPool) ClearAndAddTransactions(transactions []Transaction) *TransactionPool {
	tp.mutex.Lock()
	defer tp.mutex.Unlock()
	tp.Transactions = nil
	tp.Transactions = append(tp.Transactions, transactions...)
	return tp
}

// remove transaction from transaction pool which are in the valid block
func (tp *TransactionPool) RemoveCommonTransactions(transactions []Transaction) *TransactionPool {

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

func (tp *TransactionPool) Print() {
	tp.mutex.Lock()
	defer tp.mutex.Unlock()
	for i, t := range tp.Transactions {
		println("Transaction", i)
		t.Print()
	}
}
