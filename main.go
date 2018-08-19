package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "sc1":
		runSimpleIntraShardTransaction()
	case "sc2":
		runCrossShardTransaction()
	case "sc3":
		runDoubleSpendScenario()
	case "sc4":
		runCreateNewBlock()
	case "sc5":
		run10kTransactions()
	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage: ")
	fmt.Println("	AtomixTest.exe [scenario]")
	fmt.Println("The supported scenarios are: ")
	fmt.Println("   sc1		Simple Intra-Shard Transaction")
	fmt.Println("   sc2		Simple Cross-Shard Transaction 3 step process")
	fmt.Println("   sc3		2 Cross-Shard Transactions trying to double spend a Cross-Shard Transaction")
	fmt.Println("   sc4		10 Intra-Shard Transactions. Validate new Block Creation")
	fmt.Println("   sc5		10,000 random transactions. 70% chance Intra-Shard")
}

func runSimpleIntraShardTransaction() {
	client := Client{sampleShards, 2, 5, []*Shard{}, []*chan *Transaction{}, []*chan Gossips{}}
	client.Initialize()

	//Picks a random index... operate on the shard at that index
	chosenShard1, _ := pickRandom2ShardIndexes()
	shardID, shardAddresses := client.getShardNodeAddresses(chosenShard1)

	tx := client.buildTransaction(uint32(chosenShard1), shardAddresses[0], shardID, shardAddresses[1])
	client.ExecuteNewTransaction(tx, false)
	log.Println(fmt.Sprintf("Executed IntraShard Transaction on Shard %d. From node %s to node %s", shardID, shardAddresses[0], shardAddresses[1]))
	client.printShardUTXOStates(shardID)
}

func runCrossShardTransaction() {
	client := Client{sampleShards, 2, 5, []*Shard{}, []*chan *Transaction{}, []*chan Gossips{}}
	client.Initialize()

	//Picks 2 random indexes... operate on the shard at that index
	chosenShard1, chosenShard2 := pickRandom2ShardIndexes()
	shard1ID, shard1Addresses := client.getShardNodeAddresses(chosenShard1)
	shard2ID, shard2Addresses := client.getShardNodeAddresses(chosenShard2)

	tx := client.buildTransaction(uint32(chosenShard1), shard1Addresses[0], shard2ID, shard2Addresses[0])
	client.ExecuteNewTransaction(tx, false)
	log.Println(fmt.Sprintf("Executed Cross-Shard Transaction from Shard %d to Shard %d.", shard1ID, shard2ID))
	log.Println(fmt.Sprintf("Executed Cross-Shard Transaction from Node %s to Node %s.", shard1Addresses[0], shard2Addresses[0]))
	client.printShardUTXOStates(shard1ID)
	client.printShardUTXOStates(shard2ID)
}

func runDoubleSpendScenario() {
	client := Client{sampleShards, 2, 5, []*Shard{}, []*chan *Transaction{}, []*chan Gossips{}}
	client.Initialize()

	//Picks 2 random indexes... operate on the shard at that index
	chosenShard1, chosenShard2 := pickRandom2ShardIndexes()
	shard1ID, shard1Addresses := client.getShardNodeAddresses(chosenShard1)
	shard2ID, shard2Addresses := client.getShardNodeAddresses(chosenShard2)

	tx := client.buildTransaction(uint32(chosenShard1), shard1Addresses[0], shard2ID, shard2Addresses[0])
	client.ExecuteNewTransaction(tx, false)
	log.Println(fmt.Sprintf("Executed Cross-Shard Transaction from Shard %d to Shard %d.", shard1ID, shard2ID))
	log.Println(fmt.Sprintf("Executed Cross-Shard Transaction from Node %s to Node %s.", shard1Addresses[0], shard2Addresses[0]))
	client.printShardUTXOStates(shard1ID)
	client.printShardUTXOStates(shard2ID)

	log.Println(fmt.Sprintf("Try double spending. Run same UTXO from Node %s to Node %s.", shard1Addresses[0], shard2Addresses[0]))
	client.ExecuteNewTransaction(tx, false)
	client.printShardUTXOStates(shard1ID)
	client.printShardUTXOStates(shard2ID)
}

func runCreateNewBlock() {
	client := Client{sampleShards, 2, 5, []*Shard{}, []*chan *Transaction{}, []*chan Gossips{}}
	client.Initialize()
	//Picks a random index... operate on the shard at that index
	chosenShard1, _ := pickRandom2ShardIndexes()
	shardID, shardAddresses := client.getShardNodeAddresses(chosenShard1)

	//run 30 transactions should create 3 new blocks
	for i := 0; i < 30; i++ {
		tx := client.buildTransaction(uint32(chosenShard1), shardAddresses[i%2], shardID, shardAddresses[(i+1)%2])
		log.Println(fmt.Sprintf("Transaction %x on Shard %d. From node %s to node %s", tx.ID, shardID, shardAddresses[i%2], shardAddresses[(i+1)%2]))
		client.ExecuteNewTransaction(tx, false)
	}

	client.printShardUTXOStates(shardID)
	client.printShardBlockchainStates(shardID)
}

func run10kTransactions() {
	client := Client{sampleShards, 5, 5, []*Shard{}, []*chan *Transaction{}, []*chan Gossips{}}
	client.Initialize()
	for i := 0; i < 10000; i++ {
		//Picks 2 random indexes and 2 random addresses... operate on the shard at that index
		chosenShard1, chosenShard2 := pickShardsToUseIn10kTransactions()
		_, shard1Addresses := client.getShardNodeAddresses(chosenShard1)
		shard2ID, shard2Addresses := client.getShardNodeAddresses(chosenShard2)

		fromAddressIndex := rand.Intn(1)
		toAddressIndex := rand.Intn(1)

		tx := client.buildTransaction(uint32(chosenShard1), shard1Addresses[fromAddressIndex], shard2ID, shard2Addresses[toAddressIndex])

		if tx == nil {
			//	log.Println(fmt.Sprintf("Tried to use a Node %s in Shard %d without UTXOs", shard1Addresses[fromAddressIndex], chosenShard1))
			continue
		}

		client.ExecuteNewTransaction(tx, false)
	}

	for j := 1; j <= 10; j++ {
		client.printShardUTXOStates(uint32(j))
	}
	for j := 1; j <= 10; j++ {
		client.printShardBlockchainStates(uint32(j))
	}
}

func pickShardsToUseIn10kTransactions() (int, int) {
	shard1 := rand.Intn(10)
	shard2 := shard1

	if rand.Intn(10) < 3 {
		shard2 = rand.Intn(10)
	}

	return shard1, shard2
}

func pickRandom2ShardIndexes() (int, int) {
	shard1 := rand.Intn(10)
	shard2 := rand.Intn(10)

	for shard1 == shard2 {
		shard2 = rand.Intn(10)
	}

	return shard1, shard2
}

func (client *Client) getShardNodeAddresses(shardIndex int) (uint32, []string) {
	var shardAddresses []string
	for address := range client.shards[shardIndex].Mempool {
		shardAddresses = append(shardAddresses, address)
	}

	return client.shards[shardIndex].ShardID, shardAddresses
}
func (client *Client) buildTransaction(shardIndex uint32, fromAddress string, toShardID uint32, toAddress string) *Transaction {
	utxos := client.shards[shardIndex].Mempool[fromAddress]
	if len(utxos) == 0 {
		return nil
	}
	utxo := utxos[0]
	return NewUTXOTransaction(utxo, toAddress, toShardID, 5)
}

// shards to be used by all test scenarios
var sampleShards = map[uint32][]string{
	1:  {"a1", "b1"},
	2:  {"a2", "b2"},
	3:  {"a3", "b3"},
	4:  {"a4", "b4"},
	5:  {"a5", "b5"},
	6:  {"a6", "b6"},
	7:  {"a7", "b7"},
	8:  {"a8", "b8"},
	9:  {"a9", "b9"},
	10: {"a10", "b10"},
}
