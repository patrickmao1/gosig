// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.29.3
// source: rpc.proto

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

type SubmitTransactionReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tx *SignedTransaction `protobuf:"bytes,1,opt,name=tx,proto3" json:"tx,omitempty"`
}

func (x *SubmitTransactionReq) Reset() {
	*x = SubmitTransactionReq{}
	mi := &file_rpc_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SubmitTransactionReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitTransactionReq) ProtoMessage() {}

func (x *SubmitTransactionReq) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitTransactionReq.ProtoReflect.Descriptor instead.
func (*SubmitTransactionReq) Descriptor() ([]byte, []int) {
	return file_rpc_proto_rawDescGZIP(), []int{0}
}

func (x *SubmitTransactionReq) GetTx() *SignedTransaction {
	if x != nil {
		return x.Tx
	}
	return nil
}

type SubmitTransactionsReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Txs []*SignedTransaction `protobuf:"bytes,1,rep,name=txs,proto3" json:"txs,omitempty"`
}

func (x *SubmitTransactionsReq) Reset() {
	*x = SubmitTransactionsReq{}
	mi := &file_rpc_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SubmitTransactionsReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitTransactionsReq) ProtoMessage() {}

func (x *SubmitTransactionsReq) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitTransactionsReq.ProtoReflect.Descriptor instead.
func (*SubmitTransactionsReq) Descriptor() ([]byte, []int) {
	return file_rpc_proto_rawDescGZIP(), []int{1}
}

func (x *SubmitTransactionsReq) GetTxs() []*SignedTransaction {
	if x != nil {
		return x.Txs
	}
	return nil
}

type GetBalanceReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Account []byte `protobuf:"bytes,1,opt,name=account,proto3" json:"account,omitempty"`
}

func (x *GetBalanceReq) Reset() {
	*x = GetBalanceReq{}
	mi := &file_rpc_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetBalanceReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBalanceReq) ProtoMessage() {}

func (x *GetBalanceReq) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetBalanceReq.ProtoReflect.Descriptor instead.
func (*GetBalanceReq) Descriptor() ([]byte, []int) {
	return file_rpc_proto_rawDescGZIP(), []int{2}
}

func (x *GetBalanceReq) GetAccount() []byte {
	if x != nil {
		return x.Account
	}
	return nil
}

type GetBalanceRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Amount uint64 `protobuf:"varint,1,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (x *GetBalanceRes) Reset() {
	*x = GetBalanceRes{}
	mi := &file_rpc_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetBalanceRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBalanceRes) ProtoMessage() {}

func (x *GetBalanceRes) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetBalanceRes.ProtoReflect.Descriptor instead.
func (*GetBalanceRes) Descriptor() ([]byte, []int) {
	return file_rpc_proto_rawDescGZIP(), []int{3}
}

func (x *GetBalanceRes) GetAmount() uint64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	mi := &file_rpc_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_rpc_proto_rawDescGZIP(), []int{4}
}

type QueryTXsReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TxHashes [][]byte `protobuf:"bytes,1,rep,name=tx_hashes,json=txHashes,proto3" json:"tx_hashes,omitempty"`
}

func (x *QueryTXsReq) Reset() {
	*x = QueryTXsReq{}
	mi := &file_rpc_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueryTXsReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryTXsReq) ProtoMessage() {}

func (x *QueryTXsReq) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryTXsReq.ProtoReflect.Descriptor instead.
func (*QueryTXsReq) Descriptor() ([]byte, []int) {
	return file_rpc_proto_rawDescGZIP(), []int{5}
}

func (x *QueryTXsReq) GetTxHashes() [][]byte {
	if x != nil {
		return x.TxHashes
	}
	return nil
}

type QueryTXsRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Txs []*SignedTransaction `protobuf:"bytes,1,rep,name=txs,proto3" json:"txs,omitempty"`
}

func (x *QueryTXsRes) Reset() {
	*x = QueryTXsRes{}
	mi := &file_rpc_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueryTXsRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryTXsRes) ProtoMessage() {}

func (x *QueryTXsRes) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryTXsRes.ProtoReflect.Descriptor instead.
func (*QueryTXsRes) Descriptor() ([]byte, []int) {
	return file_rpc_proto_rawDescGZIP(), []int{6}
}

func (x *QueryTXsRes) GetTxs() []*SignedTransaction {
	if x != nil {
		return x.Txs
	}
	return nil
}

var File_rpc_proto protoreflect.FileDescriptor

var file_rpc_proto_rawDesc = []byte{
	0x0a, 0x09, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x67, 0x6f, 0x73,
	0x69, 0x67, 0x1a, 0x0b, 0x67, 0x6f, 0x73, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x40, 0x0a, 0x14, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x12, 0x28, 0x0a, 0x02, 0x74, 0x78, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x67, 0x6f, 0x73, 0x69, 0x67, 0x2e, 0x53, 0x69, 0x67, 0x6e,
	0x65, 0x64, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x02, 0x74,
	0x78, 0x22, 0x43, 0x0a, 0x15, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x12, 0x2a, 0x0a, 0x03, 0x74, 0x78,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x67, 0x6f, 0x73, 0x69, 0x67, 0x2e,
	0x53, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x03, 0x74, 0x78, 0x73, 0x22, 0x29, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x42, 0x61, 0x6c,
	0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x71, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x22, 0x27, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x52,
	0x65, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x22, 0x2a, 0x0a, 0x0b, 0x51, 0x75, 0x65, 0x72, 0x79, 0x54, 0x58, 0x73, 0x52,
	0x65, 0x71, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x78, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x65, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x08, 0x74, 0x78, 0x48, 0x61, 0x73, 0x68, 0x65, 0x73, 0x22,
	0x39, 0x0a, 0x0b, 0x51, 0x75, 0x65, 0x72, 0x79, 0x54, 0x58, 0x73, 0x52, 0x65, 0x73, 0x12, 0x2a,
	0x0a, 0x03, 0x74, 0x78, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x67, 0x6f,
	0x73, 0x69, 0x67, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x03, 0x74, 0x78, 0x73, 0x32, 0x9c, 0x02, 0x0a, 0x03, 0x52,
	0x50, 0x43, 0x12, 0x3e, 0x0a, 0x11, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x2e, 0x67, 0x6f, 0x73, 0x69, 0x67, 0x2e,
	0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x71, 0x1a, 0x0c, 0x2e, 0x67, 0x6f, 0x73, 0x69, 0x67, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x12, 0x40, 0x0a, 0x12, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x1c, 0x2e, 0x67, 0x6f, 0x73, 0x69, 0x67,
	0x2e, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x1a, 0x0c, 0x2e, 0x67, 0x6f, 0x73, 0x69, 0x67, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x12, 0x38, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x42, 0x61, 0x6c, 0x61, 0x6e,
	0x63, 0x65, 0x12, 0x14, 0x2e, 0x67, 0x6f, 0x73, 0x69, 0x67, 0x2e, 0x47, 0x65, 0x74, 0x42, 0x61,
	0x6c, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x14, 0x2e, 0x67, 0x6f, 0x73, 0x69, 0x67,
	0x2e, 0x47, 0x65, 0x74, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x73, 0x12, 0x25,
	0x0a, 0x04, 0x53, 0x65, 0x6e, 0x64, 0x12, 0x0f, 0x2e, 0x67, 0x6f, 0x73, 0x69, 0x67, 0x2e, 0x45,
	0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x1a, 0x0c, 0x2e, 0x67, 0x6f, 0x73, 0x69, 0x67, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x32, 0x0a, 0x08, 0x51, 0x75, 0x65, 0x72, 0x79, 0x54, 0x58,
	0x73, 0x12, 0x12, 0x2e, 0x67, 0x6f, 0x73, 0x69, 0x67, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x54,
	0x58, 0x73, 0x52, 0x65, 0x71, 0x1a, 0x12, 0x2e, 0x67, 0x6f, 0x73, 0x69, 0x67, 0x2e, 0x51, 0x75,
	0x65, 0x72, 0x79, 0x54, 0x58, 0x73, 0x52, 0x65, 0x73, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x2f, 0x74,
	0x79, 0x70, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc_proto_rawDescOnce sync.Once
	file_rpc_proto_rawDescData = file_rpc_proto_rawDesc
)

func file_rpc_proto_rawDescGZIP() []byte {
	file_rpc_proto_rawDescOnce.Do(func() {
		file_rpc_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_proto_rawDescData)
	})
	return file_rpc_proto_rawDescData
}

var file_rpc_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_rpc_proto_goTypes = []any{
	(*SubmitTransactionReq)(nil),  // 0: gosig.SubmitTransactionReq
	(*SubmitTransactionsReq)(nil), // 1: gosig.SubmitTransactionsReq
	(*GetBalanceReq)(nil),         // 2: gosig.GetBalanceReq
	(*GetBalanceRes)(nil),         // 3: gosig.GetBalanceRes
	(*Empty)(nil),                 // 4: gosig.Empty
	(*QueryTXsReq)(nil),           // 5: gosig.QueryTXsReq
	(*QueryTXsRes)(nil),           // 6: gosig.QueryTXsRes
	(*SignedTransaction)(nil),     // 7: gosig.SignedTransaction
	(*Envelope)(nil),              // 8: gosig.Envelope
}
var file_rpc_proto_depIdxs = []int32{
	7, // 0: gosig.SubmitTransactionReq.tx:type_name -> gosig.SignedTransaction
	7, // 1: gosig.SubmitTransactionsReq.txs:type_name -> gosig.SignedTransaction
	7, // 2: gosig.QueryTXsRes.txs:type_name -> gosig.SignedTransaction
	0, // 3: gosig.RPC.SubmitTransaction:input_type -> gosig.SubmitTransactionReq
	1, // 4: gosig.RPC.SubmitTransactions:input_type -> gosig.SubmitTransactionsReq
	2, // 5: gosig.RPC.GetBalance:input_type -> gosig.GetBalanceReq
	8, // 6: gosig.RPC.Send:input_type -> gosig.Envelope
	5, // 7: gosig.RPC.QueryTXs:input_type -> gosig.QueryTXsReq
	4, // 8: gosig.RPC.SubmitTransaction:output_type -> gosig.Empty
	4, // 9: gosig.RPC.SubmitTransactions:output_type -> gosig.Empty
	3, // 10: gosig.RPC.GetBalance:output_type -> gosig.GetBalanceRes
	4, // 11: gosig.RPC.Send:output_type -> gosig.Empty
	6, // 12: gosig.RPC.QueryTXs:output_type -> gosig.QueryTXsRes
	8, // [8:13] is the sub-list for method output_type
	3, // [3:8] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_rpc_proto_init() }
func file_rpc_proto_init() {
	if File_rpc_proto != nil {
		return
	}
	file_gosig_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_rpc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_rpc_proto_goTypes,
		DependencyIndexes: file_rpc_proto_depIdxs,
		MessageInfos:      file_rpc_proto_msgTypes,
	}.Build()
	File_rpc_proto = out.File
	file_rpc_proto_rawDesc = nil
	file_rpc_proto_goTypes = nil
	file_rpc_proto_depIdxs = nil
}
