// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v5.28.0--dev
// source: block.proto

package proto

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

// BlockHeader contains metadata about the block
type BlockHeader struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Height           uint64 `protobuf:"varint,1,opt,name=Height,proto3" json:"Height,omitempty"`                    // Block height
	ParentHash       []byte `protobuf:"bytes,2,opt,name=ParentHash,proto3" json:"ParentHash,omitempty"`             // Hash of the parent block
	StateRoot        []byte `protobuf:"bytes,3,opt,name=StateRoot,proto3" json:"StateRoot,omitempty"`               // Merkle root of state
	TransactionsRoot []byte `protobuf:"bytes,4,opt,name=TransactionsRoot,proto3" json:"TransactionsRoot,omitempty"` // Merkle root of transactions
	ReceiptRoot      []byte `protobuf:"bytes,5,opt,name=ReceiptRoot,proto3" json:"ReceiptRoot,omitempty"`           // Merkle root of receipts (optional)
	Timestamp        uint64 `protobuf:"varint,6,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`              // Unix time
	Prevrandao       []byte `protobuf:"bytes,7,opt,name=Prevrandao,proto3" json:"Prevrandao,omitempty"`             // Randomness from beacon chain (PoS replacement for difficulty)
	Proposer         []byte `protobuf:"bytes,8,opt,name=Proposer,proto3" json:"Proposer,omitempty"`                 // Address of the block proposer/leader
	Signature        []byte `protobuf:"bytes,9,opt,name=Signature,proto3" json:"Signature,omitempty"`               // Signature of the proposer
	GasUsed          uint64 `protobuf:"varint,10,opt,name=GasUsed,proto3" json:"GasUsed,omitempty"`                 // Arbitrary extra data (optional)
	ExtraData        []byte `protobuf:"bytes,11,opt,name=ExtraData,proto3" json:"ExtraData,omitempty"`
	UncleHash        []byte `protobuf:"bytes,12,opt,name=UncleHash,proto3" json:"UncleHash,omitempty"`
	Bloom            []byte `protobuf:"bytes,13,opt,name=Bloom,proto3" json:"Bloom,omitempty"`
}

func (x *BlockHeader) Reset() {
	*x = BlockHeader{}
	if protoimpl.UnsafeEnabled {
		mi := &file_block_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlockHeader) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockHeader) ProtoMessage() {}

func (x *BlockHeader) ProtoReflect() protoreflect.Message {
	mi := &file_block_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
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
	return file_block_proto_rawDescGZIP(), []int{0}
}

func (x *BlockHeader) GetHeight() uint64 {
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

func (x *BlockHeader) GetStateRoot() []byte {
	if x != nil {
		return x.StateRoot
	}
	return nil
}

func (x *BlockHeader) GetTransactionsRoot() []byte {
	if x != nil {
		return x.TransactionsRoot
	}
	return nil
}

func (x *BlockHeader) GetReceiptRoot() []byte {
	if x != nil {
		return x.ReceiptRoot
	}
	return nil
}

func (x *BlockHeader) GetTimestamp() uint64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *BlockHeader) GetPrevrandao() []byte {
	if x != nil {
		return x.Prevrandao
	}
	return nil
}

func (x *BlockHeader) GetProposer() []byte {
	if x != nil {
		return x.Proposer
	}
	return nil
}

func (x *BlockHeader) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

func (x *BlockHeader) GetGasUsed() uint64 {
	if x != nil {
		return x.GasUsed
	}
	return 0
}

func (x *BlockHeader) GetExtraData() []byte {
	if x != nil {
		return x.ExtraData
	}
	return nil
}

func (x *BlockHeader) GetUncleHash() []byte {
	if x != nil {
		return x.UncleHash
	}
	return nil
}

func (x *BlockHeader) GetBloom() []byte {
	if x != nil {
		return x.Bloom
	}
	return nil
}

// BlockHashData is used to calculate the block hash (does not include Signature)
type BlockHashData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Height           uint64 `protobuf:"varint,1,opt,name=Height,proto3" json:"Height,omitempty"`
	ParentHash       []byte `protobuf:"bytes,2,opt,name=ParentHash,proto3" json:"ParentHash,omitempty"`
	StateRoot        []byte `protobuf:"bytes,3,opt,name=StateRoot,proto3" json:"StateRoot,omitempty"`
	TransactionsRoot []byte `protobuf:"bytes,4,opt,name=TransactionsRoot,proto3" json:"TransactionsRoot,omitempty"`
	ReceiptRoot      []byte `protobuf:"bytes,5,opt,name=ReceiptRoot,proto3" json:"ReceiptRoot,omitempty"`
	Timestamp        uint64 `protobuf:"varint,6,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
	Prevrandao       []byte `protobuf:"bytes,7,opt,name=Prevrandao,proto3" json:"Prevrandao,omitempty"` // Included in hash calculation
	Proposer         []byte `protobuf:"bytes,8,opt,name=Proposer,proto3" json:"Proposer,omitempty"`
	GasUsed          uint64 `protobuf:"varint,9,opt,name=GasUsed,proto3" json:"GasUsed,omitempty"` // Arbitrary extra data (optional)
	ExtraData        []byte `protobuf:"bytes,10,opt,name=ExtraData,proto3" json:"ExtraData,omitempty"`
	UncleHash        []byte `protobuf:"bytes,11,opt,name=UncleHash,proto3" json:"UncleHash,omitempty"`
	Bloom            []byte `protobuf:"bytes,12,opt,name=Bloom,proto3" json:"Bloom,omitempty"`
}

func (x *BlockHashData) Reset() {
	*x = BlockHashData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_block_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlockHashData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockHashData) ProtoMessage() {}

func (x *BlockHashData) ProtoReflect() protoreflect.Message {
	mi := &file_block_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockHashData.ProtoReflect.Descriptor instead.
func (*BlockHashData) Descriptor() ([]byte, []int) {
	return file_block_proto_rawDescGZIP(), []int{1}
}

func (x *BlockHashData) GetHeight() uint64 {
	if x != nil {
		return x.Height
	}
	return 0
}

func (x *BlockHashData) GetParentHash() []byte {
	if x != nil {
		return x.ParentHash
	}
	return nil
}

func (x *BlockHashData) GetStateRoot() []byte {
	if x != nil {
		return x.StateRoot
	}
	return nil
}

func (x *BlockHashData) GetTransactionsRoot() []byte {
	if x != nil {
		return x.TransactionsRoot
	}
	return nil
}

func (x *BlockHashData) GetReceiptRoot() []byte {
	if x != nil {
		return x.ReceiptRoot
	}
	return nil
}

func (x *BlockHashData) GetTimestamp() uint64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *BlockHashData) GetPrevrandao() []byte {
	if x != nil {
		return x.Prevrandao
	}
	return nil
}

func (x *BlockHashData) GetProposer() []byte {
	if x != nil {
		return x.Proposer
	}
	return nil
}

func (x *BlockHashData) GetGasUsed() uint64 {
	if x != nil {
		return x.GasUsed
	}
	return 0
}

func (x *BlockHashData) GetExtraData() []byte {
	if x != nil {
		return x.ExtraData
	}
	return nil
}

func (x *BlockHashData) GetUncleHash() []byte {
	if x != nil {
		return x.UncleHash
	}
	return nil
}

func (x *BlockHashData) GetBloom() []byte {
	if x != nil {
		return x.Bloom
	}
	return nil
}

// Block contains the header and the list of transactions
type Block struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Header *BlockHeader   `protobuf:"bytes,1,opt,name=Header,proto3" json:"Header,omitempty"` // Header of the block
	Txns   []*Transaction `protobuf:"bytes,2,rep,name=Txns,proto3" json:"Txns,omitempty"`     // List of transactions
	Uncles []*BlockHeader `protobuf:"bytes,3,rep,name=Uncles,proto3" json:"Uncles,omitempty"`
}

func (x *Block) Reset() {
	*x = Block{}
	if protoimpl.UnsafeEnabled {
		mi := &file_block_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Block) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Block) ProtoMessage() {}

func (x *Block) ProtoReflect() protoreflect.Message {
	mi := &file_block_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Block.ProtoReflect.Descriptor instead.
func (*Block) Descriptor() ([]byte, []int) {
	return file_block_proto_rawDescGZIP(), []int{2}
}

func (x *Block) GetHeader() *BlockHeader {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *Block) GetTxns() []*Transaction {
	if x != nil {
		return x.Txns
	}
	return nil
}

func (x *Block) GetUncles() []*BlockHeader {
	if x != nil {
		return x.Uncles
	}
	return nil
}

type Body struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Txns   []*Transaction `protobuf:"bytes,1,rep,name=Txns,proto3" json:"Txns,omitempty"` // List of transactions
	Uncles []*BlockHeader `protobuf:"bytes,2,rep,name=Uncles,proto3" json:"Uncles,omitempty"`
}

func (x *Body) Reset() {
	*x = Body{}
	if protoimpl.UnsafeEnabled {
		mi := &file_block_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Body) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Body) ProtoMessage() {}

func (x *Body) ProtoReflect() protoreflect.Message {
	mi := &file_block_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Body.ProtoReflect.Descriptor instead.
func (*Body) Descriptor() ([]byte, []int) {
	return file_block_proto_rawDescGZIP(), []int{3}
}

func (x *Body) GetTxns() []*Transaction {
	if x != nil {
		return x.Txns
	}
	return nil
}

func (x *Body) GetUncles() []*BlockHeader {
	if x != nil {
		return x.Uncles
	}
	return nil
}

var File_block_proto protoreflect.FileDescriptor

var file_block_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x62,
	0x6c, 0x6f, 0x63, 0x6b, 0x1a, 0x11, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x95, 0x03, 0x0a, 0x0b, 0x42, 0x6c, 0x6f, 0x63,
	0x6b, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x48, 0x65, 0x69, 0x67, 0x68,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12,
	0x1e, 0x0a, 0x0a, 0x50, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x0a, 0x50, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x68, 0x12,
	0x1c, 0x0a, 0x09, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x6f, 0x6f, 0x74, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x09, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x6f, 0x6f, 0x74, 0x12, 0x2a, 0x0a,
	0x10, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x6f, 0x6f,
	0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x10, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x6f, 0x6f, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x52, 0x65, 0x63,
	0x65, 0x69, 0x70, 0x74, 0x52, 0x6f, 0x6f, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b,
	0x52, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x52, 0x6f, 0x6f, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x06, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x1e, 0x0a, 0x0a, 0x50, 0x72, 0x65,
	0x76, 0x72, 0x61, 0x6e, 0x64, 0x61, 0x6f, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x50,
	0x72, 0x65, 0x76, 0x72, 0x61, 0x6e, 0x64, 0x61, 0x6f, 0x12, 0x1a, 0x0a, 0x08, 0x50, 0x72, 0x6f,
	0x70, 0x6f, 0x73, 0x65, 0x72, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x50, 0x72, 0x6f,
	0x70, 0x6f, 0x73, 0x65, 0x72, 0x12, 0x1c, 0x0a, 0x09, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74,
	0x75, 0x72, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x47, 0x61, 0x73, 0x55, 0x73, 0x65, 0x64, 0x18, 0x0a,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x47, 0x61, 0x73, 0x55, 0x73, 0x65, 0x64, 0x12, 0x1c, 0x0a,
	0x09, 0x45, 0x78, 0x74, 0x72, 0x61, 0x44, 0x61, 0x74, 0x61, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x09, 0x45, 0x78, 0x74, 0x72, 0x61, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1c, 0x0a, 0x09, 0x55,
	0x6e, 0x63, 0x6c, 0x65, 0x48, 0x61, 0x73, 0x68, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09,
	0x55, 0x6e, 0x63, 0x6c, 0x65, 0x48, 0x61, 0x73, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x42, 0x6c, 0x6f,
	0x6f, 0x6d, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x42, 0x6c, 0x6f, 0x6f, 0x6d, 0x22,
	0xf9, 0x02, 0x0a, 0x0d, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x44, 0x61, 0x74,
	0x61, 0x12, 0x16, 0x0a, 0x06, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x06, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x50, 0x61, 0x72,
	0x65, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x50,
	0x61, 0x72, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1c, 0x0a, 0x09, 0x53, 0x74, 0x61,
	0x74, 0x65, 0x52, 0x6f, 0x6f, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x53, 0x74,
	0x61, 0x74, 0x65, 0x52, 0x6f, 0x6f, 0x74, 0x12, 0x2a, 0x0a, 0x10, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x6f, 0x6f, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x10, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52,
	0x6f, 0x6f, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x52, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x52, 0x6f,
	0x6f, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x52, 0x65, 0x63, 0x65, 0x69, 0x70,
	0x74, 0x52, 0x6f, 0x6f, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x18, 0x06, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x12, 0x1e, 0x0a, 0x0a, 0x50, 0x72, 0x65, 0x76, 0x72, 0x61, 0x6e, 0x64, 0x61,
	0x6f, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x50, 0x72, 0x65, 0x76, 0x72, 0x61, 0x6e,
	0x64, 0x61, 0x6f, 0x12, 0x1a, 0x0a, 0x08, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x65, 0x72, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x65, 0x72, 0x12,
	0x18, 0x0a, 0x07, 0x47, 0x61, 0x73, 0x55, 0x73, 0x65, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x07, 0x47, 0x61, 0x73, 0x55, 0x73, 0x65, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x45, 0x78, 0x74,
	0x72, 0x61, 0x44, 0x61, 0x74, 0x61, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x45, 0x78,
	0x74, 0x72, 0x61, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1c, 0x0a, 0x09, 0x55, 0x6e, 0x63, 0x6c, 0x65,
	0x48, 0x61, 0x73, 0x68, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x55, 0x6e, 0x63, 0x6c,
	0x65, 0x48, 0x61, 0x73, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x42, 0x6c, 0x6f, 0x6f, 0x6d, 0x18, 0x0c,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x42, 0x6c, 0x6f, 0x6f, 0x6d, 0x22, 0x8d, 0x01, 0x0a, 0x05,
	0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x2a, 0x0a, 0x06, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x42, 0x6c,
	0x6f, 0x63, 0x6b, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x52, 0x06, 0x48, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x12, 0x2c, 0x0a, 0x04, 0x54, 0x78, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x18, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x54, 0x72,
	0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x04, 0x54, 0x78, 0x6e, 0x73, 0x12,
	0x2a, 0x0a, 0x06, 0x55, 0x6e, 0x63, 0x6c, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x12, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x52, 0x06, 0x55, 0x6e, 0x63, 0x6c, 0x65, 0x73, 0x22, 0x60, 0x0a, 0x04, 0x42,
	0x6f, 0x64, 0x79, 0x12, 0x2c, 0x0a, 0x04, 0x54, 0x78, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x18, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x04, 0x54, 0x78, 0x6e,
	0x73, 0x12, 0x2a, 0x0a, 0x06, 0x55, 0x6e, 0x63, 0x6c, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x12, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x52, 0x06, 0x55, 0x6e, 0x63, 0x6c, 0x65, 0x73, 0x42, 0x08, 0x5a,
	0x06, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_block_proto_rawDescOnce sync.Once
	file_block_proto_rawDescData = file_block_proto_rawDesc
)

func file_block_proto_rawDescGZIP() []byte {
	file_block_proto_rawDescOnce.Do(func() {
		file_block_proto_rawDescData = protoimpl.X.CompressGZIP(file_block_proto_rawDescData)
	})
	return file_block_proto_rawDescData
}

var file_block_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_block_proto_goTypes = []interface{}{
	(*BlockHeader)(nil),   // 0: block.BlockHeader
	(*BlockHashData)(nil), // 1: block.BlockHashData
	(*Block)(nil),         // 2: block.Block
	(*Body)(nil),          // 3: block.Body
	(*Transaction)(nil),   // 4: transaction.Transaction
}
var file_block_proto_depIdxs = []int32{
	0, // 0: block.Block.Header:type_name -> block.BlockHeader
	4, // 1: block.Block.Txns:type_name -> transaction.Transaction
	0, // 2: block.Block.Uncles:type_name -> block.BlockHeader
	4, // 3: block.Body.Txns:type_name -> transaction.Transaction
	0, // 4: block.Body.Uncles:type_name -> block.BlockHeader
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_block_proto_init() }
func file_block_proto_init() {
	if File_block_proto != nil {
		return
	}
	file_transaction_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_block_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlockHeader); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_block_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlockHashData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_block_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Block); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_block_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Body); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_block_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_block_proto_goTypes,
		DependencyIndexes: file_block_proto_depIdxs,
		MessageInfos:      file_block_proto_msgTypes,
	}.Build()
	File_block_proto = out.File
	file_block_proto_rawDesc = nil
	file_block_proto_goTypes = nil
	file_block_proto_depIdxs = nil
}
