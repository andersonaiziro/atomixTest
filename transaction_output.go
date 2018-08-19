package main

type TXOutput struct {
	Value   int    // The amount of money
	Address string // The account address
	ShardID uint32 // The Id of the shard where this UTXO belongs
}

// TXOutputs collects TXOutput
type TXOutputs struct {
	Outputs []TXOutput
}

// NewTXOutput create a new TXOutput
func NewTXOutput(value int, address string, shardId uint32) *TXOutput {
	txo := &TXOutput{value, address, shardId}

	return txo
}
