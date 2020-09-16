// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.10.0
// source: tokens_service.proto

package tokens

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type CreateTokenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code string `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
}

func (x *CreateTokenRequest) Reset() {
	*x = CreateTokenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tokens_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTokenRequest) ProtoMessage() {}

func (x *CreateTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_tokens_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateTokenRequest.ProtoReflect.Descriptor instead.
func (*CreateTokenRequest) Descriptor() ([]byte, []int) {
	return file_tokens_service_proto_rawDescGZIP(), []int{0}
}

func (x *CreateTokenRequest) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

type Token struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *Token) Reset() {
	*x = Token{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tokens_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Token) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Token) ProtoMessage() {}

func (x *Token) ProtoReflect() protoreflect.Message {
	mi := &file_tokens_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Token.ProtoReflect.Descriptor instead.
func (*Token) Descriptor() ([]byte, []int) {
	return file_tokens_service_proto_rawDescGZIP(), []int{1}
}

func (x *Token) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

var File_tokens_service_proto protoreflect.FileDescriptor

var file_tokens_service_proto_rawDesc = []byte{
	0x0a, 0x14, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x17, 0x61, 0x70, 0x70, 0x2e, 0x6d, 0x69, 0x6e, 0x69,
	0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x2e, 0x76, 0x31, 0x1a,
	0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x28, 0x0a,
	0x12, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x22, 0x1d, 0x0a, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x32, 0x86, 0x01, 0x0a, 0x0d, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x75, 0x0a, 0x0b, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x2b, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x6d, 0x69,
	0x6e, 0x69, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x2e, 0x76,
	0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x6d, 0x69, 0x6e, 0x69, 0x62,
	0x6f, 0x61, 0x72, 0x64, 0x2e, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x19, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x13, 0x22, 0x0e, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x3a, 0x01, 0x2a, 0x42,
	0x08, 0x5a, 0x06, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_tokens_service_proto_rawDescOnce sync.Once
	file_tokens_service_proto_rawDescData = file_tokens_service_proto_rawDesc
)

func file_tokens_service_proto_rawDescGZIP() []byte {
	file_tokens_service_proto_rawDescOnce.Do(func() {
		file_tokens_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_tokens_service_proto_rawDescData)
	})
	return file_tokens_service_proto_rawDescData
}

var file_tokens_service_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_tokens_service_proto_goTypes = []interface{}{
	(*CreateTokenRequest)(nil), // 0: app.miniboard.tokens.v1.CreateTokenRequest
	(*Token)(nil),              // 1: app.miniboard.tokens.v1.Token
}
var file_tokens_service_proto_depIdxs = []int32{
	0, // 0: app.miniboard.tokens.v1.TokensService.CreateToken:input_type -> app.miniboard.tokens.v1.CreateTokenRequest
	1, // 1: app.miniboard.tokens.v1.TokensService.CreateToken:output_type -> app.miniboard.tokens.v1.Token
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_tokens_service_proto_init() }
func file_tokens_service_proto_init() {
	if File_tokens_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_tokens_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateTokenRequest); i {
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
		file_tokens_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Token); i {
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
			RawDescriptor: file_tokens_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_tokens_service_proto_goTypes,
		DependencyIndexes: file_tokens_service_proto_depIdxs,
		MessageInfos:      file_tokens_service_proto_msgTypes,
	}.Build()
	File_tokens_service_proto = out.File
	file_tokens_service_proto_rawDesc = nil
	file_tokens_service_proto_goTypes = nil
	file_tokens_service_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// TokensServiceClient is the client API for TokensService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TokensServiceClient interface {
	// Create authorization token.
	//
	// Returns an authorization token.
	CreateToken(ctx context.Context, in *CreateTokenRequest, opts ...grpc.CallOption) (*Token, error)
}

type tokensServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTokensServiceClient(cc grpc.ClientConnInterface) TokensServiceClient {
	return &tokensServiceClient{cc}
}

func (c *tokensServiceClient) CreateToken(ctx context.Context, in *CreateTokenRequest, opts ...grpc.CallOption) (*Token, error) {
	out := new(Token)
	err := c.cc.Invoke(ctx, "/app.miniboard.tokens.v1.TokensService/CreateToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TokensServiceServer is the server API for TokensService service.
type TokensServiceServer interface {
	// Create authorization token.
	//
	// Returns an authorization token.
	CreateToken(context.Context, *CreateTokenRequest) (*Token, error)
}

// UnimplementedTokensServiceServer can be embedded to have forward compatible implementations.
type UnimplementedTokensServiceServer struct {
}

func (*UnimplementedTokensServiceServer) CreateToken(context.Context, *CreateTokenRequest) (*Token, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateToken not implemented")
}

func RegisterTokensServiceServer(s *grpc.Server, srv TokensServiceServer) {
	s.RegisterService(&_TokensService_serviceDesc, srv)
}

func _TokensService_CreateToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokensServiceServer).CreateToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.miniboard.tokens.v1.TokensService/CreateToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokensServiceServer).CreateToken(ctx, req.(*CreateTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _TokensService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "app.miniboard.tokens.v1.TokensService",
	HandlerType: (*TokensServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateToken",
			Handler:    _TokensService_CreateToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "tokens_service.proto",
}
