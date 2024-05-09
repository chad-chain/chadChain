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

<<<<<<< HEAD
func BuildBlock(chn chan t.Block, transactionPool *t.TransactionPool, Miner string) (t.Block, error) {
=======
func createEmptyBlock() t.Block {
	emptyBlock := new(t.Block)
	emptyBlock.Header.ParentHash = t.GetParentHash()
	emptyBlock.Header.StateRoot = [32]byte(t.StateRootHash)
	emptyBlock.Header.TransactionsRoot = [32]byte(t.GetParentBlockTransactionsRoot())
	emptyBlock.Header.Number = t.GetParentBlockHeight() + 1
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
>>>>>>> a98a04aa5f3612e6e545980a86d8d5a3fe67c7fa
	log.Default().Println("Building Block Function")
	transactions := transactionPool.Get_all_transactions_and_clear()
	txn := storage.BadgerDB.NewTransaction(true)
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

<<<<<<< HEAD
	transactionRootBytes := [32]byte{}
	copy(transactionRootBytes[:], transactionroot)
	miner, err := rlp.EncodeData(Miner, false)
	minerBytes := [20]byte{}
	copy(minerBytes[:], miner)
	header := new(t.Header)
	header.ParentHash = t.GetParentHash()
	header.Miner = minerBytes
	// header.StateRoot = [32]byte(t.StateRootHash)
	header.TransactionsRoot = transactionRootBytes
	header.Number = t.GetParentBlockHeight() + 1
	header.Timestamp = uint64(timestamp)

	b := t.CreateBlock(*header, transactionPool.Get_all_transactions_and_clear())
	log.Default().Println("Block Created", b)
	minedBlock, err := MineBlock(b)
	if err != nil {
		log.Default().Println("Error in encoding transaction")
		emptyBlock := new(t.Block)
		emptyBlock.Header.ParentHash = t.GetParentHash()
		emptyBlock.Header.StateRoot = [32]byte(t.StateRootHash)
		emptyBlock.Header.TransactionsRoot = [32]byte(t.GetParentBlockTransactionsRoot())
		emptyBlock.Header.Number = t.GetParentBlockHeight() + 1
		emptyBlock.Header.Timestamp = uint64(timestamp)
		emptyBlock.Header.Miner = minerBytes
		return *emptyBlock, err
=======
	err := t.ComputeAndSaveRootHash()
	if err != nil {
		txn.Discard()
		chn <- createEmptyBlock()
		log.Fatalln("Error computing and saving root hash: ", err)
>>>>>>> a98a04aa5f3612e6e545980a86d8d5a3fe67c7fa
	}

	transactionRLP, err := rlp.EncodeData(transactions, false)
	if err != nil {
		txn.Discard()
		chn <- createEmptyBlock()
		log.Fatalln("Error encoding the transactions so creating empty block:  ", err)
	}

	header := t.Header{
		ParentHash:       t.GetParentHash(),
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
	chn <- b
}
