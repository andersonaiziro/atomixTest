# atomixTest
Testing Blockchain Atomix 3-Step Process

## Scenario 1: Intra-Shard Transaction
### Run: atomixTest.exe sc1
Sample Initial State:
--- Shard 1
UTXOs for Address: a1
TxID: 968f150bf5d0d07d255456811b5070fbefd7389865e13356ee6bbd259d9933d6, OutIndex: 0, Address: a1, Value: 5
TxID: 968f150bf5d0d07d255456811b5070fbefd7389865e13356ee6bbd259d9933d6, OutIndex: 2, Address: a1, Value: 5
UTXOs for Address: b1
TxID: 968f150bf5d0d07d255456811b5070fbefd7389865e13356ee6bbd259d9933d6, OutIndex: 1, Address: b1, Value: 5
TxID: 968f150bf5d0d07d255456811b5070fbefd7389865e13356ee6bbd259d9933d6, OutIndex: 3, Address: b1, Value: 5

Sample Final State: 
Executed IntraShard Transaction on Shard 1. From node a1 to node b1
--- Shard 1:
UTXOs for Address: a1
TxID: 968f150bf5d0d07d255456811b5070fbefd7389865e13356ee6bbd259d9933d6, OutIndex: 2, Address: a1, Value: 5
UTXOs for Address: b1
TxID: 968f150bf5d0d07d255456811b5070fbefd7389865e13356ee6bbd259d9933d6, OutIndex: 1, Address: b1, Value: 5
TxID: 968f150bf5d0d07d255456811b5070fbefd7389865e13356ee6bbd259d9933d6, OutIndex: 3, Address: b1, Value: 5
TxID: **327c48067a9ce9839b59e38ac2baf9eb4cdb52e5316b379519989ad230ad6cac**, OutIndex: 0, Address: b1, Value: 5

## Scenario 2: Cross-Shard Transaction
### Run: atomixTest.exe sc2
Sample Initial State:
--- Shard 4:
UTXOs for Address: b4
TxID: c24323642a6eca815ae6e489b87fd84897de5d174ab7270d9ad888ce346f3173, OutIndex: 1, Address: b4, Value: 5
TxID: c24323642a6eca815ae6e489b87fd84897de5d174ab7270d9ad888ce346f3173, OutIndex: 3, Address: b4, Value: 5
UTXOs for Address: a4
TxID: c24323642a6eca815ae6e489b87fd84897de5d174ab7270d9ad888ce346f3173, OutIndex: 0, Address: a4, Value: 5
TxID: c24323642a6eca815ae6e489b87fd84897de5d174ab7270d9ad888ce346f3173, OutIndex: 2, Address: a4, Value: 5
--- Shard 6:
UTXOs for Address: a6
TxID: 2f7d07799be47c1fe072cddf399103fddda326568b59dc57a3e1b0d3029ada03, OutIndex: 0, Address: a6, Value: 5
TxID: 2f7d07799be47c1fe072cddf399103fddda326568b59dc57a3e1b0d3029ada03, OutIndex: 2, Address: a6, Value: 5
UTXOs for Address: b6
TxID: 2f7d07799be47c1fe072cddf399103fddda326568b59dc57a3e1b0d3029ada03, OutIndex: 1, Address: b6, Value: 5
TxID: 2f7d07799be47c1fe072cddf399103fddda326568b59dc57a3e1b0d3029ada03, OutIndex: 3, Address: b6, Value: 5

Sample Final State: 
Executed Cross-Shard Transaction from Shard 4 to Shard 6.
Executed Cross-Shard Transaction from Node b4 to Node a6.
--- Shard 4:
UTXOs for Address: a4
TxID: c24323642a6eca815ae6e489b87fd84897de5d174ab7270d9ad888ce346f3173, OutIndex: 0, Address: a4, Value: 5
TxID: c24323642a6eca815ae6e489b87fd84897de5d174ab7270d9ad888ce346f3173, OutIndex: 2, Address: a4, Value: 5
UTXOs for Address: b4
TxID: c24323642a6eca815ae6e489b87fd84897de5d174ab7270d9ad888ce346f3173, OutIndex: 3, Address: b4, Value: 5
--- Shard 6:
UTXOs for Address: a6
TxID: 2f7d07799be47c1fe072cddf399103fddda326568b59dc57a3e1b0d3029ada03, OutIndex: 0, Address: a6, Value: 5
TxID: 2f7d07799be47c1fe072cddf399103fddda326568b59dc57a3e1b0d3029ada03, OutIndex: 2, Address: a6, Value: 5
**TxID: c25955ec55f0bcf97bee6624336f8edb10986b5f522811c28f4b2e69eebfd07d, OutIndex: 0, Address: a6, Value: 5**
UTXOs for Address: b6
TxID: 2f7d07799be47c1fe072cddf399103fddda326568b59dc57a3e1b0d3029ada03, OutIndex: 1, Address: b6, Value: 5
TxID: 2f7d07799be47c1fe072cddf399103fddda326568b59dc57a3e1b0d3029ada03, OutIndex: 3, Address: b6, Value: 5

## Scenario 3: 2 Cross-Shard Transactions. 2nd One tries to double-spend 
### Run: atomixTest.exe sc3
In this test, there should only be 1 new UTXO added to all shards.

## Scenario 4: 30 Intra-Shard Transactions. Validate Block Creation every 10 transactions
### Run: atomixTest.exe sc4
In this test, the blockchain for one of the shards should have 3 blocks added to the chain. 
Example Output: 
Block 0, Hash: 3e51664dfe7026bcca67f47b312fd5d6f6c511ce19cff68d7a8f1e3ec2667373. TransactionCount: 1.
Transactions: db444f76b64748ebabc315bc53d2fc31eeb515e526cde58e26cc1271c01cf29b
Block 1, Hash: afb058ae4d1873ffce225b4012654f8c0e89d83a97df363667c89b0d2b53903e. TransactionCount: 10.
Transactions: d1754603526d5ec2ac9981405688dcaa81953a18f59ccbafb06472a179ce2233, 851973882be04abb0f1f83f47fd4b7c58fd66a7488932b9492c812c797029df3, ee912a0f1119bdcce9924da4aa705b0932f25491ad6120fe9f79bc22c786ab13, d10f688280fa94673df3c659725a6d0828eb3354db0b200777ad5edb8edfaaaf, a1ddddb1081014b2c935545c632f41d4de29b6a05262f3e4d444aee540e19176, 35ac2944e31b9933837ff29eef3e08a617c310e0e69ea246895c8dc04c4f2eae, e3fa1a0b3862d069f3da3f14a336dc93f2abb84c177135858d5a1063d797d5cb, 9828e875101db89832aa4b79ffe7b63418396096d9835c9da294464bce21843c, 88b1721df3362f187e3cebf2cd21ac20d22235eeac1ef47f21ac3c94ab373828, 4897d0cc859f49650630d86c7f7995451d6e201822c764ef80726a60536f017a
Block 2, Hash: 51a76544bc8e389cca20b3ea366c1f41c11e5470a57448850aa109f311cdcd16. TransactionCount: 10.
Transactions: 7f4ddecb665e04cb1e345b02b3df295f0e02c66e3f4201904c445024535fd99b, 64d9c6004274558f68fe9978dc3001a7a2f458dc4df598323cef5f9716cb1b1e, 871f38c4c50a5083fa69162ecf9fa39f26a724fffa81e06f1f4401e2081ee940, cd57475b8280b59e13f37e93f45fe56ae57c2c7a208bee8fa7acf6cb39cab425, 0d9d12372a07e5712a2762095d84a3834f741080b03351c10c04cdb50037967d, 5f5bec87e2125878544f2dcdb71654f3d180ab2a03ba0106445bd380914d814a, 56a4b83116c6218ae81db9f3c78eea382c68665682cb53a36a480457ce2d598f, 5e168d85b0a7ac50e0671f5d2f663f90692e621d6f208cb0267095f11e6d7c01, a46745290c269b0f148127e28028dc68cdf568132e19e244e99d28d14bd9f34e, d9f039883b0f2966b94cb2d9cf6822d6ea921ec6a7af9d2049540ea15c3571c1
Block 3, Hash: 3536bf6f0c80953ef03f79d6d9af0916f663f4aafa7c727c5ee4531ad399ea3f. TransactionCount: 10.
Transactions: 446a6418d7907cd88da8b471d0c944f36de34e643733b090095cf123c3fd1778, eb3c3abfc0e4d9b661a21e0374c2730597bcb6e3c5ceed58af18fa9d097244bb, 4d5112bb990a592236ce4d2360e871344591dd529ede3664b49d2bb7904a891c, 81dfa6d566bb99d5d746c9e9e5d95736216be7289b83598c05f0f1e8effe09c6, 0b7a6da616475b35a61c1d666d8c2351c94c4b90538c47237befd533bad5f11e, d3b7d20c7a4ebf8c8b78df08ba087840860c6e5fa72557470272fa67eec2d52c, b281c1595da981de198034a1cd26c380e2d56315b725c88cdbfc43b296318ff1, 5c7ae9fc9fc0924eca17bed804718185713ba68312872804550827a5961b1711, 5ecdff0d16b27abeae7befcc4ec0ef81b14a464686db922f9ec4c9ee72c35d33, eb234c5c8dbc459bcf76a744f0321a29aa46035e43b99a08ea2bbe56d2ea319e

## Scenario 5: 10,000 Transactions. 30% of Cross-Shard Transaction. 
### Run: atomixTest.exe sc5
In this test, we'll execute 10,000 transactions. 30% of them are cross-sharded. 
In order to validate the correctness of the output, on average each shard should have between 100 and 130 blocks. 
