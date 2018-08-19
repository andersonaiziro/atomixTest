package main

import (
	"fmt"
	"log"
)

// Blockchain is an array of Blocks
type Blockchain struct {
	Blocks []Block
}

// NewBlockchain creates a new Blockchain with genesis Block
func NewBlockchain() *Blockchain {

	cbtx := NewCoinbaseTX("blockchainRoot", 0, 1)
	genesis := NewGenesisBlock([]*Transaction{cbtx})

	bc := Blockchain{Blocks: []Block{*genesis}}

	return &bc
}

// MineBlock mines a new block with the provided transactions
func (bc *Blockchain) MineBlock(transactions []*Transaction, shardID uint32) *Block {
	var lastHash [32]byte

	lastHash = bc.getLastBlock().Hash

	newBlock := NewBlock(transactions, lastHash, shardID)

	log.Println(fmt.Sprintf("New Block mined for Shard %d. Block hash: %x", shardID, newBlock.Hash))
	return newBlock
}

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) AddBlock(minedBlock *Block) {
	bc.Blocks = append(bc.Blocks, *minedBlock)
}

func (bc *Blockchain) getLastBlock() Block {
	lastBlock := bc.Blocks[len(bc.Blocks)-1]
	return lastBlock
}
