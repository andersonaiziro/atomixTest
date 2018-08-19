package main

type UTXO struct {
	TxID        [32]byte
	OutputIndex int
	Value       int    // The amount of money
	Address     string // The account address
	ShardID     uint32 // The Id of the shard where this UTXO belongs
}

// NewUTXO creates a new UTXO to be added to a mempool
func NewUTXO(transaction Transaction, outIndex int) (string, *UTXO) {
	txout := transaction.TxOutput[outIndex]
	utxo := UTXO{TxID: transaction.ID,
		OutputIndex: outIndex,
		Value:       txout.Value,
		Address:     txout.Address,
		ShardID:     txout.ShardID}

	return txout.Address, &utxo
}
