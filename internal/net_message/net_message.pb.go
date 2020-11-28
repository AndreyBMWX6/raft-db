// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: net_message.proto

package net_message

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

type BaseRaftMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OwnerIp   []byte `protobuf:"bytes,1,opt,name=OwnerIp,proto3" json:"OwnerIp,omitempty"`
	OwnerPort uint32 `protobuf:"varint,2,opt,name=OwnerPort,proto3" json:"OwnerPort,omitempty"`
	DestIp    []byte `protobuf:"bytes,3,opt,name=DestIp,proto3" json:"DestIp,omitempty"`
	DestPort  uint32 `protobuf:"varint,4,opt,name=DestPort,proto3" json:"DestPort,omitempty"`
	CurrTerm  uint32 `protobuf:"varint,5,opt,name=CurrTerm,proto3" json:"CurrTerm,omitempty"`
}

func (x *BaseRaftMessage) Reset() {
	*x = BaseRaftMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_net_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BaseRaftMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BaseRaftMessage) ProtoMessage() {}

func (x *BaseRaftMessage) ProtoReflect() protoreflect.Message {
	mi := &file_net_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BaseRaftMessage.ProtoReflect.Descriptor instead.
func (*BaseRaftMessage) Descriptor() ([]byte, []int) {
	return file_net_message_proto_rawDescGZIP(), []int{0}
}

func (x *BaseRaftMessage) GetOwnerIp() []byte {
	if x != nil {
		return x.OwnerIp
	}
	return nil
}

func (x *BaseRaftMessage) GetOwnerPort() uint32 {
	if x != nil {
		return x.OwnerPort
	}
	return 0
}

func (x *BaseRaftMessage) GetDestIp() []byte {
	if x != nil {
		return x.DestIp
	}
	return nil
}

func (x *BaseRaftMessage) GetDestPort() uint32 {
	if x != nil {
		return x.DestPort
	}
	return 0
}

func (x *BaseRaftMessage) GetCurrTerm() uint32 {
	if x != nil {
		return x.CurrTerm
	}
	return 0
}

type Entry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term  uint32 `protobuf:"varint,1,opt,name=Term,proto3" json:"Term,omitempty"`
	Query []byte `protobuf:"bytes,2,opt,name=Query,proto3" json:"Query,omitempty"`
}

func (x *Entry) Reset() {
	*x = Entry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_net_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Entry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Entry) ProtoMessage() {}

func (x *Entry) ProtoReflect() protoreflect.Message {
	mi := &file_net_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Entry.ProtoReflect.Descriptor instead.
func (*Entry) Descriptor() ([]byte, []int) {
	return file_net_message_proto_rawDescGZIP(), []int{1}
}

func (x *Entry) GetTerm() uint32 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *Entry) GetQuery() []byte {
	if x != nil {
		return x.Query
	}
	return nil
}

type AppendEntries struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg      *BaseRaftMessage `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	PrevTerm uint32           `protobuf:"varint,2,opt,name=PrevTerm,proto3" json:"PrevTerm,omitempty"`
	NewIndex uint32           `protobuf:"varint,3,opt,name=NewIndex,proto3" json:"NewIndex,omitempty"`
	Entries  []*Entry         `protobuf:"bytes,4,rep,name=Entries,proto3" json:"Entries,omitempty"`
}

func (x *AppendEntries) Reset() {
	*x = AppendEntries{}
	if protoimpl.UnsafeEnabled {
		mi := &file_net_message_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppendEntries) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppendEntries) ProtoMessage() {}

func (x *AppendEntries) ProtoReflect() protoreflect.Message {
	mi := &file_net_message_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppendEntries.ProtoReflect.Descriptor instead.
func (*AppendEntries) Descriptor() ([]byte, []int) {
	return file_net_message_proto_rawDescGZIP(), []int{2}
}

func (x *AppendEntries) GetMsg() *BaseRaftMessage {
	if x != nil {
		return x.Msg
	}
	return nil
}

func (x *AppendEntries) GetPrevTerm() uint32 {
	if x != nil {
		return x.PrevTerm
	}
	return 0
}

func (x *AppendEntries) GetNewIndex() uint32 {
	if x != nil {
		return x.NewIndex
	}
	return 0
}

func (x *AppendEntries) GetEntries() []*Entry {
	if x != nil {
		return x.Entries
	}
	return nil
}

type AppendAck struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg       *BaseRaftMessage `protobuf:"bytes,1,opt,name=Msg,proto3" json:"Msg,omitempty"`
	Appended  bool             `protobuf:"varint,2,opt,name=Appended,proto3" json:"Appended,omitempty"`
	Heartbeat bool             `protobuf:"varint,3,opt,name=Heartbeat,proto3" json:"Heartbeat,omitempty"`
}

func (x *AppendAck) Reset() {
	*x = AppendAck{}
	if protoimpl.UnsafeEnabled {
		mi := &file_net_message_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppendAck) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppendAck) ProtoMessage() {}

func (x *AppendAck) ProtoReflect() protoreflect.Message {
	mi := &file_net_message_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppendAck.ProtoReflect.Descriptor instead.
func (*AppendAck) Descriptor() ([]byte, []int) {
	return file_net_message_proto_rawDescGZIP(), []int{3}
}

func (x *AppendAck) GetMsg() *BaseRaftMessage {
	if x != nil {
		return x.Msg
	}
	return nil
}

func (x *AppendAck) GetAppended() bool {
	if x != nil {
		return x.Appended
	}
	return false
}

func (x *AppendAck) GetHeartbeat() bool {
	if x != nil {
		return x.Heartbeat
	}
	return false
}

type RequestVote struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg      *BaseRaftMessage `protobuf:"bytes,1,opt,name=Msg,proto3" json:"Msg,omitempty"`
	TopIndex uint32           `protobuf:"varint,2,opt,name=TopIndex,proto3" json:"TopIndex,omitempty"`
	TopTerm  uint32           `protobuf:"varint,3,opt,name=TopTerm,proto3" json:"TopTerm,omitempty"`
}

func (x *RequestVote) Reset() {
	*x = RequestVote{}
	if protoimpl.UnsafeEnabled {
		mi := &file_net_message_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestVote) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestVote) ProtoMessage() {}

func (x *RequestVote) ProtoReflect() protoreflect.Message {
	mi := &file_net_message_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestVote.ProtoReflect.Descriptor instead.
func (*RequestVote) Descriptor() ([]byte, []int) {
	return file_net_message_proto_rawDescGZIP(), []int{4}
}

func (x *RequestVote) GetMsg() *BaseRaftMessage {
	if x != nil {
		return x.Msg
	}
	return nil
}

func (x *RequestVote) GetTopIndex() uint32 {
	if x != nil {
		return x.TopIndex
	}
	return 0
}

func (x *RequestVote) GetTopTerm() uint32 {
	if x != nil {
		return x.TopTerm
	}
	return 0
}

type RequestAck struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg   *BaseRaftMessage `protobuf:"bytes,1,opt,name=Msg,proto3" json:"Msg,omitempty"`
	Voted bool             `protobuf:"varint,2,opt,name=Voted,proto3" json:"Voted,omitempty"`
}

func (x *RequestAck) Reset() {
	*x = RequestAck{}
	if protoimpl.UnsafeEnabled {
		mi := &file_net_message_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestAck) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestAck) ProtoMessage() {}

func (x *RequestAck) ProtoReflect() protoreflect.Message {
	mi := &file_net_message_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestAck.ProtoReflect.Descriptor instead.
func (*RequestAck) Descriptor() ([]byte, []int) {
	return file_net_message_proto_rawDescGZIP(), []int{5}
}

func (x *RequestAck) GetMsg() *BaseRaftMessage {
	if x != nil {
		return x.Msg
	}
	return nil
}

func (x *RequestAck) GetVoted() bool {
	if x != nil {
		return x.Voted
	}
	return false
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to RaftMessage:
	//	*Message_AppendEntries
	//	*Message_AppendAck
	//	*Message_RequestVote
	//	*Message_RequestAck
	RaftMessage isMessage_RaftMessage `protobuf_oneof:"RaftMessage"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_net_message_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_net_message_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_net_message_proto_rawDescGZIP(), []int{6}
}

func (m *Message) GetRaftMessage() isMessage_RaftMessage {
	if m != nil {
		return m.RaftMessage
	}
	return nil
}

func (x *Message) GetAppendEntries() *AppendEntries {
	if x, ok := x.GetRaftMessage().(*Message_AppendEntries); ok {
		return x.AppendEntries
	}
	return nil
}

func (x *Message) GetAppendAck() *AppendAck {
	if x, ok := x.GetRaftMessage().(*Message_AppendAck); ok {
		return x.AppendAck
	}
	return nil
}

func (x *Message) GetRequestVote() *RequestVote {
	if x, ok := x.GetRaftMessage().(*Message_RequestVote); ok {
		return x.RequestVote
	}
	return nil
}

func (x *Message) GetRequestAck() *RequestAck {
	if x, ok := x.GetRaftMessage().(*Message_RequestAck); ok {
		return x.RequestAck
	}
	return nil
}

type isMessage_RaftMessage interface {
	isMessage_RaftMessage()
}

type Message_AppendEntries struct {
	AppendEntries *AppendEntries `protobuf:"bytes,1,opt,name=AppendEntries,proto3,oneof"`
}

type Message_AppendAck struct {
	AppendAck *AppendAck `protobuf:"bytes,2,opt,name=AppendAck,proto3,oneof"`
}

type Message_RequestVote struct {
	RequestVote *RequestVote `protobuf:"bytes,3,opt,name=RequestVote,proto3,oneof"`
}

type Message_RequestAck struct {
	RequestAck *RequestAck `protobuf:"bytes,4,opt,name=RequestAck,proto3,oneof"`
}

func (*Message_AppendEntries) isMessage_RaftMessage() {}

func (*Message_AppendAck) isMessage_RaftMessage() {}

func (*Message_RequestVote) isMessage_RaftMessage() {}

func (*Message_RequestAck) isMessage_RaftMessage() {}

var File_net_message_proto protoreflect.FileDescriptor

var file_net_message_proto_rawDesc = []byte{
	0x0a, 0x11, 0x6e, 0x65, 0x74, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x6e, 0x65, 0x74, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x22, 0x99, 0x01, 0x0a, 0x0f, 0x42, 0x61, 0x73, 0x65, 0x52, 0x61, 0x66, 0x74, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x49, 0x70, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x49, 0x70, 0x12, 0x1c,
	0x0a, 0x09, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x50, 0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x09, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x16, 0x0a, 0x06,
	0x44, 0x65, 0x73, 0x74, 0x49, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x44, 0x65,
	0x73, 0x74, 0x49, 0x70, 0x12, 0x1a, 0x0a, 0x08, 0x44, 0x65, 0x73, 0x74, 0x50, 0x6f, 0x72, 0x74,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x44, 0x65, 0x73, 0x74, 0x50, 0x6f, 0x72, 0x74,
	0x12, 0x1a, 0x0a, 0x08, 0x43, 0x75, 0x72, 0x72, 0x54, 0x65, 0x72, 0x6d, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x08, 0x43, 0x75, 0x72, 0x72, 0x54, 0x65, 0x72, 0x6d, 0x22, 0x31, 0x0a, 0x05,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x04, 0x54, 0x65, 0x72, 0x6d, 0x12, 0x14, 0x0a, 0x05, 0x51, 0x75, 0x65,
	0x72, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x51, 0x75, 0x65, 0x72, 0x79, 0x22,
	0xa5, 0x01, 0x0a, 0x0d, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65,
	0x73, 0x12, 0x2e, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c,
	0x2e, 0x6e, 0x65, 0x74, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x42, 0x61, 0x73,
	0x65, 0x52, 0x61, 0x66, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x03, 0x6d, 0x73,
	0x67, 0x12, 0x1a, 0x0a, 0x08, 0x50, 0x72, 0x65, 0x76, 0x54, 0x65, 0x72, 0x6d, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x08, 0x50, 0x72, 0x65, 0x76, 0x54, 0x65, 0x72, 0x6d, 0x12, 0x1a, 0x0a,
	0x08, 0x4e, 0x65, 0x77, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x08, 0x4e, 0x65, 0x77, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x2c, 0x0a, 0x07, 0x45, 0x6e, 0x74,
	0x72, 0x69, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6e, 0x65, 0x74,
	0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07,
	0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x22, 0x75, 0x0a, 0x09, 0x41, 0x70, 0x70, 0x65, 0x6e,
	0x64, 0x41, 0x63, 0x6b, 0x12, 0x2e, 0x0a, 0x03, 0x4d, 0x73, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1c, 0x2e, 0x6e, 0x65, 0x74, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e,
	0x42, 0x61, 0x73, 0x65, 0x52, 0x61, 0x66, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52,
	0x03, 0x4d, 0x73, 0x67, 0x12, 0x1a, 0x0a, 0x08, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x65, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x65, 0x64,
	0x12, 0x1c, 0x0a, 0x09, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x09, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x22, 0x73,
	0x0a, 0x0b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x12, 0x2e, 0x0a,
	0x03, 0x4d, 0x73, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x6e, 0x65, 0x74,
	0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x42, 0x61, 0x73, 0x65, 0x52, 0x61, 0x66,
	0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x03, 0x4d, 0x73, 0x67, 0x12, 0x1a, 0x0a,
	0x08, 0x54, 0x6f, 0x70, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x08, 0x54, 0x6f, 0x70, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x18, 0x0a, 0x07, 0x54, 0x6f, 0x70,
	0x54, 0x65, 0x72, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x54, 0x6f, 0x70, 0x54,
	0x65, 0x72, 0x6d, 0x22, 0x52, 0x0a, 0x0a, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x41, 0x63,
	0x6b, 0x12, 0x2e, 0x0a, 0x03, 0x4d, 0x73, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c,
	0x2e, 0x6e, 0x65, 0x74, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x42, 0x61, 0x73,
	0x65, 0x52, 0x61, 0x66, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x03, 0x4d, 0x73,
	0x67, 0x12, 0x14, 0x0a, 0x05, 0x56, 0x6f, 0x74, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x05, 0x56, 0x6f, 0x74, 0x65, 0x64, 0x22, 0x8d, 0x02, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x42, 0x0a, 0x0d, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74,
	0x72, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x6e, 0x65, 0x74,
	0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45,
	0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x48, 0x00, 0x52, 0x0d, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64,
	0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x12, 0x36, 0x0a, 0x09, 0x41, 0x70, 0x70, 0x65, 0x6e,
	0x64, 0x41, 0x63, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x6e, 0x65, 0x74,
	0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x41,
	0x63, 0x6b, 0x48, 0x00, 0x52, 0x09, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x41, 0x63, 0x6b, 0x12,
	0x3c, 0x0a, 0x0b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x6e, 0x65, 0x74, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x48, 0x00,
	0x52, 0x0b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x12, 0x39, 0x0a,
	0x0a, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x41, 0x63, 0x6b, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x17, 0x2e, 0x6e, 0x65, 0x74, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x41, 0x63, 0x6b, 0x48, 0x00, 0x52, 0x0a, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x41, 0x63, 0x6b, 0x42, 0x0d, 0x0a, 0x0b, 0x52, 0x61, 0x66, 0x74,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_net_message_proto_rawDescOnce sync.Once
	file_net_message_proto_rawDescData = file_net_message_proto_rawDesc
)

func file_net_message_proto_rawDescGZIP() []byte {
	file_net_message_proto_rawDescOnce.Do(func() {
		file_net_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_net_message_proto_rawDescData)
	})
	return file_net_message_proto_rawDescData
}

var file_net_message_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_net_message_proto_goTypes = []interface{}{
	(*BaseRaftMessage)(nil), // 0: net_message.BaseRaftMessage
	(*Entry)(nil),           // 1: net_message.Entry
	(*AppendEntries)(nil),   // 2: net_message.AppendEntries
	(*AppendAck)(nil),       // 3: net_message.AppendAck
	(*RequestVote)(nil),     // 4: net_message.RequestVote
	(*RequestAck)(nil),      // 5: net_message.RequestAck
	(*Message)(nil),         // 6: net_message.Message
}
var file_net_message_proto_depIdxs = []int32{
	0, // 0: net_message.AppendEntries.msg:type_name -> net_message.BaseRaftMessage
	1, // 1: net_message.AppendEntries.Entries:type_name -> net_message.Entry
	0, // 2: net_message.AppendAck.Msg:type_name -> net_message.BaseRaftMessage
	0, // 3: net_message.RequestVote.Msg:type_name -> net_message.BaseRaftMessage
	0, // 4: net_message.RequestAck.Msg:type_name -> net_message.BaseRaftMessage
	2, // 5: net_message.Message.AppendEntries:type_name -> net_message.AppendEntries
	3, // 6: net_message.Message.AppendAck:type_name -> net_message.AppendAck
	4, // 7: net_message.Message.RequestVote:type_name -> net_message.RequestVote
	5, // 8: net_message.Message.RequestAck:type_name -> net_message.RequestAck
	9, // [9:9] is the sub-list for method output_type
	9, // [9:9] is the sub-list for method input_type
	9, // [9:9] is the sub-list for extension type_name
	9, // [9:9] is the sub-list for extension extendee
	0, // [0:9] is the sub-list for field type_name
}

func init() { file_net_message_proto_init() }
func file_net_message_proto_init() {
	if File_net_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_net_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BaseRaftMessage); i {
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
		file_net_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Entry); i {
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
		file_net_message_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppendEntries); i {
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
		file_net_message_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppendAck); i {
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
		file_net_message_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestVote); i {
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
		file_net_message_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestAck); i {
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
		file_net_message_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
	file_net_message_proto_msgTypes[6].OneofWrappers = []interface{}{
		(*Message_AppendEntries)(nil),
		(*Message_AppendAck)(nil),
		(*Message_RequestVote)(nil),
		(*Message_RequestAck)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_net_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_net_message_proto_goTypes,
		DependencyIndexes: file_net_message_proto_depIdxs,
		MessageInfos:      file_net_message_proto_msgTypes,
	}.Build()
	File_net_message_proto = out.File
	file_net_message_proto_rawDesc = nil
	file_net_message_proto_goTypes = nil
	file_net_message_proto_depIdxs = nil
}