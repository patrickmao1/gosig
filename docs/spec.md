# Gosig Implementation Specification

## Tech Stack

### Language

We use Golang. Go is by far among the most dominant language in blockchain implementations or even the broader distributed system landscape. It has an easy-to-use coroutine primitive for concurrency and it has many mature libraries we need such as bls12381 and goleveldb.

### Database

Blockchain nodes are often equiped with in-process key-value DBs. We will use **goleveldb** for it's a robust choice and proven usage by Ethereum and other major chains.

An alternative would be storing everything in memory. But in that case we would need to implement all the concurrency handlings. Using goleveldb which handles concurrent access for us makes it a lot simpler. Plus, using a persistent DB also helps debugging a lot as we can query the DB even after the program crashes.

### Toolings

Docker and Docker Compose are used for local testing and benchmarking before deploying to a real cluster.

## Messages/Data Types

### Block Related Messages

```protobuf
// sets the state of the blockchain at `key` to `val`
message SetOp {
    bytes key = 1;
    bytes val = 2;
}

// Transaction represents an ordered list of simple set(key, val) operation
message Transaction {
    repeated SetOp set_ops = 1;
}

message SignedTransaction {
    Transaction tx = 1;
    bytes sig = 2;
}

message TransactionHashes {
    repeated bytes tx_hashes = 1;
}

message SignedTransactions {
    repeated SignedTransaction txs = 1;
}

message BlockHeader {
    uint32 height = 1;
    bytes parent_hash = 2;
    Certificate proposal_cert = 3;
    // sign(privKey, Q^h), the preimage of the proposer score L^r
    bytes proposer_proof = 4;
    // a flat hash of all transaction hashes in this block. 
    // i.e. hash(hash(tx1) || hash(tx2) || ... || hash(txN))
    bytes tx_root = 5;
}
```

### _SignedMessage_

```protobuf
message SignedMessage {
    bytes sig = 1;
    int validator_index = 2;
    oneof message_types {
        BlockProposal proposal = 3;
        Prepare prepare = 4;
        TentativeCommit tc = 5;
    }
}
```

### _BlockProposal_

```protobuf
message BlockProposal {
    int round = 1;
    BlockHeader block_header = 2;
}
```

### _Prepare_

```protobuf
message Prepare {
    bytes block_hash = 1;
    uint32 block_height = 2;
}
```

### _TentativeCommit_

```protobuf
message TentativeCommit {
    bytes block_hash = 1;
    uint32 block_height = 2;
}
```

### _Certificate_

```protobuf
message Certificate {
    bytes agg_sig = 1; // aggregated signatures
    repeated uint32 sig_counts = 2; // an array of how many times each validator provided the signatures
    int round = 3; // the round in which this certificate is formed
}
```

## Networking

### Broadcast

Gosig relies on gossip to disseminate all types messages. The original implementation uses grpc-java but we think using connection based networking is unnecessary because of the repeated nature of gossiping. We implement a simple gossip network over UDP. The gossip degree is 3 (contacts 3 random nodes in the membership list in one loop) as specified in the paper.

### Membership

For simplicity, the member list is read from config files and stays sorted and fixed. i.e. no dynamic membership. Every node will know every other node from the start. Each member is identified by its public key.

### Top Level Interface

```go
type MessageHandler func(msg *proto.SignedMessage) error

interface Network {
    // Starts listening on the given port. `handler` is called every time a message arrives.
    Listen(port int, handler MessageHandler) error
    // Disseminates the message to the network. Procceeds in a preconfigured number of gossip rounds.
    Broadcast(msg *proto.SignedMessage) error
} 
```

## Crypto Primitives

### Signature Aggregation

We use BLS signature as specified in the paper. Signature aggregation is one of the main techniques of the paper (it's in the name). In Gosig, nodes only save two things for each message: the aggregated signature, and an integer array where each element at index i indicates how many times a validator of index i has signed the message. The reason why a validator can sign a message multiple times lies in how the signatures are aggregated. Here is the signature aggregation algorithm from a node's perspective:

I have a state for a message M: `s = (agg_sig, counts)` where `agg_sig` is the aggregated signature and `counts` is the aforementioned integer array.

When I receive a message M with `(incoming_agg_sig, incoming_counts)` from others:

1. `s.agg_sig += incoming_agg_sig`
2. For each `s.counts[i]`, `s.counts[i] += incoming_counts[i]`

Note that `len(counts) == len(incoming_counts)` because member list length is fixed.

To verify an aggregated signature:

1. Sum all `memberPubKeys[i] * s.counts[i]`, we get `aggPubKeys`
2. Check the bilinear mappings `e(g, s.agg_sig) ?= e(aggPubKeys, M)` where `g` is the generator.

Note that the paper does not specify a concrete curve from the BLS family. We will use BLS12-381 as it's has a well-tested library for Go.

### Hashing

The paper does not suggest a hash function. We will use SHA3-256 for it's pseudorandomness.

### Verifiable Random Function

VRF is used to randomly for selecting block proposers. We implement a VRF module as the paper specifies.

The initial `seed` is a random number that's hardcoded into every node. We use `initSeed = sha3("Gosig")`

```go
type VRF interface {
    Generate(privKey []byte, seed []byte) (rng, proof []byte)
    Verify(pubKey []byte, proposerProof []byte, signData []byte) bool
}
```

## The Blockchain Module

The blockchain module is an abstraction around the blockchain data structure.

```go
// TODO
type Blockchain interface {
    GetLatestCommitedHeader() *proto.BlockHeader
}
```

## System Parameters

- Round interval = 6s
- Proposal stage interval = 2s
- Agreement stage interval = 4s
- Gossip interval = 500ms
- Gossip degree = 3

## Tentative Commit State

Each node has a locally managed state for recording uncommitted but tentatively committed blocks as specified by the paper.

```protobuf
// internally used messages, for serialization convenience
message TcState {
  BlockHeader header = 1;
  uint32 round = 2;
  Certificate prepare_cert = 3;
}
```

TcState is persisted into the DB

```plaintext
"tc_state" -> TcState
```

## Round State

Each node maintains a state specific for the current round. The state is reset after going into a new round.

- `round`: the current round number, starting from 0 (the genesis round)
- `roundSeed`: fetched from the latest committed block = `round ++ block.ProposerProof`
- `phase`: the phase the node is in within the current round. Can be `Init`, `Prepared`, `TentativelyCommitted`, `Committed`
- `validProposals`: the valid proposals collected during this round. A "valid" `proposal` is defined as
  - `proposal.round == round`
  - `sha3(proposal.BlockHeader.ProposerProof)[:32] < T`
  - `vrf.Verify(pubKey, proposal.BlockHeader.ProposerProof, roundSeed) == true`
  - the TC-cert or P-cert is valid (signatures are valid)

## Consensus

### Random Seed Generation

Our implementation deviates from Gosig in what data the new seed depends on. Instead of having it depend on only the 
previous seed, we use `proposerProof` of block `h` as the seed for generating seed for height `h+1`. For proposing a new block $B^{h+1}$, $B^{h}$'s `proposerProof` has the form $SIG_{l^{h}}(r^{h},Q^{h-1})$ where $r^{h-1}$ is the round in which $B^{h-1}$ is proposed, which is just:

1. as random because it also depends on $Q^{h-1}$
2. as unpredictable because it's signed using a private key

Therefore, using `proposerProof` directly as $Q^{h-1}$ in the computation is just as secure and it removes redundancy.

Gosig does not specify when the random seed generated by the proposer of block `h` is revealed. If it's only revealed after a block is proposed, then if the proposer goes offline before revealing it, and the block eventually gets committed, the seed chain would break. Therefore, the only valid place to reveal the new seed without changing the protocol is in the block.

### Proposer Selection

1. Compute the `proposerScore` for the current `round`:
   - `seed = round ++ prevSeed` where `prevSeed` is just the `proposerProof` of the previous committeed block.
   - `proposerScore, proposerProof = vrf.Generate(privKey, seed)`
   - > Note that we must include the `round` in the seed so that `seed` is guaranteed to be different for every round. This is not discussed in the paper but it's likely designed to prevent the following attack: say we do not incorporate `r` into the computation of `seed`, i.e. `seed = prevSeed`. 1. a malicious proposer proposes two blocks and send to different honest parties. 2. No block is committed for this round because no block has 2f+1 votes. 3. The next round's `seed` will be the same because no new seed is revealed in the previous round. This means the proposer will be the same malicious validator. 4. Repeat these steps to compromise liveness.
2. Take the higher 32 bits of the score: `proposerScore[:32]`
3. The proposer threshold is defined as $T = 7/N * 2^{32}$ where N is the total number nodes (the authors of Gosig found 7/N is good enough)
4. If `proposerScore < T`, then I'm a proposer candidate

> Note: the last round's seed is effectively determined once I (a node) gather enough TC votes for a block (in other words, when I commit a block). At that time, the block can no longer be changed and so is the `proposerProof` in the block, which is used as the `prevSeed` for next round's seed computation.

### Proposing a Block

For each round, from a proposer's perspective:

1. If I have previously TC'd a block `state.TcBlock` but never got enough TC messages from others to commit it:
    1. Use `state.TcBlock` as the `newBlock`
    2. `newProposalCert = db.Get(sha3(state.TcBlock) + "/p_cert")`
2. Else:
    1. Pick some transactions `var txs []SignedTransaction` from the local pool
    2. Compute the root of `txs` (flat hash)
    3. Use the last committed block's TC-cert as the proposal cert `newProposalCert = db.Get(sha3(blockchain.GetLatestCommitted()) + "/tc_cert")` (This is probably to ensure that we are always building on a valid prev block. Note that the other branch doesn't need this because having a P-cert of a previously TC'd block implies that the proposer of that round did this step)
3. Build the `BlockHeader` with `height`, `txRoot`, `newProposalCert`, and `proposerProof`
4. Sign and gossip a message `&proto.Proposal{BlockHeader}`

When receiving a PR message, proceed with Algorithm 1 & 2 defined in the paper.

### Proposal Reduction

There can be multiple proposals proposed in a round. Correct nodes follow the following algorithm to decide on a single proposal to "prepare".

From a receiving node's perspective, at the end of the proposal phase:

1. If `len(validProposals) == 0`, then do nothing and return.
2. If there are some `validProposals`, then decide on which `proposal` to prepare:
   1. For all `validProposals`, find `candidateProposal` such that `candidateProposal.proposerScore` is the minimum among all.
   2. If my `state.TcBlock == nil`
      1. Reset `state` and decide on the block in `candidateProposal`.
   3. Else
      1. If `candidateProposal.BlockHeader.ProposalCert.Round > state.TcRound`:
         1. Reset `state` and decide on the block in `candidateProposal`.
      2. Else if there is any `validProposals[i].BlockHeader == state.TcBlock` and `validProposals[i].BlockHeader.ProposalCert.Round >= state.TcBlock`
         1. `state.TcRound = validProposals[i].BlockHeader.ProposalCert.Round`
         2. `state.TcCert = validProposals[i].BlockHeader.ProposalCert`
      3. Else: decide on no block for this round (skip the round and do not react to any message)
   4. If decided on some block, follow Algorithm 1.

<!-- ### Voting for Prepare and TentativeCommit Messages

#### Loop 1: Handling Prepare Messages

For each received prepare message

1. If the `phase != Prepared`, ignore the message
2. Find `Prepare` -->

### Round Synchronization

The concensus process proceeds in rounds. A round consists of two stages: the proposal phase and the agreement phase. `roundDuration` is implicitly defined as `proposalStageDuration + agreementStageDuration`. The current round is defined as `round = (now - genesisTime) / roundDuration`. `nextRoundTime = genesisTime + (round + 1) * roundDuration`. Corrects node should behave according to the following rules:

- Only perform round transition according to the node's local clock.
- When round transition happens (round -> round + 1), the node should immediately abandon everything it's doing in the previous round and move to the new round.
- When stage transition happens within a round (proposal stage -> agreement stage), the node should also immediately move to the next stage. This means it will derive the chosen proposal from the list of valid proposals it has at the transition time.

Note that because the clocks are not perfectly synchronized, it is very likely that a node can receive Prepare messages while it thinks it's still in the proposal stage. The node doesn't need to immediately transition to agreement stage when this happens. This is ok because missing several gossip messages of Prepare doesn't matter at much. Prepare (or TC) messages have aggregated signatures, so messages later in the round will carry the effects of the earlier messages.

### Block Storage

Block data is persisted into disk. goleveldb is used for this task. Storage layout is defined here:

```plaintext
"block/{blockHash}/header" -> block
"block/{blockHash}/transactions" -> transaction hashes
"block/{blockHash}/p_cert" -> prepare certificate
"block/{blockHash}/tc_cert" -> tc certificate
"block/head" -> blockHash
"transaction/{txHash}" -> transcation
```

### Catching Up on Missing Blocks

From the consensus process it's easy to see that if a node does not have the latest committed block, then it cannot verify anything because it needs to know the latest seed to verify proposals.

A node can discover that it doesn't have the latest committed block once it receives a `BlockProposal` in which the block height is greater than the node's local height + 1. In this case, the lagging node needs to catch up by querying peers for blocks `[height+1, proposal.block.height)`

Peers to query are chosen randomly.

1. Query X random peers asking for the missing range
2. Peers respond with the range they have
3. Download from the peer with the longest range
4. Update the missing range
5. Repeat from step 1 until head is updated to (proposal.block.height - 1)

## Transactions & Accounts

For transactions, we want to keep it as simple as possible because it's not in the scope of a consensus protocol. We implement a simple payment system where users can send money to each other. We do not implement a global transaction mempool. Instead, each node will have its own pool. For testing convenience, we implement a client that sends transactions to all nodes all at once so that it is guaranteed to be picked up by a random proposer in the next round.

A user of the blockchain is identified by their public key and so is their account on the chain. An account can only hold one type of state is their balance. A transaction is valid if:

1. The sender has enough money
2. The transaction is correctly signed

Each user's account state is saved in the DB as

```plaintext
"{pubKey}" -> balance
```

The node should implement a GRPC server to receive transactions. For testing and benchmark purposes, there are also RPCs for querying transaction states and statistics.

```protobuf
service NodeServer {
    rpc SubmitTransaction(SubmitTransactionReq) returns (SubmitTransactionRes);
    rpc TransactionStatus(TransactionStatusReq) returns (TransactionStatusRes);
}

message SubmitTransactionReq {
    SignedTransaction tx = 1;
}

message SubmitTransactionRes {
    optional string err = 1;
}

message TransactionStatusReq {
    bytes tx_hash = 1;
}

message TransactionStatusRes {
    optional string err = 1;
    string status = 2;
}
```

## Pipelining

### Background Block Data Dissemination

Gosig separates the dissenmination of block headers and block data so that the latency-bound block proposal phase can be shortened and nodes can start sending Prepare while

## Benchmarking

TODO

## Optimizations Not Implemented

There are many optimizations that are left out in this implementation to make the scope approachable. These are listed here.

- Block body segmentation: splitting transactions into batches for fine-grained pipelining.
