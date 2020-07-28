// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.24.0
// 	protoc        v3.11.4
// source: pkg/server/pb/fetcher.proto

package server

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Measure struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID       int32  `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	URL      string `protobuf:"bytes,2,opt,name=URL,proto3" json:"URL,omitempty"`
	Interval int32  `protobuf:"varint,3,opt,name=interval,proto3" json:"interval,omitempty"`
}

func (x *Measure) Reset() {
	*x = Measure{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_server_pb_fetcher_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Measure) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Measure) ProtoMessage() {}

func (x *Measure) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_server_pb_fetcher_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Measure.ProtoReflect.Descriptor instead.
func (*Measure) Descriptor() ([]byte, []int) {
	return file_pkg_server_pb_fetcher_proto_rawDescGZIP(), []int{0}
}

func (x *Measure) GetID() int32 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *Measure) GetURL() string {
	if x != nil {
		return x.URL
	}
	return ""
}

func (x *Measure) GetInterval() int32 {
	if x != nil {
		return x.Interval
	}
	return 0
}

type GetMeasuresRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetMeasuresRequest) Reset() {
	*x = GetMeasuresRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_server_pb_fetcher_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMeasuresRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMeasuresRequest) ProtoMessage() {}

func (x *GetMeasuresRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_server_pb_fetcher_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMeasuresRequest.ProtoReflect.Descriptor instead.
func (*GetMeasuresRequest) Descriptor() ([]byte, []int) {
	return file_pkg_server_pb_fetcher_proto_rawDescGZIP(), []int{1}
}

type GetMeasuresResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Measures []*Measure `protobuf:"bytes,1,rep,name=measures,proto3" json:"measures,omitempty"`
}

func (x *GetMeasuresResponse) Reset() {
	*x = GetMeasuresResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_server_pb_fetcher_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMeasuresResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMeasuresResponse) ProtoMessage() {}

func (x *GetMeasuresResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_server_pb_fetcher_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMeasuresResponse.ProtoReflect.Descriptor instead.
func (*GetMeasuresResponse) Descriptor() ([]byte, []int) {
	return file_pkg_server_pb_fetcher_proto_rawDescGZIP(), []int{2}
}

func (x *GetMeasuresResponse) GetMeasures() []*Measure {
	if x != nil {
		return x.Measures
	}
	return nil
}

var File_pkg_server_pb_fetcher_proto protoreflect.FileDescriptor

var file_pkg_server_pb_fetcher_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x70, 0x6b, 0x67, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x70, 0x62, 0x2f,
	0x66, 0x65, 0x74, 0x63, 0x68, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x22, 0x47, 0x0a, 0x07, 0x4d, 0x65, 0x61, 0x73, 0x75, 0x72, 0x65,
	0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x49, 0x44,
	0x12, 0x10, 0x0a, 0x03, 0x55, 0x52, 0x4c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x55,
	0x52, 0x4c, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x22, 0x14,
	0x0a, 0x12, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x61, 0x73, 0x75, 0x72, 0x65, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x22, 0x42, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x61, 0x73, 0x75,
	0x72, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2b, 0x0a, 0x08, 0x6d,
	0x65, 0x61, 0x73, 0x75, 0x72, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x4d, 0x65, 0x61, 0x73, 0x75, 0x72, 0x65, 0x52, 0x08,
	0x6d, 0x65, 0x61, 0x73, 0x75, 0x72, 0x65, 0x73, 0x32, 0x58, 0x0a, 0x0e, 0x46, 0x65, 0x74, 0x63,
	0x68, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x46, 0x0a, 0x0b, 0x47, 0x65,
	0x74, 0x4d, 0x65, 0x61, 0x73, 0x75, 0x72, 0x65, 0x73, 0x12, 0x1a, 0x2e, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x61, 0x73, 0x75, 0x72, 0x65, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x47,
	0x65, 0x74, 0x4d, 0x65, 0x61, 0x73, 0x75, 0x72, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x42, 0x0c, 0x5a, 0x0a, 0x70, 0x6b, 0x67, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_server_pb_fetcher_proto_rawDescOnce sync.Once
	file_pkg_server_pb_fetcher_proto_rawDescData = file_pkg_server_pb_fetcher_proto_rawDesc
)

func file_pkg_server_pb_fetcher_proto_rawDescGZIP() []byte {
	file_pkg_server_pb_fetcher_proto_rawDescOnce.Do(func() {
		file_pkg_server_pb_fetcher_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_server_pb_fetcher_proto_rawDescData)
	})
	return file_pkg_server_pb_fetcher_proto_rawDescData
}

var file_pkg_server_pb_fetcher_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_pkg_server_pb_fetcher_proto_goTypes = []interface{}{
	(*Measure)(nil),             // 0: server.Measure
	(*GetMeasuresRequest)(nil),  // 1: server.GetMeasuresRequest
	(*GetMeasuresResponse)(nil), // 2: server.GetMeasuresResponse
}
var file_pkg_server_pb_fetcher_proto_depIdxs = []int32{
	0, // 0: server.GetMeasuresResponse.measures:type_name -> server.Measure
	1, // 1: server.FetcherService.GetMeasures:input_type -> server.GetMeasuresRequest
	2, // 2: server.FetcherService.GetMeasures:output_type -> server.GetMeasuresResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_pkg_server_pb_fetcher_proto_init() }
func file_pkg_server_pb_fetcher_proto_init() {
	if File_pkg_server_pb_fetcher_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_server_pb_fetcher_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Measure); i {
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
		file_pkg_server_pb_fetcher_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetMeasuresRequest); i {
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
		file_pkg_server_pb_fetcher_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetMeasuresResponse); i {
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
			RawDescriptor: file_pkg_server_pb_fetcher_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_server_pb_fetcher_proto_goTypes,
		DependencyIndexes: file_pkg_server_pb_fetcher_proto_depIdxs,
		MessageInfos:      file_pkg_server_pb_fetcher_proto_msgTypes,
	}.Build()
	File_pkg_server_pb_fetcher_proto = out.File
	file_pkg_server_pb_fetcher_proto_rawDesc = nil
	file_pkg_server_pb_fetcher_proto_goTypes = nil
	file_pkg_server_pb_fetcher_proto_depIdxs = nil
}
