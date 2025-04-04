// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.29.3
// source: gosig.proto

package types

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Transaction represents a simple payment of money. `from“ and `to` are
// pubkeys of users
type Transaction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	From   []byte `protobuf:"bytes,1,opt,name=from,proto3" json:"from,omitempty"`
	To     []byte `protobuf:"bytes,2,opt,name=to,proto3" json:"to,omitempty"`
	Amount uint64 `protobuf:"varint,3,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (x *Transaction) Reset() {
	*x = Transaction{}
	mi := &file_gosig_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Transaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transaction) ProtoMessage() {}

func (x *Transaction) ProtoReflect() protoreflect.Message {
	mi := &file_gosig_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Transaction.ProtoReflect.Descriptor instead.
func (*Transaction) Descriptor() ([]byte, []int) {
	return file_gosig_proto_rawDescGZIP(), []int{0}
}

func (x *Transaction) GetFrom() []byte {
	if x != nil {
		return x.From
	}
	return nil
}

func (x *Transaction) GetTo() []byte {
	if x != nil {
		return x.To
	}
	return nil
}

func (x *Transaction) GetAmount() uint64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

type SignedTransaction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tx  *Transaction `protobuf:"bytes,1,opt,name=tx,proto3" json:"tx,omitempty"`
	Sig []byte       `protobuf:"bytes,2,opt,name=sig,proto3" json:"sig,omitempty"`
}

func (x *SignedTransaction) Reset() {
	*x = SignedTransaction{}
	mi := &file_gosig_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SignedTransaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignedTransaction) ProtoMessage() {}

func (x *SignedTransaction) ProtoReflect() protoreflect.Message {
	mi := &file_gosig_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignedTransaction.ProtoReflect.Descriptor instead.
func (*SignedTransaction) Descriptor() ([]byte, []int) {
	return file_gosig_proto_rawDescGZIP(), []int{1}
}

func (x *SignedTransaction) GetTx() *Transaction {
	if x != nil {
		return x.Tx
	}
	return nil
}

func (x *SignedTransaction) GetSig() []byte {
	if x != nil {
		return x.Sig
	}
	return nil
}

type TransactionHashes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TxHashes [][]byte `protobuf:"bytes,1,rep,name=tx_hashes,json=txHashes,proto3" json:"tx_hashes,omitempty"`
	Root     []byte   `protobuf:"bytes,2,opt,name=root,proto3" json:"root,omitempty"` // root = hash(tx_hashes), for convenience
}

func (x *TransactionHashes) Reset() {
	*x = TransactionHashes{}
	mi := &file_gosig_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TransactionHashes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransactionHashes) ProtoMessage() {}

func (x *TransactionHashes) ProtoReflect() protoreflect.Message {
	mi := &file_gosig_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransactionHashes.ProtoReflect.Descriptor instead.
func (*TransactionHashes) Descriptor() ([]byte, []int) {
	return file_gosig_proto_rawDescGZIP(), []int{2}
}

func (x *TransactionHashes) GetTxHashes() [][]byte {
	if x != nil {
		return x.TxHashes
	}
	return nil
}

func (x *TransactionHashes) GetRoot() []byte {
	if x != nil {
		return x.Root
	}
	return nil
}

type SignedTransactions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Txs []*SignedTransaction `protobuf:"bytes,1,rep,name=txs,proto3" json:"txs,omitempty"`
}

func (x *SignedTransactions) Reset() {
	*x = SignedTransactions{}
	mi := &file_gosig_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SignedTransactions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignedTransactions) ProtoMessage() {}

func (x *SignedTransactions) ProtoReflect() protoreflect.Message {
	mi := &file_gosig_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignedTransactions.ProtoReflect.Descriptor instead.
func (*SignedTransactions) Descriptor() ([]byte, []int) {
	return file_gosig_proto_rawDescGZIP(), []int{3}
}

func (x *SignedTransactions) GetTxs() []*SignedTransaction {
	if x != nil {
		return x.Txs
	}
	return nil
}

type BlockHeader struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Height       uint32       `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
	ParentHash   []byte       `protobuf:"bytes,2,opt,name=parent_hash,json=parentHash,proto3" json:"parent_hash,omitempty"`
	ProposalCert *Certificate `protobuf:"bytes,3,opt,name=proposal_cert,json=proposalCert,proto3" json:"proposal_cert,omitempty"`
	// sign(privKey, Q^h), the preimage of the proposer score L^r
	ProposerProof []byte `protobuf:"bytes,4,opt,name=proposer_proof,json=proposerProof,proto3" json:"proposer_proof,omitempty"`
	// a flat hash of all transaction hashes in this block.
	// i.e. hash(hash(tx1) || hash(tx2) || ... || hash(txN))
	TxRoot []byte `protobuf:"bytes,5,opt,name=tx_root,json=txRoot,proto3" json:"tx_root,omitempty"`
}

func (x *BlockHeader) Reset() {
	*x = BlockHeader{}
	mi := &file_gosig_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BlockHeader) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockHeader) ProtoMessage() {}

func (x *BlockHeader) ProtoReflect() protoreflect.Message {
	mi := &file_gosig_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockHeader.ProtoReflect.Descriptor instead.
func (*BlockHeader) Descriptor() ([]byte, []int) {
	return file_gosig_proto_rawDescGZIP(), []int{4}
}

func (x *BlockHeader) GetHeight() uint32 {
	if x != nil {
		return x.Height
	}
	return 0
}

func (x *BlockHeader) GetParentHash() []byte {
	if x != nil {
		return x.ParentHash
	}
	return nil
}

func (x *BlockHeader) GetProposalCert() *Certificate {
	if x != nil {
		return x.ProposalCert
	}
	return nil
}

func (x *BlockHeader) GetProposerProof() []byte {
	if x != nil {
		return x.ProposerProof
	}
	return nil
}

func (x *BlockHeader) GetTxRoot() []byte {
	if x != nil {
		return x.TxRoot
	}
	return nil
}

type SignedMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sig            []byte `protobuf:"bytes,1,opt,name=sig,proto3" json:"sig,omitempty"`
	ValidatorIndex uint32 `protobuf:"varint,2,opt,name=validator_index,json=validatorIndex,proto3" json:"validator_index,omitempty"`
	// Types that are assignable to Message:
	//
	//	*SignedMessage_Bytes
	//	*SignedMessage_Proposal
	//	*SignedMessage_Prepare
	//	*SignedMessage_Tc
	Message  isSignedMessage_Message `protobuf_oneof:"message"`
	Deadline int64                   `protobuf:"varint,7,opt,name=deadline,proto3" json:"deadline,omitempty"` // deadline after which the message will no longer be relayed
}

func (x *SignedMessage) Reset() {
	*x = SignedMessage{}
	mi := &file_gosig_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SignedMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignedMessage) ProtoMessage() {}

func (x *SignedMessage) ProtoReflect() protoreflect.Message {
	mi := &file_gosig_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignedMessage.ProtoReflect.Descriptor instead.
func (*SignedMessage) Descriptor() ([]byte, []int) {
	return file_gosig_proto_rawDescGZIP(), []int{5}
}

func (x *SignedMessage) GetSig() []byte {
	if x != nil {
		return x.Sig
	}
	return nil
}

func (x *SignedMessage) GetValidatorIndex() uint32 {
	if x != nil {
		return x.ValidatorIndex
	}
	return 0
}

func (m *SignedMessage) GetMessage() isSignedMessage_Message {
	if m != nil {
		return m.Message
	}
	return nil
}

func (x *SignedMessage) GetBytes() []byte {
	if x, ok := x.GetMessage().(*SignedMessage_Bytes); ok {
		return x.Bytes
	}
	return nil
}

func (x *SignedMessage) GetProposal() *BlockProposal {
	if x, ok := x.GetMessage().(*SignedMessage_Proposal); ok {
		return x.Proposal
	}
	return nil
}

func (x *SignedMessage) GetPrepare() *PrepareCertificate {
	if x, ok := x.GetMessage().(*SignedMessage_Prepare); ok {
		return x.Prepare
	}
	return nil
}

func (x *SignedMessage) GetTc() *TentativeCommitCertificate {
	if x, ok := x.GetMessage().(*SignedMessage_Tc); ok {
		return x.Tc
	}
	return nil
}

func (x *SignedMessage) GetDeadline() int64 {
	if x != nil {
		return x.Deadline
	}
	return 0
}

type isSignedMessage_Message interface {
	isSignedMessage_Message()
}

type SignedMessage_Bytes struct {
	Bytes []byte `protobuf:"bytes,3,opt,name=bytes,proto3,oneof"` // arbitrary bytes for testing purpose
}

type SignedMessage_Proposal struct {
	Proposal *BlockProposal `protobuf:"bytes,4,opt,name=proposal,proto3,oneof"`
}

type SignedMessage_Prepare struct {
	Prepare *PrepareCertificate `protobuf:"bytes,5,opt,name=prepare,proto3,oneof"`
}

type SignedMessage_Tc struct {
	Tc *TentativeCommitCertificate `protobuf:"bytes,6,opt,name=tc,proto3,oneof"`
}

func (*SignedMessage_Bytes) isSignedMessage_Message() {}

func (*SignedMessage_Proposal) isSignedMessage_Message() {}

func (*SignedMessage_Prepare) isSignedMessage_Message() {}

func (*SignedMessage_Tc) isSignedMessage_Message() {}

type SignedMessages struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msgs []*SignedMessage `protobuf:"bytes,1,rep,name=msgs,proto3" json:"msgs,omitempty"`
}

func (x *SignedMessages) Reset() {
	*x = SignedMessages{}
	mi := &file_gosig_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SignedMessages) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignedMessages) ProtoMessage() {}

func (x *SignedMessages) ProtoReflect() protoreflect.Message {
	mi := &file_gosig_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignedMessages.ProtoReflect.Descriptor instead.
func (*SignedMessages) Descriptor() ([]byte, []int) {
	return file_gosig_proto_rawDescGZIP(), []int{6}
}

func (x *SignedMessages) GetMsgs() []*SignedMessage {
	if x != nil {
		return x.Msgs
	}
	return nil
}

type BlockProposal struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Round       uint32       `protobuf:"varint,1,opt,name=round,proto3" json:"round,omitempty"`
	BlockHeader *BlockHeader `protobuf:"bytes,2,opt,name=block_header,json=blockHeader,proto3" json:"block_header,omitempty"`
}

func (x *BlockProposal) Reset() {
	*x = BlockProposal{}
	mi := &file_gosig_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BlockProposal) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockProposal) ProtoMessage() {}

func (x *BlockProposal) ProtoReflect() protoreflect.Message {
	mi := &file_gosig_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockProposal.ProtoReflect.Descriptor instead.
func (*BlockProposal) Descriptor() ([]byte, []int) {
	return file_gosig_proto_rawDescGZIP(), []int{7}
}

func (x *BlockProposal) GetRound() uint32 {
	if x != nil {
		return x.Round
	}
	return 0
}

func (x *BlockProposal) GetBlockHeader() *BlockHeader {
	if x != nil {
		return x.BlockHeader
	}
	return nil
}

type PrepareCertificate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg  *Prepare     `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	Cert *Certificate `protobuf:"bytes,2,opt,name=cert,proto3" json:"cert,omitempty"` // can be a partial certificate (e.g. signed by less than 2/3)
}

func (x *PrepareCertificate) Reset() {
	*x = PrepareCertificate{}
	mi := &file_gosig_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PrepareCertificate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrepareCertificate) ProtoMessage() {}

func (x *PrepareCertificate) ProtoReflect() protoreflect.Message {
	mi := &file_gosig_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PrepareCertificate.ProtoReflect.Descriptor instead.
func (*PrepareCertificate) Descriptor() ([]byte, []int) {
	return file_gosig_proto_rawDescGZIP(), []int{8}
}

func (x *PrepareCertificate) GetMsg() *Prepare {
	if x != nil {
		return x.Msg
	}
	return nil
}

func (x *PrepareCertificate) GetCert() *Certificate {
	if x != nil {
		return x.Cert
	}
	return nil
}

type Prepare struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BlockHash   []byte `protobuf:"bytes,1,opt,name=block_hash,json=blockHash,proto3" json:"block_hash,omitempty"`
	BlockHeight uint32 `protobuf:"varint,2,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
}

func (x *Prepare) Reset() {
	*x = Prepare{}
	mi := &file_gosig_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Prepare) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Prepare) ProtoMessage() {}

func (x *Prepare) ProtoReflect() protoreflect.Message {
	mi := &file_gosig_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Prepare.ProtoReflect.Descriptor instead.
func (*Prepare) Descriptor() ([]byte, []int) {
	return file_gosig_proto_rawDescGZIP(), []int{9}
}

func (x *Prepare) GetBlockHash() []byte {
	if x != nil {
		return x.BlockHash
	}
	return nil
}

func (x *Prepare) GetBlockHeight() uint32 {
	if x != nil {
		return x.BlockHeight
	}
	return 0
}

type TentativeCommitCertificate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg  *TentativeCommit `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	Cert *Certificate     `protobuf:"bytes,2,opt,name=cert,proto3" json:"cert,omitempty"` // can be a partial certificate (e.g. signed by less than 2/3)
}

func (x *TentativeCommitCertificate) Reset() {
	*x = TentativeCommitCertificate{}
	mi := &file_gosig_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TentativeCommitCertificate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TentativeCommitCertificate) ProtoMessage() {}

func (x *TentativeCommitCertificate) ProtoReflect() protoreflect.Message {
	mi := &file_gosig_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TentativeCommitCertificate.ProtoReflect.Descriptor instead.
func (*TentativeCommitCertificate) Descriptor() ([]byte, []int) {
	return file_gosig_proto_rawDescGZIP(), []int{10}
}

func (x *TentativeCommitCertificate) GetMsg() *TentativeCommit {
	if x != nil {
		return x.Msg
	}
	return nil
}

func (x *TentativeCommitCertificate) GetCert() *Certificate {
	if x != nil {
		return x.Cert
	}
	return nil
}

type TentativeCommit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BlockHash   []byte `protobuf:"bytes,1,opt,name=block_hash,json=blockHash,proto3" json:"block_hash,omitempty"`
	BlockHeight uint32 `protobuf:"varint,2,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
}

func (x *TentativeCommit) Reset() {
	*x = TentativeCommit{}
	mi := &file_gosig_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TentativeCommit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TentativeCommit) ProtoMessage() {}

func (x *TentativeCommit) ProtoReflect() protoreflect.Message {
	mi := &file_gosig_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TentativeCommit.ProtoReflect.Descriptor instead.
func (*TentativeCommit) Descriptor() ([]byte, []int) {
	return file_gosig_proto_rawDescGZIP(), []int{11}
}

func (x *TentativeCommit) GetBlockHash() []byte {
	if x != nil {
		return x.BlockHash
	}
	return nil
}

func (x *TentativeCommit) GetBlockHeight() uint32 {
	if x != nil {
		return x.BlockHeight
	}
	return 0
}

type Certificate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AggSig    []byte   `protobuf:"bytes,1,opt,name=agg_sig,json=aggSig,proto3" json:"agg_sig,omitempty"`                  // aggregated signatures
	SigCounts []uint32 `protobuf:"varint,2,rep,packed,name=sig_counts,json=sigCounts,proto3" json:"sig_counts,omitempty"` // an array of how many times each validator provided the signatures
	Round     uint32   `protobuf:"varint,3,opt,name=round,proto3" json:"round,omitempty"`                                 // the round in which this certificate is formed
}

func (x *Certificate) Reset() {
	*x = Certificate{}
	mi := &file_gosig_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Certificate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Certificate) ProtoMessage() {}

func (x *Certificate) ProtoReflect() protoreflect.Message {
	mi := &file_gosig_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Certificate.ProtoReflect.Descriptor instead.
func (*Certificate) Descriptor() ([]byte, []int) {
	return file_gosig_proto_rawDescGZIP(), []int{12}
}

func (x *Certificate) GetAggSig() []byte {
	if x != nil {
		return x.AggSig
	}
	return nil
}

func (x *Certificate) GetSigCounts() []uint32 {
	if x != nil {
		return x.SigCounts
	}
	return nil
}

func (x *Certificate) GetRound() uint32 {
	if x != nil {
		return x.Round
	}
	return 0
}

var File_gosig_proto protoreflect.FileDescriptor

var file_gosig_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x67, 0x6f, 0x73, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x68,
	0x79, 0x64, 0x66, 0x73, 0x22, 0x49, 0x0a, 0x0b, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x74, 0x6f, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x02, 0x74, 0x6f, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22,
	0x49, 0x0a, 0x11, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x22, 0x0a, 0x02, 0x74, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x12, 0x2e, 0x68, 0x79, 0x64, 0x66, 0x73, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x02, 0x74, 0x78, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x69, 0x67, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x73, 0x69, 0x67, 0x22, 0x44, 0x0a, 0x11, 0x54, 0x72,
	0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x61, 0x73, 0x68, 0x65, 0x73, 0x12,
	0x1b, 0x0a, 0x09, 0x74, 0x78, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0c, 0x52, 0x08, 0x74, 0x78, 0x48, 0x61, 0x73, 0x68, 0x65, 0x73, 0x12, 0x12, 0x0a, 0x04,
	0x72, 0x6f, 0x6f, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x72, 0x6f, 0x6f, 0x74,
	0x22, 0x40, 0x0a, 0x12, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x2a, 0x0a, 0x03, 0x74, 0x78, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x68, 0x79, 0x64, 0x66, 0x73, 0x2e, 0x53, 0x69, 0x67, 0x6e,
	0x65, 0x64, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x03, 0x74,
	0x78, 0x73, 0x22, 0xbf, 0x01, 0x0a, 0x0b, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x61,
	0x72, 0x65, 0x6e, 0x74, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x0a, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x68, 0x12, 0x37, 0x0a, 0x0d, 0x70,
	0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x5f, 0x63, 0x65, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x12, 0x2e, 0x68, 0x79, 0x64, 0x66, 0x73, 0x2e, 0x43, 0x65, 0x72, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x0c, 0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c,
	0x43, 0x65, 0x72, 0x74, 0x12, 0x25, 0x0a, 0x0e, 0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x65, 0x72,
	0x5f, 0x70, 0x72, 0x6f, 0x6f, 0x66, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0d, 0x70, 0x72,
	0x6f, 0x70, 0x6f, 0x73, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x12, 0x17, 0x0a, 0x07, 0x74,
	0x78, 0x5f, 0x72, 0x6f, 0x6f, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x74, 0x78,
	0x52, 0x6f, 0x6f, 0x74, 0x22, 0xa9, 0x02, 0x0a, 0x0d, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x69, 0x67, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x03, 0x73, 0x69, 0x67, 0x12, 0x27, 0x0a, 0x0f, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x6f, 0x72, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x0e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x49, 0x6e, 0x64, 0x65,
	0x78, 0x12, 0x16, 0x0a, 0x05, 0x62, 0x79, 0x74, 0x65, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c,
	0x48, 0x00, 0x52, 0x05, 0x62, 0x79, 0x74, 0x65, 0x73, 0x12, 0x32, 0x0a, 0x08, 0x70, 0x72, 0x6f,
	0x70, 0x6f, 0x73, 0x61, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x68, 0x79,
	0x64, 0x66, 0x73, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61,
	0x6c, 0x48, 0x00, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x12, 0x35, 0x0a,
	0x07, 0x70, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19,
	0x2e, 0x68, 0x79, 0x64, 0x66, 0x73, 0x2e, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x43, 0x65,
	0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x48, 0x00, 0x52, 0x07, 0x70, 0x72, 0x65,
	0x70, 0x61, 0x72, 0x65, 0x12, 0x33, 0x0a, 0x02, 0x74, 0x63, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x21, 0x2e, 0x68, 0x79, 0x64, 0x66, 0x73, 0x2e, 0x54, 0x65, 0x6e, 0x74, 0x61, 0x74, 0x69,
	0x76, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x65, 0x48, 0x00, 0x52, 0x02, 0x74, 0x63, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x65, 0x61,
	0x64, 0x6c, 0x69, 0x6e, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x64, 0x65, 0x61,
	0x64, 0x6c, 0x69, 0x6e, 0x65, 0x42, 0x09, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x22, 0x3a, 0x0a, 0x0e, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x73, 0x12, 0x28, 0x0a, 0x04, 0x6d, 0x73, 0x67, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x14, 0x2e, 0x68, 0x79, 0x64, 0x66, 0x73, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x04, 0x6d, 0x73, 0x67, 0x73, 0x22, 0x5c, 0x0a, 0x0d,
	0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x12, 0x14, 0x0a,
	0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x72, 0x6f,
	0x75, 0x6e, 0x64, 0x12, 0x35, 0x0a, 0x0c, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x68, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x68, 0x79, 0x64, 0x66,
	0x73, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x52, 0x0b, 0x62,
	0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x22, 0x5e, 0x0a, 0x12, 0x50, 0x72,
	0x65, 0x70, 0x61, 0x72, 0x65, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65,
	0x12, 0x20, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e,
	0x68, 0x79, 0x64, 0x66, 0x73, 0x2e, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x52, 0x03, 0x6d,
	0x73, 0x67, 0x12, 0x26, 0x0a, 0x04, 0x63, 0x65, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x12, 0x2e, 0x68, 0x79, 0x64, 0x66, 0x73, 0x2e, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x65, 0x52, 0x04, 0x63, 0x65, 0x72, 0x74, 0x22, 0x4b, 0x0a, 0x07, 0x50, 0x72,
	0x65, 0x70, 0x61, 0x72, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x68,
	0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b,
	0x48, 0x61, 0x73, 0x68, 0x12, 0x21, 0x0a, 0x0c, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x68, 0x65,
	0x69, 0x67, 0x68, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0b, 0x62, 0x6c, 0x6f, 0x63,
	0x6b, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x22, 0x6e, 0x0a, 0x1a, 0x54, 0x65, 0x6e, 0x74, 0x61,
	0x74, 0x69, 0x76, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x65, 0x12, 0x28, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x16, 0x2e, 0x68, 0x79, 0x64, 0x66, 0x73, 0x2e, 0x54, 0x65, 0x6e, 0x74, 0x61,
	0x74, 0x69, 0x76, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x12,
	0x26, 0x0a, 0x04, 0x63, 0x65, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e,
	0x68, 0x79, 0x64, 0x66, 0x73, 0x2e, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x65, 0x52, 0x04, 0x63, 0x65, 0x72, 0x74, 0x22, 0x53, 0x0a, 0x0f, 0x54, 0x65, 0x6e, 0x74, 0x61,
	0x74, 0x69, 0x76, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x62, 0x6c,
	0x6f, 0x63, 0x6b, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x12, 0x21, 0x0a, 0x0c, 0x62, 0x6c, 0x6f,
	0x63, 0x6b, 0x5f, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x0b, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x22, 0x5b, 0x0a, 0x0b,
	0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x61,
	0x67, 0x67, 0x5f, 0x73, 0x69, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x61, 0x67,
	0x67, 0x53, 0x69, 0x67, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x69, 0x67, 0x5f, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x09, 0x73, 0x69, 0x67, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x2f, 0x74,
	0x79, 0x70, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gosig_proto_rawDescOnce sync.Once
	file_gosig_proto_rawDescData = file_gosig_proto_rawDesc
)

func file_gosig_proto_rawDescGZIP() []byte {
	file_gosig_proto_rawDescOnce.Do(func() {
		file_gosig_proto_rawDescData = protoimpl.X.CompressGZIP(file_gosig_proto_rawDescData)
	})
	return file_gosig_proto_rawDescData
}

var file_gosig_proto_msgTypes = make([]protoimpl.MessageInfo, 13)
var file_gosig_proto_goTypes = []any{
	(*Transaction)(nil),                // 0: hydfs.Transaction
	(*SignedTransaction)(nil),          // 1: hydfs.SignedTransaction
	(*TransactionHashes)(nil),          // 2: hydfs.TransactionHashes
	(*SignedTransactions)(nil),         // 3: hydfs.SignedTransactions
	(*BlockHeader)(nil),                // 4: hydfs.BlockHeader
	(*SignedMessage)(nil),              // 5: hydfs.SignedMessage
	(*SignedMessages)(nil),             // 6: hydfs.SignedMessages
	(*BlockProposal)(nil),              // 7: hydfs.BlockProposal
	(*PrepareCertificate)(nil),         // 8: hydfs.PrepareCertificate
	(*Prepare)(nil),                    // 9: hydfs.Prepare
	(*TentativeCommitCertificate)(nil), // 10: hydfs.TentativeCommitCertificate
	(*TentativeCommit)(nil),            // 11: hydfs.TentativeCommit
	(*Certificate)(nil),                // 12: hydfs.Certificate
}
var file_gosig_proto_depIdxs = []int32{
	0,  // 0: hydfs.SignedTransaction.tx:type_name -> hydfs.Transaction
	1,  // 1: hydfs.SignedTransactions.txs:type_name -> hydfs.SignedTransaction
	12, // 2: hydfs.BlockHeader.proposal_cert:type_name -> hydfs.Certificate
	7,  // 3: hydfs.SignedMessage.proposal:type_name -> hydfs.BlockProposal
	8,  // 4: hydfs.SignedMessage.prepare:type_name -> hydfs.PrepareCertificate
	10, // 5: hydfs.SignedMessage.tc:type_name -> hydfs.TentativeCommitCertificate
	5,  // 6: hydfs.SignedMessages.msgs:type_name -> hydfs.SignedMessage
	4,  // 7: hydfs.BlockProposal.block_header:type_name -> hydfs.BlockHeader
	9,  // 8: hydfs.PrepareCertificate.msg:type_name -> hydfs.Prepare
	12, // 9: hydfs.PrepareCertificate.cert:type_name -> hydfs.Certificate
	11, // 10: hydfs.TentativeCommitCertificate.msg:type_name -> hydfs.TentativeCommit
	12, // 11: hydfs.TentativeCommitCertificate.cert:type_name -> hydfs.Certificate
	12, // [12:12] is the sub-list for method output_type
	12, // [12:12] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_gosig_proto_init() }
func file_gosig_proto_init() {
	if File_gosig_proto != nil {
		return
	}
	file_gosig_proto_msgTypes[5].OneofWrappers = []any{
		(*SignedMessage_Bytes)(nil),
		(*SignedMessage_Proposal)(nil),
		(*SignedMessage_Prepare)(nil),
		(*SignedMessage_Tc)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_gosig_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   13,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_gosig_proto_goTypes,
		DependencyIndexes: file_gosig_proto_depIdxs,
		MessageInfos:      file_gosig_proto_msgTypes,
	}.Build()
	File_gosig_proto = out.File
	file_gosig_proto_rawDesc = nil
	file_gosig_proto_goTypes = nil
	file_gosig_proto_depIdxs = nil
}
