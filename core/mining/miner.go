package mining

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	cry "github.com/chad-chain/chadChain/core/crypto"
	"github.com/chad-chain/chadChain/core/network"
	"github.com/chad-chain/chadChain/core/storage"
	t "github.com/chad-chain/chadChain/core/types"
	rlp "github.com/chad-chain/chadChain/core/utils"
	"github.com/chad-chain/chadChain/core/validator"
	"github.com/dgraph-io/badger/v4"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	expectedMiner chan string
)

func createEmptyBlock() t.Block {
	emptyBlock := new(t.Block)
	emptyBlock.Header.ParentHash = t.LatestBlockHash
	emptyBlock.Header.StateRoot = [32]byte(t.StateRootHash)
	emptyBlock.Header.TransactionsRoot = [32]byte(t.LatestBlock.Header.TransactionsRoot)
	emptyBlock.Header.Number = t.LatestBlock.Header.Number + 1
	emptyBlock.Header.Timestamp = uint64(time.Now().Unix())
	sig, err := cry.SignHeader(&emptyBlock.Header)
	if err != nil {
		log.Fatalln("Failed to sign header: ", err)
	}
	emptyBlock.Header.ExtraData = sig
	return *emptyBlock
}

func ExecuteTransaction(transaction *t.Transaction, txn *badger.Txn) error {
	fmt.Println("Executing Transaction")
	senderAddress, err := cry.VerifyTxSign(transaction)
	if err != nil {
		return err
	}
	receiverAccount, err := t.GetAccount(transaction.To)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			receiverAccount = t.Account{
				Address: transaction.To,
				Nonce:   0,
				Balance: 0,
			}
			_, _, err := receiverAccount.AddAccount()
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	senderAccount, err := t.GetAccount(senderAddress)
	if err != nil {
		return err
	}

	senderAccount.Balance -= transaction.Value
	receiverAccount.Balance += transaction.Value

	fmt.Println("Sender Account Balance: ", senderAccount.Balance)
	fmt.Println("Receiver Account Balance: ", receiverAccount.Balance)

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
	fmt.Println("Building Block Function")
	transactions := transactionPool.Get_all_transactions()
	fmt.Println("Mining Block with transactions:")
	txn := storage.BadgerDB.NewTransaction(true)
	defer txn.Discard()

	// Execute all the transactions in the transaction pool
	for i, tx := range transactions {
		fmt.Println("Transaction ", i, " : ")
		fmt.Println("Transaction:")
		fmt.Println("To:", tx.To.Hex())
		fmt.Println("Value:", tx.Value)
		fmt.Println("Nonce:", tx.Nonce)
		err := ExecuteTransaction(&tx, txn)
		if err != nil {
			txn.Discard()
			chn <- createEmptyBlock()
			fmt.Println("Error executing transaction: ", err)
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
	sig, err := cry.SignHeader(&header)
	if err != nil {
		txn.Discard()
		chn <- createEmptyBlock()
		log.Fatalln("Failed to sign header: ", err)
	}
	header.ExtraData = sig

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

// func to add block received from other miners to the blockchain
func AddBlockToChain(block t.Block) {
	// check if the block is valid
	if !validator.VerifyBlock(&block) {
		fmt.Println("Block is invalid")
		return
	}
	txn := storage.BadgerDB.NewTransaction(true)
	defer txn.Discard()
	for _, tx := range block.Transactions {
		err := ExecuteTransaction(&tx, txn)
		if err != nil {
			txn.Discard()
			fmt.Println("Error executing transaction: ", err)
			return
		}
	}
	err := t.ComputeAndSaveRootHash()
	if err != nil {
		txn.Discard()
		fmt.Println("Error computing and saving root hash: ", err)
		return
	}
	if block.Header.StateRoot != [32]byte(t.StateRootHash) {
		fmt.Println("State root hash is not correct")
		return
	}
	if block.Header.Number != t.LatestBlock.Header.Number+1 {
		fmt.Println("Block number is not correct")
		return
	}
	err = txn.Commit()
	if err != nil {
		txn.Discard()
		fmt.Println("Error committing transaction: ", err)
		return
	}

	err = block.PersistBlock()
	if err != nil {
		fmt.Println("Error persisting block: ", err)
		return
	} else {
		fmt.Println("Block added to chain")
	}
}

func MiningInit(expectedMiner chan string, peerAddrs *[]string, selfAdr string) { // add transactionpool as argument

	// ch := make(chan t.Block)
	chn := make(chan t.Block)
	timerCh := make(chan string)

	go Timer(timerCh, peerAddrs)
	fmt.Println("Both Chanells Created")

	for {
		select {
		case miner := <-timerCh: // string value of miner
			fmt.Println("Miner selected", miner)

			// write in a global veriable or in expectedMiner channel
			if strings.Compare(miner, selfAdr) == 0 {
				go MineBlock(chn, t.TransactionPoolVar)
			}
		case blk := <-chn: // getting mined block
			fmt.Println("Mined Block: ")
			blk.Print()
			network.SendBlock(blk)
			blk.PersistBlock()
		}
	}
}

func Timer(timerCh chan string, miners *[]string) {
	for {
		if len(*miners) > 1 {
			break
		}
	}
	sort.Strings(*miners)
	fmt.Println("Timer started")
	index := len(*miners) - 1
	numberOfMiners := len(*miners)
	time.Sleep(time.Duration(0) * time.Second)
	timerCh <- (*miners)[index]

	for {
		numberOfMiners = len(*miners)        // Update the number of miners
		index = (index + 1) % numberOfMiners // Calculate the index
		timerCh <- (*miners)[index]
		time.Sleep(time.Duration(10) * time.Second)
	}
}
