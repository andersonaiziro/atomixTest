package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

// Block represents a block in the blockchain
type Block struct {
	PrevBlockHash   [32]byte
	NumTransactions int32
	TransactionIds  [][32]byte
	Transactions    []*Transaction
	ShardId         uint32   // The shard this block belongs to
	Hash            [32]byte // The shard Id this block belongs to
}

// NewBlock creates and returns Block
func NewBlock(transactions []*Transaction, prevBlockHash [32]byte, shardID uint32) *Block {
	var transactionIds [][32]byte
	for _, transaction := range transactions {
		transactionIds = append(transactionIds, transaction.ID)
	}
	block := &Block{prevBlockHash, int32(len(transactions)), transactionIds, transactions, shardID, [32]byte{}}

	block.Hash = block.GetBlockHash()
	return block
}

// NewGenesisBlock creates and returns genesis Block
func NewGenesisBlock(coinbaseTXs []*Transaction) *Block {
	return NewBlock(coinbaseTXs, [32]byte{}, 0)
}

// GetBlockHash serializes the block
func (b Block) GetBlockHash() [32]byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return sha256.Sum256(result.Bytes())
}
