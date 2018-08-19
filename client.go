package main

import (
	"log"
)

// Client responsible for creating blockchain and executing the transactions
type Client struct {
	shardAddressMap   map[uint32][]string
	initialTXOutCount int
	initialTXOutValue int
	shards            []*Shard
	channelsIn        []*chan *Transaction
	channelsOut       []*chan Gossips
}

//Initialize is responsible for setup of the blockchain
func (client *Client) Initialize() {
	client.initializeBlockchain()
}

//ExecuteNewTransaction runs the 3-Step Omniledger Transaction
func (client *Client) ExecuteNewTransaction(transaction *Transaction, printShardStates bool) {
	//Step 1: Initialize and Broadcast transaction
	client.initializeAtomixTransaction()
	client.broadcastAtomixTransaction(transaction)

	//Step 2: Process Shard Proof Of Acceptance or Proof of Rejection
	gossipResult := client.processGossips()

	//Step 3: Unlock Shard Transactions
	client.broadcastTransactionResult(gossipResult, transaction)
}

// Step 1a: Prepare all shards to receive new transaction
func (client *Client) initializeAtomixTransaction() {
	for index, shard := range client.shards {
		inChannel := client.channelsIn[index]
		outChannel := client.channelsOut[index]
		go shard.ProcessNewTransaction(*inChannel, *outChannel)
	}
}

// Step 1b: Broadcast transaction to all shards
func (client *Client) broadcastAtomixTransaction(transaction *Transaction) {
	for _, channelIn := range client.channelsIn {
		*channelIn <- transaction
	}
}

// Step 2: Process Gossips from Channels
func (client *Client) processGossips() Gossips {
	gossipResult := UnlockToCommit
	for _, channelOut := range client.channelsOut {
		channelGossip := <-*channelOut

		switch channelGossip {
		case ProofOfRejection:
			gossipResult = UnlockToAbort
		default:
		}
	}

	return gossipResult
}

// Step 3: Unlock: unlock to abort or unlock to commit
func (client *Client) broadcastTransactionResult(gossipResult Gossips, transaction *Transaction) {
	for index, shard := range client.shards {
		inChannel := client.channelsIn[index]
		outChannel := client.channelsOut[index]
		switch gossipResult {
		case UnlockToCommit:
			go shard.CommitTransaction(*inChannel, *outChannel)
		case UnlockToAbort:
			go shard.AbortTransaction(*inChannel, *outChannel)
		}

		*inChannel <- transaction
	}

	for _, channelOut := range client.channelsOut {
		<-*channelOut
	}
}

func (client *Client) printShardUTXOStates(shardID uint32) {
	for _, shard := range client.shards {
		if (shard.ShardID == shardID) || (shardID == 0) {
			log.Println(shard.GetMempoolString())
		}
	}
}

func (client *Client) printShardBlockchainStates(shardID uint32) {
	for _, shard := range client.shards {
		if (shard.ShardID == shardID) || (shardID == 0) {
			log.Println(shard.GetBlockchainString())
		}
	}
}

func (client *Client) initializeBlockchain() {
	//Create blockchain and setup shards with nodes and initial token balances
	client.channelsIn = []*chan *Transaction{}
	client.channelsOut = []*chan Gossips{}

	for shardID, addresses := range client.shardAddressMap {
		bc := NewBlockchain()
		shard := NewShard(shardID, bc)

		initialTokenDistribution := NewDistributionTXs(addresses, shardID, client.initialTXOutValue, client.initialTXOutCount)
		in := make(chan *Transaction)
		out := make(chan Gossips)

		//log.Println(initialTokenDistribution.String())
		go shard.ProcessNewTransaction(in, out)

		in <- initialTokenDistribution
		<-out
		//log.Println(result.String())
		//log.Println(shard.GetMempoolString())
		client.shards = append(client.shards, shard)
		client.channelsIn = append(client.channelsIn, &in)
		client.channelsOut = append(client.channelsOut, &out)
	}
}
