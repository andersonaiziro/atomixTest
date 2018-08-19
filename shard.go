package main

import (
	"fmt"
	"log"
	"strings"
)

//Shard defines a node in the network
type Shard struct {
	ShardID          uint32
	Blockchain       *Blockchain
	NextTransactions []Transaction
	//maps between a transactionKey and the UTXOs belonging to that node
	Mempool map[string][]*UTXO
}

// NewShard creates a new shard for a given shardID
func NewShard(shardID uint32, blockchain *Blockchain) *Shard {
	shard := &Shard{shardID, blockchain, []Transaction{}, map[string][]*UTXO{}}
	return shard
}

// ProcessNewTransaction receives a new transaction request
func (shard *Shard) ProcessNewTransaction(inTransaction <-chan *Transaction, c chan<- Gossips) {
	transaction := <-inTransaction

	if transaction.IsFromCoinbase() {
		shard.internalCommitTransaction(transaction)
		c <- ProofOfAcceptance
		return
	}

	if !shard.validateTXIns(transaction) {
		c <- ProofOfRejection
		return
	}

	// If transaction uses valid UTXOs, lock UTXOs needed to fund the transaction
	shard.lockTransactionUTXOs(transaction)

	// If reached this point, UTXOs are locked, and
	c <- ProofOfAcceptance
}

// CommitTransaction is called when all shards agree the transaction should be commited and funds spent / received
func (shard *Shard) CommitTransaction(inTransaction <-chan *Transaction, c chan<- Gossips) {
	transaction := <-inTransaction
	shard.internalCommitTransaction(transaction)
	c <- TransactionCommited
}

func (shard *Shard) internalCommitTransaction(transaction *Transaction) {
	// Process All the TXIns and TXOuts that use this Shard
	if shard.updateMemPoolCommit(transaction) {
		shard.NextTransactions = append(shard.NextTransactions, *transaction)
	}

	// Check if a new block should be mined
	if shard.shouldCreateNewBlock() {
		shard.mineNewBlock()
	}
}

// AbortTransaction is called when the transaction should be aborted
func (shard *Shard) AbortTransaction(inTransaction <-chan *Transaction, c chan<- Gossips) {
	transaction := <-inTransaction

	for _, tx := range shard.getShardTXIns(transaction) {

		lockupKey := GetTransactionLockKey(tx.TxID)
		_, utxo := shard.getMempoolUTXO(tx.TxID, tx.OutputIndex, tx.Address, lockupKey)

		//Unlock the utxo to the address owner
		if utxo != nil {
			utxo.Address = tx.Address
		}
	}

	c <- TransactionCommited
}

func (shard *Shard) updateMemPoolCommit(transaction *Transaction) bool {
	hasMempoolChanged := false

	//check all the TXIn UTXOs referenced by this transaction and delete them from the mempool... they have been spent
	for _, tx := range shard.getShardTXIns(transaction) {
		hasMempoolChanged = true

		index, _ := shard.getMempoolUTXO(tx.TxID, tx.OutputIndex, tx.Address, GetTransactionLockKey(transaction.ID))

		shard.removeUTXO(tx.Address, index)
	}

	//check all the TXOut UTXOs referenced by this transaction and add them to the mempool...
	for _, txoutIndex := range shard.getShardTXOutsIndexes(transaction) {
		hasMempoolChanged = true
		addressKey, newUTXO := NewUTXO(*transaction, txoutIndex)
		shard.Mempool[addressKey] = append(shard.Mempool[addressKey], newUTXO)
	}

	return hasMempoolChanged
}

func (shard *Shard) removeUTXO(address string, indexToRemove int) {
	s := shard.Mempool[address]
	s[len(s)-1], s[indexToRemove] = s[indexToRemove], s[len(s)-1]
	shard.Mempool[address] = s[:len(s)-1]
}

func (shard *Shard) lockTransactionUTXOs(transaction *Transaction) {
	//check all the UTXOs referenced by this transaction can actually be spent
	for _, tx := range shard.getShardTXIns(transaction) {

		_, utxo := shard.getMempoolUTXO(tx.TxID, tx.OutputIndex, tx.Address, tx.Address)

		//Lock the UTXO to this transaction
		utxo.Address = GetTransactionLockKey(transaction.ID)
	}
}

func (shard *Shard) validateTXIns(transaction *Transaction) bool {
	//check all the UTXOs referenced by this transaction can actually be spent
	for _, tx := range shard.getShardTXIns(transaction) {
		_, utxo := shard.getMempoolUTXO(tx.TxID, tx.OutputIndex, tx.Address, tx.Address)

		if utxo == nil {
			return false
		}
	}

	return true
}

//retrieves UTXO referenced by transaction
func (shard *Shard) getMempoolUTXO(txID [32]byte, outputIndex int, utxoAddress, lockAddress string) (int, *UTXO) {
	utxos, isKeyPresent := shard.Mempool[utxoAddress]

	if !isKeyPresent {
		return -1, nil
	}

	for index, utxo := range utxos {
		if (utxo.TxID == txID) && (utxo.OutputIndex == outputIndex) && (utxo.Address == lockAddress) {
			return index, utxo
		}
	}

	return -1, nil
}

func (shard *Shard) getShardTXOutsIndexes(transaction *Transaction) []int {
	var shardTXOutIndexes []int
	for txIndex, tx := range transaction.TxOutput {
		//transaction input to another shard... nothing to do
		if tx.ShardID != shard.ShardID {
			continue
		}

		shardTXOutIndexes = append(shardTXOutIndexes, txIndex)
	}

	return shardTXOutIndexes
}

func (shard *Shard) getShardTXIns(transaction *Transaction) []TXInput {
	var shardTXIns []TXInput
	for _, tx := range transaction.TxInput {
		//transaction input to another shard... nothing to do
		if tx.ShardID != shard.ShardID {
			continue
		}

		shardTXIns = append(shardTXIns, tx)
	}

	return shardTXIns
}

// mineNewBlock will create a new block for a shard that has 10 transactions
func (shard *Shard) mineNewBlock() {
	var transactionsToCommit []*Transaction
	for _, transaction := range shard.NextTransactions[:10] {
		tx := transaction
		transactionsToCommit = append(transactionsToCommit, &tx)
	}

	newBlock := shard.Blockchain.MineBlock(transactionsToCommit, shard.ShardID)

	shard.Blockchain.AddBlock(newBlock)

	if len(shard.NextTransactions) > 10 {
		shard.NextTransactions = shard.NextTransactions[10:]
	} else {
		shard.NextTransactions = []Transaction{}
	}
}

func (shard *Shard) shouldCreateNewBlock() bool {
	if len(shard.NextTransactions) == 10 {
		log.Println(fmt.Sprintf("Shard %d ready to mine new block", shard.ShardID))
		return true
	}

	return false
}

// GetMempoolString returns a string representation of a shard mempool
func (shard *Shard) GetMempoolString() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Shard %x:", shard.ShardID))

	for address, utxos := range shard.Mempool {
		lines = append(lines, fmt.Sprintf("UTXOs for Address: %s", address))
		for _, utxo := range utxos {
			if utxo.Address != address {
				lines = append(lines, fmt.Sprintf("LOCKED --- TxID: %xAddress: %s", utxo.TxID, utxo.Address))
			}
			lines = append(lines, fmt.Sprintf("TxID: %x, OutIndex: %d, Address: %s, Value: %d", utxo.TxID, utxo.OutputIndex, utxo.Address, utxo.Value))
		}

	}

	return strings.Join(lines, "\n")
}

// GetBlockchainString returns a string representation of a shard blockchain
func (shard *Shard) GetBlockchainString() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Shard %x:", shard.ShardID))

	for index, block := range shard.Blockchain.Blocks {
		lines = append(lines, fmt.Sprintf("Block %d, Hash: %x. TransactionCount: %d.", index, block.Hash, block.NumTransactions))
		var txIds []string
		for _, txID := range block.TransactionIds {
			txIds = append(txIds, fmt.Sprintf("%x", txID))
		}
		lines = append(lines, fmt.Sprintf("Transactions: %s", strings.Join(txIds, ", ")))
	}

	return strings.Join(lines, "\n")
}
