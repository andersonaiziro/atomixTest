package main

// Gossips contains the messages switched between the client and the shards/nodes
type Gossips int

// ProofOfAcceptance: used when there is enough funds to process the transaction, or when the transaction doesn't involve any node in that shard
// ProofOfRejection: used when there is not enough fund to process the transaction
// UnlockToCommit: used when transactions should be commited
// UnlockToAbort: used when transactions should be reverted
const (
	ProofOfAcceptance   Gossips = iota // 0
	ProofOfRejection                   // 1
	UnlockToCommit                     // 2
	UnlockToAbort                      // 3
	TransactionCommited                // 4
)

func (gossip Gossips) String() string {
	switch gossip {
	case ProofOfAcceptance:
		return "ProofOfAcceptance"
	case ProofOfRejection:
		return "ProofOfRejection"
	case UnlockToCommit:
		return "UnlockToCommit"
	case UnlockToAbort:
		return "UnlockToAbort"
	case TransactionCommited:
		return "TransactionCommited"
	default:
		return "Unknown"
	}
}
