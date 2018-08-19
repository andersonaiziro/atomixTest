package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"strings"
)

// Transaction represents a new UTXO transaction between nodes
type Transaction struct {
	ID       [32]byte
	TxInput  []TXInput
	TxOutput []TXOutput
	RandVal  int
}

// Hash returns the hash of the Transaction
func (tx *Transaction) Hash() [32]byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = [32]byte{}

	hash = sha256.Sum256(txCopy.Serialize())

	return hash
}

// GetTransactionLockKey Retrieves an Unique Key for a specific OutputIndex in a Transaction
func GetTransactionLockKey(txID [32]byte) string {
	b := md5.Sum(txID[:])
	lockKey := hex.EncodeToString(b[:])

	return lockKey
}

// Serialize returns a serialized Transaction
func (tx Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

// NewCoinbaseTX creates a new coinbase transaction
func NewCoinbaseTX(to string, shardID uint32, value int) *Transaction {

	txin := TXInput{[32]byte{}, -1, "", 0}
	txout := NewTXOutput(value, to, shardID)
	tx := Transaction{[32]byte{}, []TXInput{txin}, []TXOutput{*txout}, rand.Intn(100000000)}
	tx.ID = tx.Hash()

	return &tx
}

// IsFromCoinbase is an auxiliary method to avoid validating genesis transactions
func (tx *Transaction) IsFromCoinbase() bool {
	isFromCoinbase := false

	for _, txIns := range tx.TxInput {
		if txIns.ShardID == 0 {
			isFromCoinbase = true
		}
	}

	return isFromCoinbase
}

// NewDistributionTXs creates a new coinbase transaction to initialize the Shard
func NewDistributionTXs(toAddress []string, shardID uint32, value int, txoutCount int) *Transaction {

	txin := TXInput{[32]byte{}, -1, "", 0}

	var txouts []TXOutput

	for i := 0; i < txoutCount; i++ {
		for _, address := range toAddress {
			txouts = append(txouts, *NewTXOutput(value, address, shardID))
		}
	}

	tx := Transaction{[32]byte{}, []TXInput{txin}, txouts, rand.Intn(100000000)}
	tx.ID = tx.Hash()

	return &tx
}

// NewUTXOTransaction creates a new transaction
func NewUTXOTransaction(utxo *UTXO, toAddress string, toShard uint32, amount int) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	inputs = append(inputs, TXInput{utxo.TxID, utxo.OutputIndex, utxo.Address, utxo.ShardID})
	outputs = append(outputs, TXOutput{amount, toAddress, toShard})

	tx := Transaction{[32]byte{}, inputs, outputs, rand.Intn(100000000)}
	tx.ID = tx.Hash()

	return &tx
}

// String returns a human-readable representation of a transaction
func (tx Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.ID))

	for i, input := range tx.TxInput {
		lines = append(lines, fmt.Sprintf("     Input: %d, Address: %s, Out: %d, ShardId: %d, TxID: %x", i, input.Address, input.OutputIndex, input.ShardID, input.TxID))
	}

	for i, output := range tx.TxOutput {
		lines = append(lines, fmt.Sprintf("     Output: %d, Value: %d, ShardID: %d, Address: %s", i, output.Value, output.ShardID, output.Address))
	}

	return strings.Join(lines, "\n")
}
