package mining

import (
	"log"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/ethereum/go-ethereum/crypto"
	cry "github.com/malay44/chadChain/core/crypto"
	"github.com/malay44/chadChain/core/storage"
	t "github.com/malay44/chadChain/core/types"
	rlp "github.com/malay44/chadChain/core/utils"
)

func createEmptyBlock() t.Block {
	emptyBlock := new(t.Block)
	emptyBlock.Header.ParentHash = t.LatestBlockHash
	emptyBlock.Header.StateRoot = [32]byte(t.StateRootHash)
	emptyBlock.Header.TransactionsRoot = [32]byte(t.LatestBlock.Header.TransactionsRoot)
	emptyBlock.Header.Number = t.LatestBlock.Header.Number + 1
	emptyBlock.Header.Timestamp = uint64(time.Now().Unix())
	return *emptyBlock
}

func ExecuteTransaction(transaction *t.Transaction, txn *badger.Txn) error {
	senderAddress, err := cry.VerifyTxSign(transaction)
	if err != nil {
		return err
	}
	receiverAccount, err := t.GetAccount(transaction.To)
	if err != nil {
		return err
	}
	senderAccount, err := t.GetAccount(senderAddress)
	if err != nil {
		return err
	}

	senderAccount.Balance -= transaction.Value
	receiverAccount.Balance += transaction.Value

	// This is done to maintain the atomicity of the transaction
	// If any of the updates fail, the transaction will be rolled back
	// and the error will be returned
	updateFn, err := receiverAccount.UpdateAccount()
	if err != nil {
		return err
	}
	err = updateFn(txn)
	if err != nil {
		return err
	}
	updateFn, err = senderAccount.UpdateAccount()
	if err != nil {
		return err
	}
	err = updateFn(txn)
	if err != nil {
		return err
	}
	return nil
}

// MineBlock mines a block and adds it to the channel chn.
// Don't call this function directly from the main function
// spin up a go routine to call this function
func MineBlock(chn chan t.Block, transactionPool *t.TransactionPool) {
	log.Default().Println("Building Block Function")
	transactions := transactionPool.Get_all_transactions()
	log.Default().Println("Transactions: ", transactions)
	txn := storage.BadgerDB.NewTransaction(true)
	log.Default().Println("txn: ", txn)
	defer txn.Discard()

	// Execute all the transactions in the transaction pool
	for _, tx := range transactions {
		err := ExecuteTransaction(&tx, txn)
		if err != nil {
			txn.Discard()
			chn <- createEmptyBlock()
			log.Fatalln("Error executing transaction: ", err)
		}
	}

	err := t.ComputeAndSaveRootHash()
	if err != nil {
		txn.Discard()
		chn <- createEmptyBlock()
		log.Fatalln("Error computing and saving root hash: ", err)
	}

	transactionRLP, err := rlp.EncodeData(transactions, false)
	if err != nil {
		txn.Discard()
		chn <- createEmptyBlock()
		log.Fatalln("Error encoding the transactions so creating empty block:  ", err)
	}

	header := t.Header{
		ParentHash:       t.LatestBlockHash,
		Miner:            cry.MinerAddress,
		StateRoot:        [32]byte(t.StateRootHash),
		TransactionsRoot: [32]byte(crypto.Keccak256(transactionRLP)),
		Number:           t.LatestBlock.Header.Number + 1,
		Timestamp:        uint64(time.Now().Unix()),
	}
	cry.SignHeader(&header)

	b := t.Block{
		Header:       header,
		Transactions: transactions,
	}
	// TODO: if adding block to chain fails, we should add them back to the transaction pool
	if err != nil {
		txn.Discard()
		chn <- createEmptyBlock()
		log.Fatalln("Error adding block to chain: ", err)
	}

	err = txn.Commit()
	if err != nil {
		txn.Discard()
		chn <- createEmptyBlock()
		log.Fatalln("Error committing transaction: ", err)
	}
	transactionPool.RemoveCommonTransactions(transactions)
	chn <- b
}
