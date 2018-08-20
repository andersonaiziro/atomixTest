# Project Structure
Describe the project structure of the files included in the solution

## Main.go
### Entry point for the solution. Takes 1 Argument from the set: {"sc1", "sc2", "sc3", "sc4", "sc5"}
Main.go is the responsible for generating the scenarios to be run.
It sends transactions to be run by the client. There are 5 scenarios supported to validate the correct behavior.
1. sc1: intra-shard transaction; 
2. sc2: cross-shard transaction; 
3. sc3: cross-shard double spend trial; 
4. sc4: block creation every 10 transactions; 
5. sc5: 10,000 transactions with 30% chance of cross-shard 

## Client.go
### Implements the 3 Step Process described in Omniledger for Cross-Shard Transactions
Client.go communicates with all the shards in a multi-thread broadcasting of the 3-Step Protocol through Go Channels. 
Client sends the transactions to each of the shards, and receive Gossips as return values.
Client computes the aggregate consensus of the transaction and broadcasts the outcome to all shards through Go Channels. 

## Shard.go 
### Responsible for maintaining a mempool of UTXOs for all addresses present in that shard
Shards receives the transactions from the client and respond with Gossips. 
Client also maintains a list of commited transactions. Once that list reaches 10 transactions, shard will create a new block to be added to the blockchain.
Shards will use the Address field of the UTXO to denominate the owner of the UTXO. 
Once a UTXO is locked for a specific transaction, we'll modify that UTXO address to a unique key by getting the md5 hash of the transaction

If the client decides to commit the transaction, the shard will add the transaction to the list and remove all locked UTXOs from the mempool.
If the client decides to abort the transaction, the shard will revert the locked UTXOs to the previous address owner. 

## Transaction.go
### Generates the transactions to be used in this simulation
There are 3 types of transactions that can be generated:
1. Coinbase transaction (```NewCoinbaseTX```) to be used in the Genesis block
2. Distribution transactions (```NewDistributionTXs```) to populate initial shard states 
3. Cross-Shard / Intra-Shard transactions (```NewUTXOTransaction```) used to send funds between addresses

## Block.go
### Generates a new block for a set of transactions
There are 2 types of blocks generated: 
1. Genesis block (```NewGenesisBlock```) to be created when the blockchain is created
2. ```NewBlock``` is created by each shard's blockchain when the shard reaches 10 complete transactions 

## Gossips.go
### Types of responses that the shard can send to the client
According to the Omniledger paper, the gossips are the responses from the shards to the client. 
Also using this to compute the node agreement and decision to ```UnlockToCommit``` or ```UnlockToAbort```

## Blockchain.go
### Keeps a list of all blocks mined by each shard
```MineBlock``` is called by the shards when 10 transactions are completed

## UTXO.go
### Structure to maintain the unspent outputs in each shard mempool
This structure needs to know what was the origin transaction: ```TxID``` and ```OutputIndex```
This structure needs to know what was the destination address: ```Value```, ```Address```, and ```ShardID```
This is used by each shard to be able to quickly lookup what's the available balance for each address belonging to the shard. 
