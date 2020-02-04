// Code generated by protoc-gen-go. DO NOT EDIT.
// source: codes_service.proto

package codes

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type CreateCodeRequest struct {
	// Email address.
	Email string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	// Response type.
	ResponseType         string   `protobuf:"bytes,2,opt,name=response_type,json=responseType,proto3" json:"response_type,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateCodeRequest) Reset()         { *m = CreateCodeRequest{} }
func (m *CreateCodeRequest) String() string { return proto.CompactTextString(m) }
func (*CreateCodeRequest) ProtoMessage()    {}
func (*CreateCodeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e5a73e0b96836af0, []int{0}
}

func (m *CreateCodeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateCodeRequest.Unmarshal(m, b)
}
func (m *CreateCodeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateCodeRequest.Marshal(b, m, deterministic)
}
func (m *CreateCodeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateCodeRequest.Merge(m, src)
}
func (m *CreateCodeRequest) XXX_Size() int {
	return xxx_messageInfo_CreateCodeRequest.Size(m)
}
func (m *CreateCodeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateCodeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateCodeRequest proto.InternalMessageInfo

func (m *CreateCodeRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *CreateCodeRequest) GetResponseType() string {
	if m != nil {
		return m.ResponseType
	}
	return ""
}

type Code struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Code) Reset()         { *m = Code{} }
func (m *Code) String() string { return proto.CompactTextString(m) }
func (*Code) ProtoMessage()    {}
func (*Code) Descriptor() ([]byte, []int) {
	return fileDescriptor_e5a73e0b96836af0, []int{1}
}

func (m *Code) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Code.Unmarshal(m, b)
}
func (m *Code) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Code.Marshal(b, m, deterministic)
}
func (m *Code) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Code.Merge(m, src)
}
func (m *Code) XXX_Size() int {
	return xxx_messageInfo_Code.Size(m)
}
func (m *Code) XXX_DiscardUnknown() {
	xxx_messageInfo_Code.DiscardUnknown(m)
}

var xxx_messageInfo_Code proto.InternalMessageInfo

func init() {
	proto.RegisterType((*CreateCodeRequest)(nil), "app.miniboard.codes.v1.CreateCodeRequest")
	proto.RegisterType((*Code)(nil), "app.miniboard.codes.v1.Code")
}

func init() { proto.RegisterFile("codes_service.proto", fileDescriptor_e5a73e0b96836af0) }

var fileDescriptor_e5a73e0b96836af0 = []byte{
	// 184 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4e, 0xce, 0x4f, 0x49,
	0x2d, 0x8e, 0x2f, 0x4e, 0x2d, 0x2a, 0xcb, 0x4c, 0x4e, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17,
	0x12, 0x4b, 0x2c, 0x28, 0xd0, 0xcb, 0xcd, 0xcc, 0xcb, 0x4c, 0xca, 0x4f, 0x2c, 0x4a, 0xd1, 0x03,
	0x2b, 0xd1, 0x2b, 0x33, 0x54, 0xf2, 0xe3, 0x12, 0x74, 0x2e, 0x4a, 0x4d, 0x2c, 0x49, 0x75, 0xce,
	0x4f, 0x49, 0x0d, 0x4a, 0x2d, 0x2c, 0x4d, 0x2d, 0x2e, 0x11, 0x12, 0xe1, 0x62, 0x4d, 0xcd, 0x4d,
	0xcc, 0xcc, 0x91, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x82, 0x70, 0x84, 0x94, 0xb9, 0x78, 0x8b,
	0x52, 0x8b, 0x0b, 0xf2, 0xf3, 0x8a, 0x53, 0xe3, 0x4b, 0x2a, 0x0b, 0x52, 0x25, 0x98, 0xc0, 0xb2,
	0x3c, 0x30, 0xc1, 0x90, 0xca, 0x82, 0x54, 0x25, 0x36, 0x2e, 0x16, 0x90, 0x49, 0x46, 0xe9, 0x5c,
	0x3c, 0x20, 0xba, 0x38, 0x18, 0xe2, 0x0a, 0xa1, 0x70, 0x2e, 0x2e, 0x84, 0x3d, 0x42, 0x9a, 0x7a,
	0xd8, 0x9d, 0xa3, 0x87, 0xe1, 0x16, 0x29, 0x19, 0x9c, 0x4a, 0xf3, 0x53, 0x52, 0x95, 0x18, 0x9c,
	0xd8, 0xa3, 0x58, 0xc1, 0x42, 0x49, 0x6c, 0x60, 0x8f, 0x1a, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff,
	0xd7, 0xc4, 0x37, 0x23, 0xff, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// CodesServiceClient is the client API for CodesService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CodesServiceClient interface {
	// Create authorization code.
	//
	// Returns an authorization code.
	CreateCode(ctx context.Context, in *CreateCodeRequest, opts ...grpc.CallOption) (*Code, error)
}

type codesServiceClient struct {
	cc *grpc.ClientConn
}

func NewCodesServiceClient(cc *grpc.ClientConn) CodesServiceClient {
	return &codesServiceClient{cc}
}

func (c *codesServiceClient) CreateCode(ctx context.Context, in *CreateCodeRequest, opts ...grpc.CallOption) (*Code, error) {
	out := new(Code)
	err := c.cc.Invoke(ctx, "/app.miniboard.codes.v1.CodesService/CreateCode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CodesServiceServer is the server API for CodesService service.
type CodesServiceServer interface {
	// Create authorization code.
	//
	// Returns an authorization code.
	CreateCode(context.Context, *CreateCodeRequest) (*Code, error)
}

// UnimplementedCodesServiceServer can be embedded to have forward compatible implementations.
type UnimplementedCodesServiceServer struct {
}

func (*UnimplementedCodesServiceServer) CreateCode(ctx context.Context, req *CreateCodeRequest) (*Code, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCode not implemented")
}

func RegisterCodesServiceServer(s *grpc.Server, srv CodesServiceServer) {
	s.RegisterService(&_CodesService_serviceDesc, srv)
}

func _CodesService_CreateCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CodesServiceServer).CreateCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.miniboard.codes.v1.CodesService/CreateCode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CodesServiceServer).CreateCode(ctx, req.(*CreateCodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _CodesService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "app.miniboard.codes.v1.CodesService",
	HandlerType: (*CodesServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateCode",
			Handler:    _CodesService_CreateCode_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "codes_service.proto",
}
