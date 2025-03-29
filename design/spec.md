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
    Certificate proposal_cert = 2;
    // sign(privKey, Q^h), the preimage of the proposer score L^r
    bytes proposer_proof = 3;
    // a flat hash of all transaction hashes in this block. 
    // i.e. hash(hash(tx1) || hash(tx2) || ... || hash(txN))
    bytes tx_root = 4;
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
    BlockHeader block_header = 1;
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
    bytes signers = 2; // a bitmap of who provided the signatures
}
```

## Networking

### Broadcast

Because Gosig uses gossip, we use UDP to implement a very basic gossip protocol. The gossip degree $k$ is adjustable. Default k = 3 (contacts 3 random nodes in the membership list in one loop) as specified in the paper.

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

The paper does not suggest a hash function. We will use SHA3-256 as it's the latest gen, fast and secure.

### Verifiable Random Function

VRF is used to randomly select block proposers. We implement a VRF module as the paper specifies.

The initial `seed` is a random number that's hardcoded into every node. We use `initSeed = sha3("Gosig")`

```go
type VRF interface {
    Generate(privKey []byte, seed []byte) (rng, proof []byte)
    Verify(pubKey []byte, num []byte, seed []byte) bool
}
```

## The Blockchain Module

The blockchain module is an abstraction around the blockchain data structure.

```go
// TODO
type Blockchain interface {
    GetLatestCommited() Block
}
```

## Local State

Each node has a local state as specified by the paper.

```go
type LocalState struct {
    TcBlock *proto.BlockHeader
    TcRound int
    TcCert  *proto.Certificate
}

var state LocalState
```

## System Parameters

- Round interval = 5s
- Proposal stage interval = 2s
- Agreement stage interval = 3s
- Gossip interval = 300ms
- Gossip degree = 3

## Block Proposal

### Proposer Selection

After the proposer of the last round reveal the seed of the previous block `prevSeed`

1. Compute the `proposerScore` for the current `round`:
   1. `seed = r ++ roundSeed`
   2. `proposerScore, proposerProof = vrf.Generate(privKey, seed)`
2. Take the higher 32 bits of the score: `proposerScore[:32]`
3. The proposer threshold is defined as $T = 7/N * 2^{32}$ where N is the total number nodes (the authors of Gosig found 7/N is good enough)
4. If `proposerScore < T`, then I'm a proposer

> Note: the last round's seed reveal is effectively determined once I (a node) gather enough TC votes for a block (in other words, when I commit a block). At that time, the block can no longer be changed and so is the `seed_reveal` in the block.

### Proposing a Block

For each round, from a proposer's perspective:

1. If I have previously TC'd a block `state.TcBlock` but never got enough TC messages from others to commit it:
    1. Use `state.TcBlock` as the `newBlock`
    2. `newProposalCert = db.Get(sha3(state.TcBlock) + "/p_cert")`
    3. (Note: This check ensures that if prepared(B), then it will eventually be committed)
2. Else:
    1. Pick some transactions `var txs []SignedTransaction` from the local pool
    2. Compute the root of `txs` (flat hash)
    3. Use the last committed block's TC-cert as the proposal cert `newProposalCert = db.Get(sha3(blockchain.GetLatestCommitted()) + "/tc_cert")` (This is probably to ensure that we are always building on a valid prev block. Note that the other branch doesn't need this because having a P-cert of a previously TC'd block implies that the proposer of that round did this step)
3. Build the `BlockHeader` with `height`, `txRoot`, `newProposalCert`, and `proposerProof`
4. Sign and gossip a message `&proto.Proposal{BlockHeader}`

When receiving a PR message, proceed with Algorithm 1 & 2 defined in the paper.

> ### Unresolved Question
>
> When do I reveal the seed for the next round?
> If I only reveal it after a block is committed, then if I go offline and don't reveal, the seed chain would break.
> If I reveal it directly in block, then why can't we just use the proposer proof for next seed generation?
> Maybe look more into Algorand's method for this since Gosig's VRF is taken from Algorand.

## Transactions

For transactions, we want to keep it as simple as possible because it's not in the scope of a consensus protocol. We do not implement a global transaction mempool. Instead, each node will have its own pool. For testing convenience, we implement a client that sends transactions to all nodes all at once so that it is guaranteed to be picked up by a random proposer in the next round.

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

### Background Block Dissemination

## Benchmarking

TODO

## Optimizations Not Implemented

There are many optimizations that are left out in this implementation to make the scope approachable. These are listed here.

- Block body segmentation: splitting transactions into batches for fine-grained pipelining.
-
