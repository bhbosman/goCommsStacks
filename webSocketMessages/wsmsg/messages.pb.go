// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.10
// 	protoc        v3.14.0
// source: messages.proto

package wsmsg

import (
	stream "github.com/bhbosman/gocommon/stream"
	goerrors "github.com/bhbosman/goerrors"
	goprotoextra "github.com/bhbosman/goprotoextra"
	proto "google.golang.org/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type WebSocketMessage_OpCodeEnum int32

const (
	WebSocketMessage_OpContinuation WebSocketMessage_OpCodeEnum = 0
	WebSocketMessage_OpText         WebSocketMessage_OpCodeEnum = 1
	WebSocketMessage_OpBinary       WebSocketMessage_OpCodeEnum = 2
	WebSocketMessage_OpClose        WebSocketMessage_OpCodeEnum = 8
	WebSocketMessage_OpPing         WebSocketMessage_OpCodeEnum = 9
	WebSocketMessage_OpPong         WebSocketMessage_OpCodeEnum = 10
	WebSocketMessage_OpStartLoop    WebSocketMessage_OpCodeEnum = 14
	WebSocketMessage_OpEndLoop      WebSocketMessage_OpCodeEnum = 15
)

// Enum value maps for WebSocketMessage_OpCodeEnum.
var (
	WebSocketMessage_OpCodeEnum_name = map[int32]string{
		0:  "OpContinuation",
		1:  "OpText",
		2:  "OpBinary",
		8:  "OpClose",
		9:  "OpPing",
		10: "OpPong",
		14: "OpStartLoop",
		15: "OpEndLoop",
	}
	WebSocketMessage_OpCodeEnum_value = map[string]int32{
		"OpContinuation": 0,
		"OpText":         1,
		"OpBinary":       2,
		"OpClose":        8,
		"OpPing":         9,
		"OpPong":         10,
		"OpStartLoop":    14,
		"OpEndLoop":      15,
	}
)

func (x WebSocketMessage_OpCodeEnum) Enum() *WebSocketMessage_OpCodeEnum {
	p := new(WebSocketMessage_OpCodeEnum)
	*p = x
	return p
}

func (x WebSocketMessage_OpCodeEnum) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (WebSocketMessage_OpCodeEnum) Descriptor() protoreflect.EnumDescriptor {
	return file_messages_proto_enumTypes[0].Descriptor()
}

func (WebSocketMessage_OpCodeEnum) Type() protoreflect.EnumType {
	return &file_messages_proto_enumTypes[0]
}

func (x WebSocketMessage_OpCodeEnum) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use WebSocketMessage_OpCodeEnum.Descriptor instead.
func (WebSocketMessage_OpCodeEnum) EnumDescriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{0, 0}
}

type WebSocketMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OpCode  WebSocketMessage_OpCodeEnum `protobuf:"varint,1,opt,name=OpCode,proto3,enum=golang.example.policy.WebSocketMessage_OpCodeEnum" json:"OpCode,omitempty"`
	Message []byte                      `protobuf:"bytes,2,opt,name=Message,proto3" json:"Message,omitempty"`
}

func (x *WebSocketMessage) Reset() {
	*x = WebSocketMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WebSocketMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WebSocketMessage) ProtoMessage() {}

func (x *WebSocketMessage) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WebSocketMessage.ProtoReflect.Descriptor instead.
func (*WebSocketMessage) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{0}
}

func (x *WebSocketMessage) GetOpCode() WebSocketMessage_OpCodeEnum {
	if x != nil {
		return x.OpCode
	}
	return WebSocketMessage_OpContinuation
}

func (x *WebSocketMessage) GetMessage() []byte {
	if x != nil {
		return x.Message
	}
	return nil
}

func (self *WebSocketMessage) TypeCode() uint32 {
	return WebSocketMessageTypeCode
}

var file_messages_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         50000,
		Name:          "golang.example.policy.non_sensitive",
		Tag:           "varint,50000,opt,name=non_sensitive",
		Filename:      "messages.proto",
	},
}

// Extension fields to descriptorpb.FieldOptions.
var (
	// optional bool non_sensitive = 50000;
	E_NonSensitive = &file_messages_proto_extTypes[0]
)

var File_messages_proto protoreflect.FileDescriptor

var file_messages_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x15, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65,
	0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xf9, 0x01, 0x0a, 0x10, 0x57, 0x65,
	0x62, 0x53, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x4a,
	0x0a, 0x06, 0x4f, 0x70, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x32,
	0x2e, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e,
	0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x57, 0x65, 0x62, 0x53, 0x6f, 0x63, 0x6b, 0x65, 0x74,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x4f, 0x70, 0x43, 0x6f, 0x64, 0x65, 0x45, 0x6e,
	0x75, 0x6d, 0x52, 0x06, 0x4f, 0x70, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x22, 0x7f, 0x0a, 0x0a, 0x4f, 0x70, 0x43, 0x6f, 0x64, 0x65, 0x45, 0x6e,
	0x75, 0x6d, 0x12, 0x12, 0x0a, 0x0e, 0x4f, 0x70, 0x43, 0x6f, 0x6e, 0x74, 0x69, 0x6e, 0x75, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x4f, 0x70, 0x54, 0x65, 0x78, 0x74,
	0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x4f, 0x70, 0x42, 0x69, 0x6e, 0x61, 0x72, 0x79, 0x10, 0x02,
	0x12, 0x0b, 0x0a, 0x07, 0x4f, 0x70, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x10, 0x08, 0x12, 0x0a, 0x0a,
	0x06, 0x4f, 0x70, 0x50, 0x69, 0x6e, 0x67, 0x10, 0x09, 0x12, 0x0a, 0x0a, 0x06, 0x4f, 0x70, 0x50,
	0x6f, 0x6e, 0x67, 0x10, 0x0a, 0x12, 0x0f, 0x0a, 0x0b, 0x4f, 0x70, 0x53, 0x74, 0x61, 0x72, 0x74,
	0x4c, 0x6f, 0x6f, 0x70, 0x10, 0x0e, 0x12, 0x0d, 0x0a, 0x09, 0x4f, 0x70, 0x45, 0x6e, 0x64, 0x4c,
	0x6f, 0x6f, 0x70, 0x10, 0x0f, 0x3a, 0x44, 0x0a, 0x0d, 0x6e, 0x6f, 0x6e, 0x5f, 0x73, 0x65, 0x6e,
	0x73, 0x69, 0x74, 0x69, 0x76, 0x65, 0x12, 0x1d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xd0, 0x86, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x6e,
	0x6f, 0x6e, 0x53, 0x65, 0x6e, 0x73, 0x69, 0x74, 0x69, 0x76, 0x65, 0x42, 0x08, 0x5a, 0x06, 0x2f,
	0x77, 0x73, 0x6d, 0x73, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_messages_proto_rawDescOnce sync.Once
	file_messages_proto_rawDescData = file_messages_proto_rawDesc
)

func file_messages_proto_rawDescGZIP() []byte {
	file_messages_proto_rawDescOnce.Do(func() {
		file_messages_proto_rawDescData = protoimpl.X.CompressGZIP(file_messages_proto_rawDescData)
	})
	return file_messages_proto_rawDescData
}

var file_messages_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_messages_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_messages_proto_goTypes = []interface{}{
	(WebSocketMessage_OpCodeEnum)(0),  // 0: golang.example.policy.WebSocketMessage.OpCodeEnum
	(*WebSocketMessage)(nil),          // 1: golang.example.policy.WebSocketMessage
	(*descriptorpb.FieldOptions)(nil), // 2: google.protobuf.FieldOptions
}
var file_messages_proto_depIdxs = []int32{
	0, // 0: golang.example.policy.WebSocketMessage.OpCode:type_name -> golang.example.policy.WebSocketMessage.OpCodeEnum
	2, // 1: golang.example.policy.non_sensitive:extendee -> google.protobuf.FieldOptions
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	1, // [1:2] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_messages_proto_init() }
func file_messages_proto_init() {
	if File_messages_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_messages_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WebSocketMessage); i {
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
			RawDescriptor: file_messages_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 1,
			NumServices:   0,
		},
		GoTypes:           file_messages_proto_goTypes,
		DependencyIndexes: file_messages_proto_depIdxs,
		EnumInfos:         file_messages_proto_enumTypes,
		MessageInfos:      file_messages_proto_msgTypes,
		ExtensionInfos:    file_messages_proto_extTypes,
	}.Build()
	File_messages_proto = out.File
	file_messages_proto_rawDesc = nil
	file_messages_proto_goTypes = nil
	file_messages_proto_depIdxs = nil
}

// Typecode generated from: "WebSocketMessage"
const WebSocketMessageTypeCode uint32 = 287629772

//true
//true
//false
//false
type WebSocketMessageWrapper struct {
	goprotoextra.BaseMessageWrapper
	Data *WebSocketMessage
}

func (self *WebSocketMessageWrapper) Message() interface{} {
	return self.Data
}

func (self *WebSocketMessageWrapper) messageWrapper() interface{} {
	return self
}

func NewWebSocketMessageWrapper(
	data *WebSocketMessage) *WebSocketMessageWrapper {
	return &WebSocketMessageWrapper{
		BaseMessageWrapper: goprotoextra.NewBaseMessageWrapper(),
		Data:               data,
	}
}

var _ = stream.Register(
	WebSocketMessageTypeCode,
	stream.TypeCodeData{
		CreateMessage: func() proto.Message {
			return &WebSocketMessage{}
		},
		CreateWrapper: func(
			data proto.Message) (goprotoextra.IMessageWrapper, error) {
			if msg, ok := data.(*WebSocketMessage); ok {
				return NewWebSocketMessageWrapper(
					msg), nil
			}
			return nil, goerrors.InvalidParam
		}})
