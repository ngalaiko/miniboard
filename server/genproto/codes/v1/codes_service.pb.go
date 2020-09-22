// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.10.0
// source: codes/v1/codes_service.proto

package codes

import (
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type CreateCodeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Email address.
	Email string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	// Response type.
	ResponseType string `protobuf:"bytes,2,opt,name=response_type,json=responseType,proto3" json:"response_type,omitempty"`
}

func (x *CreateCodeRequest) Reset() {
	*x = CreateCodeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_codes_v1_codes_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateCodeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateCodeRequest) ProtoMessage() {}

func (x *CreateCodeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_codes_v1_codes_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateCodeRequest.ProtoReflect.Descriptor instead.
func (*CreateCodeRequest) Descriptor() ([]byte, []int) {
	return file_codes_v1_codes_service_proto_rawDescGZIP(), []int{0}
}

func (x *CreateCodeRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *CreateCodeRequest) GetResponseType() string {
	if x != nil {
		return x.ResponseType
	}
	return ""
}

type Code struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Code) Reset() {
	*x = Code{}
	if protoimpl.UnsafeEnabled {
		mi := &file_codes_v1_codes_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Code) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Code) ProtoMessage() {}

func (x *Code) ProtoReflect() protoreflect.Message {
	mi := &file_codes_v1_codes_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Code.ProtoReflect.Descriptor instead.
func (*Code) Descriptor() ([]byte, []int) {
	return file_codes_v1_codes_service_proto_rawDescGZIP(), []int{1}
}

var File_codes_v1_codes_service_proto protoreflect.FileDescriptor

var file_codes_v1_codes_service_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x63, 0x6f, 0x64, 0x65, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x64, 0x65, 0x73,
	0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16,
	0x61, 0x70, 0x70, 0x2e, 0x6d, 0x69, 0x6e, 0x69, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x63, 0x6f,
	0x64, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4e, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6f,
	0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61,
	0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12,
	0x23, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x54, 0x79, 0x70, 0x65, 0x22, 0x06, 0x0a, 0x04, 0x43, 0x6f, 0x64, 0x65, 0x32, 0x7f, 0x0a, 0x0c,
	0x43, 0x6f, 0x64, 0x65, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x6f, 0x0a, 0x0a,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x29, 0x2e, 0x61, 0x70, 0x70,
	0x2e, 0x6d, 0x69, 0x6e, 0x69, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x63, 0x6f, 0x64, 0x65, 0x73,
	0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x6d, 0x69, 0x6e, 0x69,
	0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x63, 0x6f, 0x64, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x43,
	0x6f, 0x64, 0x65, 0x22, 0x18, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x12, 0x22, 0x0d, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x64, 0x65, 0x73, 0x3a, 0x01, 0x2a, 0x42, 0x3e, 0x5a,
	0x3c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x67, 0x61, 0x6c,
	0x61, 0x69, 0x6b, 0x6f, 0x2f, 0x6d, 0x69, 0x6e, 0x69, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2f, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x67, 0x65, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63,
	0x6f, 0x64, 0x65, 0x73, 0x2f, 0x76, 0x31, 0x3b, 0x63, 0x6f, 0x64, 0x65, 0x73, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_codes_v1_codes_service_proto_rawDescOnce sync.Once
	file_codes_v1_codes_service_proto_rawDescData = file_codes_v1_codes_service_proto_rawDesc
)

func file_codes_v1_codes_service_proto_rawDescGZIP() []byte {
	file_codes_v1_codes_service_proto_rawDescOnce.Do(func() {
		file_codes_v1_codes_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_codes_v1_codes_service_proto_rawDescData)
	})
	return file_codes_v1_codes_service_proto_rawDescData
}

var file_codes_v1_codes_service_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_codes_v1_codes_service_proto_goTypes = []interface{}{
	(*CreateCodeRequest)(nil), // 0: app.miniboard.codes.v1.CreateCodeRequest
	(*Code)(nil),              // 1: app.miniboard.codes.v1.Code
}
var file_codes_v1_codes_service_proto_depIdxs = []int32{
	0, // 0: app.miniboard.codes.v1.CodesService.CreateCode:input_type -> app.miniboard.codes.v1.CreateCodeRequest
	1, // 1: app.miniboard.codes.v1.CodesService.CreateCode:output_type -> app.miniboard.codes.v1.Code
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_codes_v1_codes_service_proto_init() }
func file_codes_v1_codes_service_proto_init() {
	if File_codes_v1_codes_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_codes_v1_codes_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateCodeRequest); i {
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
		file_codes_v1_codes_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Code); i {
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
			RawDescriptor: file_codes_v1_codes_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_codes_v1_codes_service_proto_goTypes,
		DependencyIndexes: file_codes_v1_codes_service_proto_depIdxs,
		MessageInfos:      file_codes_v1_codes_service_proto_msgTypes,
	}.Build()
	File_codes_v1_codes_service_proto = out.File
	file_codes_v1_codes_service_proto_rawDesc = nil
	file_codes_v1_codes_service_proto_goTypes = nil
	file_codes_v1_codes_service_proto_depIdxs = nil
}
