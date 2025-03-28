# Gosig Implementation Specification

## Messages

### Block Related Messages

```protobuf
// sets state of the blockchain at `key` to `val`
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
    // sign(privKey, r || Q^h), the preimage of the proposer score L^r
    bytes proposer_proof = 3;
    // a flat hash of all transaction hashes in this block. 
    // i.e. hash(hash(tx1) || hash(tx2) || ... || hash(txN))
    bytes tx_root = 4;
    bytes seed_reveal = 5;
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

Because Gosig uses gossip, we use UDP to implement a very basic gossip protocol. The gossip degree $k$ is adjustable. Default k = 3 (contacts 3 random nodes in the membership list in one loop).

### Membership

For simplicity, the member list is read from config files and stays sorted and fixed. i.e. no dynamic membership. Each member is identified by its public key.

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

We use BLS12-381 as specified in the paper. Signature aggregation is one of the main contribution of the paper. In Gosig, nodes only save two things for each message: the aggregated signature, and an integer array that indicates how many times each validator has signed the message. The reason why a validator can sign a message multiple times lies in how the signatures are aggregated. Here is the signature aggregation algorithm:

I have a state for a message M: `s = (agg_sig, counts)` where `counts` is the aforementioned integer array.

When I receive a message M with `incoming_agg_sig, incoming_counts` from others:

1. `s.agg_sig += incoming_agg_sig`
2. For each `s.counts[i]`, `s.counts[i] += incoming_counts[i]`

Note that `len(counts) == len(incoming_counts)` because member list length is fixed.

To verify an aggregated signature:

1. Sum all `memberPubKeys[i] * s.counts[i]`, we get `aggPubKeys`
2. Check `e(g, s.agg_sig) ?= e(aggPubKeys, M)` where `g` is the generator.

Every time a node gossips with others (e.g. to disseminate the Prepare message), it always send these two things.

### Hashing

The paper does not suggest a hash function. We will use SHA3-256 as it's the latest gen, fast and secure.

### Verifiable Random Function

VRF is used to randomly select block proposers. We implement a VRF module as the paper specifies.

The initial seed $Q^0$ is a random number that's hardcoded into every node. We use $Q^0 = hash("Gosig")$

```go
interface VRF {

}
```

## The Blockchain Module

The blockchain module is an abstraction around the blockchain data structure.

```go
// TODO
```

## Local State

Each node has a local state as specified by the paper.

```go
type LocalState struct {
    tc_block *proto.BlockHeader
    tc_round int
    tc_cert  *proto.Certificate
}
```

## Block Proposal

### Proposer Selection

After the proposer of the last round reveal the seed of the previous block $Q^{h-1}$

1. Compute The Proposer Score $L^r$ for a round $r$: $L^r = SIG(r, Q^h)$
2. Take the higher 32 bits of $L^r$
3. The proposer threshold is defined as $T = 7/N * 2^{32}$ where N is the total number nodes (the authors of Gosig found 7/N is good enough)
4. If $L^r < T$, then I'm a proposer

> Note: the last round's seed reveal is effectively determined once I (a node) gather enough TC votes for a block (in other words, when I commit a block). At that time, the block can no longer be changed and so is the `seed_reveal` in the block.

### Proposing a Block

For each round, from a proposer's perspective:

1. If I have previously TC'd a block $B'$ but never got enough TC messages from others to commit it:
    1. Use this previously TC'd block as the new block $B=B'$
    2. Use the previously TC'd block's P-certificate as the new proposal certificate $c$
    3. (Note: This check ensures that if prepared(B), then it will eventually be committed)
2. Else:
    1. Pack some transactions $txs$ into a new block $B$
    2. Use the last committed block's TC-cert as the proposal cert $c$ (This is probably to ensure that we are always building on a valid prev block. Note that the other branch doesn't need this because having a P-cert of a previously TC'd block implies that the proposer of that round did this step)
3. Set $pp=SIG(r,Q^h)$ which is the pre-image of my $L^r$. Other nodes can check $H(pp) < T$ to ensure that I am indeed a valid proposer.
4. Sign and gossip a $PR\langle H(B),h,c,pp \rangle$ message

When receiving a PR message, proceed with Algorithm 1 & 2 defined in the paper.

## Transactions

For transactions, we want to keep it as simple as possible because it's not in the scope of a consensus protocol. We do not implement a global transaction mempool. Instead, each node will have its own pool. For testing convenience, we implement a client that sends transactions to all nodes all at once so that it is guaranteed to be picked up by a random proposer in the next round.

## Pipelining

### Background Block Dissemination

## Benchmarking

TODO
