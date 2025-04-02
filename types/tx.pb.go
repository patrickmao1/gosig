// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.29.3
// source: tx.proto

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
	mi := &file_tx_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SubmitTransactionReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitTransactionReq) ProtoMessage() {}

func (x *SubmitTransactionReq) ProtoReflect() protoreflect.Message {
	mi := &file_tx_proto_msgTypes[0]
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
	return file_tx_proto_rawDescGZIP(), []int{0}
}

func (x *SubmitTransactionReq) GetTx() *SignedTransaction {
	if x != nil {
		return x.Tx
	}
	return nil
}

type SubmitTransactionRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Err *string `protobuf:"bytes,1,opt,name=err,proto3,oneof" json:"err,omitempty"`
}

func (x *SubmitTransactionRes) Reset() {
	*x = SubmitTransactionRes{}
	mi := &file_tx_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SubmitTransactionRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitTransactionRes) ProtoMessage() {}

func (x *SubmitTransactionRes) ProtoReflect() protoreflect.Message {
	mi := &file_tx_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitTransactionRes.ProtoReflect.Descriptor instead.
func (*SubmitTransactionRes) Descriptor() ([]byte, []int) {
	return file_tx_proto_rawDescGZIP(), []int{1}
}

func (x *SubmitTransactionRes) GetErr() string {
	if x != nil && x.Err != nil {
		return *x.Err
	}
	return ""
}

type TransactionStatusReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TxHash []byte `protobuf:"bytes,1,opt,name=tx_hash,json=txHash,proto3" json:"tx_hash,omitempty"`
}

func (x *TransactionStatusReq) Reset() {
	*x = TransactionStatusReq{}
	mi := &file_tx_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TransactionStatusReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransactionStatusReq) ProtoMessage() {}

func (x *TransactionStatusReq) ProtoReflect() protoreflect.Message {
	mi := &file_tx_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransactionStatusReq.ProtoReflect.Descriptor instead.
func (*TransactionStatusReq) Descriptor() ([]byte, []int) {
	return file_tx_proto_rawDescGZIP(), []int{2}
}

func (x *TransactionStatusReq) GetTxHash() []byte {
	if x != nil {
		return x.TxHash
	}
	return nil
}

type TransactionStatusRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Err    *string `protobuf:"bytes,1,opt,name=err,proto3,oneof" json:"err,omitempty"`
	Status string  `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *TransactionStatusRes) Reset() {
	*x = TransactionStatusRes{}
	mi := &file_tx_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TransactionStatusRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransactionStatusRes) ProtoMessage() {}

func (x *TransactionStatusRes) ProtoReflect() protoreflect.Message {
	mi := &file_tx_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransactionStatusRes.ProtoReflect.Descriptor instead.
func (*TransactionStatusRes) Descriptor() ([]byte, []int) {
	return file_tx_proto_rawDescGZIP(), []int{3}
}

func (x *TransactionStatusRes) GetErr() string {
	if x != nil && x.Err != nil {
		return *x.Err
	}
	return ""
}

func (x *TransactionStatusRes) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

var File_tx_proto protoreflect.FileDescriptor

var file_tx_proto_rawDesc = []byte{
	0x0a, 0x08, 0x74, 0x78, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x68, 0x79, 0x64, 0x66,
	0x73, 0x1a, 0x0b, 0x67, 0x6f, 0x73, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x40,
	0x0a, 0x14, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x12, 0x28, 0x0a, 0x02, 0x74, 0x78, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x18, 0x2e, 0x68, 0x79, 0x64, 0x66, 0x73, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x65,
	0x64, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x02, 0x74, 0x78,
	0x22, 0x35, 0x0a, 0x14, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x12, 0x15, 0x0a, 0x03, 0x65, 0x72, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x03, 0x65, 0x72, 0x72, 0x88, 0x01, 0x01, 0x42,
	0x06, 0x0a, 0x04, 0x5f, 0x65, 0x72, 0x72, 0x22, 0x2f, 0x0a, 0x14, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x12,
	0x17, 0x0a, 0x07, 0x74, 0x78, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x06, 0x74, 0x78, 0x48, 0x61, 0x73, 0x68, 0x22, 0x4d, 0x0a, 0x14, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73,
	0x12, 0x15, 0x0a, 0x03, 0x65, 0x72, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52,
	0x03, 0x65, 0x72, 0x72, 0x88, 0x01, 0x01, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x42,
	0x06, 0x0a, 0x04, 0x5f, 0x65, 0x72, 0x72, 0x32, 0xa4, 0x01, 0x0a, 0x04, 0x4e, 0x6f, 0x64, 0x65,
	0x12, 0x4d, 0x0a, 0x11, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x2e, 0x68, 0x79, 0x64, 0x66, 0x73, 0x2e, 0x53, 0x75,
	0x62, 0x6d, 0x69, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x71, 0x1a, 0x1b, 0x2e, 0x68, 0x79, 0x64, 0x66, 0x73, 0x2e, 0x53, 0x75, 0x62, 0x6d, 0x69,
	0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x12,
	0x4d, 0x0a, 0x11, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x1b, 0x2e, 0x68, 0x79, 0x64, 0x66, 0x73, 0x2e, 0x54, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65,
	0x71, 0x1a, 0x1b, 0x2e, 0x68, 0x79, 0x64, 0x66, 0x73, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x42, 0x09,
	0x5a, 0x07, 0x2e, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_tx_proto_rawDescOnce sync.Once
	file_tx_proto_rawDescData = file_tx_proto_rawDesc
)

func file_tx_proto_rawDescGZIP() []byte {
	file_tx_proto_rawDescOnce.Do(func() {
		file_tx_proto_rawDescData = protoimpl.X.CompressGZIP(file_tx_proto_rawDescData)
	})
	return file_tx_proto_rawDescData
}

var file_tx_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_tx_proto_goTypes = []any{
	(*SubmitTransactionReq)(nil), // 0: hydfs.SubmitTransactionReq
	(*SubmitTransactionRes)(nil), // 1: hydfs.SubmitTransactionRes
	(*TransactionStatusReq)(nil), // 2: hydfs.TransactionStatusReq
	(*TransactionStatusRes)(nil), // 3: hydfs.TransactionStatusRes
	(*SignedTransaction)(nil),    // 4: hydfs.SignedTransaction
}
var file_tx_proto_depIdxs = []int32{
	4, // 0: hydfs.SubmitTransactionReq.tx:type_name -> hydfs.SignedTransaction
	0, // 1: hydfs.Node.SubmitTransaction:input_type -> hydfs.SubmitTransactionReq
	2, // 2: hydfs.Node.TransactionStatus:input_type -> hydfs.TransactionStatusReq
	1, // 3: hydfs.Node.SubmitTransaction:output_type -> hydfs.SubmitTransactionRes
	3, // 4: hydfs.Node.TransactionStatus:output_type -> hydfs.TransactionStatusRes
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_tx_proto_init() }
func file_tx_proto_init() {
	if File_tx_proto != nil {
		return
	}
	file_gosig_proto_init()
	file_tx_proto_msgTypes[1].OneofWrappers = []any{}
	file_tx_proto_msgTypes[3].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_tx_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_tx_proto_goTypes,
		DependencyIndexes: file_tx_proto_depIdxs,
		MessageInfos:      file_tx_proto_msgTypes,
	}.Build()
	File_tx_proto = out.File
	file_tx_proto_rawDesc = nil
	file_tx_proto_goTypes = nil
	file_tx_proto_depIdxs = nil
}
