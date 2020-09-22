// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.10.0
// source: sources/v1/sources_service.proto

package sources

import (
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	longrunning "google.golang.org/genproto/googleapis/longrunning"
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

type CreateSourceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Source to create.
	Source *Source `protobuf:"bytes,1,opt,name=source,proto3" json:"source,omitempty"`
}

func (x *CreateSourceRequest) Reset() {
	*x = CreateSourceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sources_v1_sources_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateSourceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateSourceRequest) ProtoMessage() {}

func (x *CreateSourceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_sources_v1_sources_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateSourceRequest.ProtoReflect.Descriptor instead.
func (*CreateSourceRequest) Descriptor() ([]byte, []int) {
	return file_sources_v1_sources_service_proto_rawDescGZIP(), []int{0}
}

func (x *CreateSourceRequest) GetSource() *Source {
	if x != nil {
		return x.Source
	}
	return nil
}

type SourceList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sources []*Source `protobuf:"bytes,1,rep,name=sources,proto3" json:"sources,omitempty"`
}

func (x *SourceList) Reset() {
	*x = SourceList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sources_v1_sources_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SourceList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SourceList) ProtoMessage() {}

func (x *SourceList) ProtoReflect() protoreflect.Message {
	mi := &file_sources_v1_sources_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SourceList.ProtoReflect.Descriptor instead.
func (*SourceList) Descriptor() ([]byte, []int) {
	return file_sources_v1_sources_service_proto_rawDescGZIP(), []int{1}
}

func (x *SourceList) GetSources() []*Source {
	if x != nil {
		return x.Sources
	}
	return nil
}

type Source struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Source url.
	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	// Raw is a base64 encoded raw source data.
	Raw []byte `protobuf:"bytes,2,opt,name=raw,proto3" json:"raw,omitempty"`
}

func (x *Source) Reset() {
	*x = Source{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sources_v1_sources_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Source) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Source) ProtoMessage() {}

func (x *Source) ProtoReflect() protoreflect.Message {
	mi := &file_sources_v1_sources_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Source.ProtoReflect.Descriptor instead.
func (*Source) Descriptor() ([]byte, []int) {
	return file_sources_v1_sources_service_proto_rawDescGZIP(), []int{2}
}

func (x *Source) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *Source) GetRaw() []byte {
	if x != nil {
		return x.Raw
	}
	return nil
}

var File_sources_v1_sources_service_proto protoreflect.FileDescriptor

var file_sources_v1_sources_service_proto_rawDesc = []byte{
	0x0a, 0x20, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x1e, 0x61, 0x70, 0x70, 0x2e, 0x6d, 0x69, 0x6e, 0x69, 0x62, 0x6f, 0x61, 0x72,
	0x64, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e,
	0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61,
	0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x23, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x6c, 0x6f, 0x6e, 0x67, 0x72, 0x75, 0x6e,
	0x6e, 0x69, 0x6e, 0x67, 0x2f, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x55, 0x0a, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x3e, 0x0a, 0x06,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x61,
	0x70, 0x70, 0x2e, 0x6d, 0x69, 0x6e, 0x69, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x75, 0x73, 0x65,
	0x72, 0x73, 0x2e, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x22, 0x4e, 0x0a, 0x0a,
	0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x40, 0x0a, 0x07, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x61, 0x70,
	0x70, 0x2e, 0x6d, 0x69, 0x6e, 0x69, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x75, 0x73, 0x65, 0x72,
	0x73, 0x2e, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x52, 0x07, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x22, 0x2c, 0x0a, 0x06,
	0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x10, 0x0a, 0x03, 0x72, 0x61, 0x77, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x72, 0x61, 0x77, 0x32, 0x96, 0x01, 0x0a, 0x0e, 0x53,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x83, 0x01,
	0x0a, 0x0c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x33,
	0x2e, 0x61, 0x70, 0x70, 0x2e, 0x6d, 0x69, 0x6e, 0x69, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x75,
	0x73, 0x65, 0x72, 0x73, 0x2e, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x6c, 0x6f, 0x6e,
	0x67, 0x72, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x2e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x22, 0x1f, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x19, 0x22, 0x0f, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x76, 0x31, 0x2f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x3a, 0x06, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x42, 0x42, 0x5a, 0x40, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x6e, 0x67, 0x61, 0x6c, 0x61, 0x69, 0x6b, 0x6f, 0x2f, 0x6d, 0x69, 0x6e, 0x69, 0x62,
	0x6f, 0x61, 0x72, 0x64, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x67, 0x65, 0x6e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x76, 0x31, 0x3b,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_sources_v1_sources_service_proto_rawDescOnce sync.Once
	file_sources_v1_sources_service_proto_rawDescData = file_sources_v1_sources_service_proto_rawDesc
)

func file_sources_v1_sources_service_proto_rawDescGZIP() []byte {
	file_sources_v1_sources_service_proto_rawDescOnce.Do(func() {
		file_sources_v1_sources_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_sources_v1_sources_service_proto_rawDescData)
	})
	return file_sources_v1_sources_service_proto_rawDescData
}

var file_sources_v1_sources_service_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_sources_v1_sources_service_proto_goTypes = []interface{}{
	(*CreateSourceRequest)(nil),   // 0: app.miniboard.users.sources.v1.CreateSourceRequest
	(*SourceList)(nil),            // 1: app.miniboard.users.sources.v1.SourceList
	(*Source)(nil),                // 2: app.miniboard.users.sources.v1.Source
	(*longrunning.Operation)(nil), // 3: google.longrunning.Operation
}
var file_sources_v1_sources_service_proto_depIdxs = []int32{
	2, // 0: app.miniboard.users.sources.v1.CreateSourceRequest.source:type_name -> app.miniboard.users.sources.v1.Source
	2, // 1: app.miniboard.users.sources.v1.SourceList.sources:type_name -> app.miniboard.users.sources.v1.Source
	0, // 2: app.miniboard.users.sources.v1.SourcesService.CreateSource:input_type -> app.miniboard.users.sources.v1.CreateSourceRequest
	3, // 3: app.miniboard.users.sources.v1.SourcesService.CreateSource:output_type -> google.longrunning.Operation
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_sources_v1_sources_service_proto_init() }
func file_sources_v1_sources_service_proto_init() {
	if File_sources_v1_sources_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_sources_v1_sources_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateSourceRequest); i {
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
		file_sources_v1_sources_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SourceList); i {
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
		file_sources_v1_sources_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Source); i {
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
			RawDescriptor: file_sources_v1_sources_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_sources_v1_sources_service_proto_goTypes,
		DependencyIndexes: file_sources_v1_sources_service_proto_depIdxs,
		MessageInfos:      file_sources_v1_sources_service_proto_msgTypes,
	}.Build()
	File_sources_v1_sources_service_proto = out.File
	file_sources_v1_sources_service_proto_rawDesc = nil
	file_sources_v1_sources_service_proto_goTypes = nil
	file_sources_v1_sources_service_proto_depIdxs = nil
}
