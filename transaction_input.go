package main

type TXInput struct {
	TxID        [32]byte
	OutputIndex int    // The index of the TXOutput this Input is referring to
	Address     string // The account address
	ShardID     uint32 // The shard this block belongs to
}
