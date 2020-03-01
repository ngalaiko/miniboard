// Code generated by protoc-gen-go. DO NOT EDIT.
// source: rss_service.proto

package rss

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	grpc "google.golang.org/grpc"
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

type Feed struct {
	// Name of the resource, for example "users/user1/feeds/feed1".
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Time of the last feed fetch by the reader.
	LastFetched          *timestamp.Timestamp `protobuf:"bytes,2,opt,name=last_fetched,json=lastFetched,proto3" json:"last_fetched,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Feed) Reset()         { *m = Feed{} }
func (m *Feed) String() string { return proto.CompactTextString(m) }
func (*Feed) ProtoMessage()    {}
func (*Feed) Descriptor() ([]byte, []int) {
	return fileDescriptor_2940ccdfc1dcaa8e, []int{0}
}

func (m *Feed) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Feed.Unmarshal(m, b)
}
func (m *Feed) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Feed.Marshal(b, m, deterministic)
}
func (m *Feed) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Feed.Merge(m, src)
}
func (m *Feed) XXX_Size() int {
	return xxx_messageInfo_Feed.Size(m)
}
func (m *Feed) XXX_DiscardUnknown() {
	xxx_messageInfo_Feed.DiscardUnknown(m)
}

var xxx_messageInfo_Feed proto.InternalMessageInfo

func (m *Feed) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Feed) GetLastFetched() *timestamp.Timestamp {
	if m != nil {
		return m.LastFetched
	}
	return nil
}

func init() {
	proto.RegisterType((*Feed)(nil), "app.miniboard.users.rss.v1.Feed")
}

func init() { proto.RegisterFile("rss_service.proto", fileDescriptor_2940ccdfc1dcaa8e) }

var fileDescriptor_2940ccdfc1dcaa8e = []byte{
	// 182 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x34, 0x8e, 0xb1, 0x6e, 0x83, 0x30,
	0x14, 0x45, 0x45, 0x4b, 0x2b, 0xd5, 0xb0, 0xd4, 0x13, 0x62, 0x29, 0xea, 0xc4, 0xf4, 0x50, 0xc8,
	0x9c, 0x25, 0x03, 0x1f, 0x00, 0x59, 0x92, 0x05, 0x19, 0x78, 0x10, 0x4b, 0x18, 0x5b, 0x7e, 0x86,
	0xef, 0x8f, 0x82, 0xc5, 0x76, 0x87, 0x73, 0x8e, 0x2e, 0xfb, 0xb5, 0x44, 0x2d, 0xa1, 0xdd, 0x64,
	0x8f, 0x60, 0xac, 0x76, 0x9a, 0xa7, 0xc2, 0x18, 0x50, 0x72, 0x91, 0x9d, 0x16, 0x76, 0x80, 0x95,
	0xd0, 0x12, 0x58, 0x22, 0xd8, 0x4e, 0xe9, 0xdf, 0xa4, 0xf5, 0x34, 0x63, 0xb1, 0x93, 0xdd, 0x3a,
	0x16, 0x4e, 0x2a, 0x24, 0x27, 0x94, 0xf1, 0xf2, 0xff, 0x9d, 0x85, 0x15, 0xe2, 0xc0, 0x39, 0x0b,
	0x17, 0xa1, 0x30, 0x09, 0xb2, 0x20, 0xff, 0xa9, 0xf7, 0xcd, 0x2f, 0x2c, 0x9e, 0x05, 0xb9, 0x76,
	0x44, 0xd7, 0x3f, 0x71, 0x48, 0x3e, 0xb2, 0x20, 0x8f, 0xca, 0x14, 0x7c, 0x13, 0x8e, 0x26, 0xdc,
	0x8e, 0x66, 0x1d, 0xbd, 0xf9, 0xca, 0xe3, 0x65, 0xcc, 0x58, 0xdd, 0x34, 0x8d, 0xff, 0x7a, 0xfd,
	0x7a, 0x7c, 0x5a, 0xa2, 0xee, 0x7b, 0xb7, 0xce, 0xaf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xea, 0x57,
	0x74, 0xc9, 0xc8, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// RSSServiceClient is the client API for RSSService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RSSServiceClient interface {
}

type rSSServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRSSServiceClient(cc grpc.ClientConnInterface) RSSServiceClient {
	return &rSSServiceClient{cc}
}

// RSSServiceServer is the server API for RSSService service.
type RSSServiceServer interface {
}

// UnimplementedRSSServiceServer can be embedded to have forward compatible implementations.
type UnimplementedRSSServiceServer struct {
}

func RegisterRSSServiceServer(s *grpc.Server, srv RSSServiceServer) {
	s.RegisterService(&_RSSService_serviceDesc, srv)
}

var _RSSService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "app.miniboard.users.rss.v1.RSSService",
	HandlerType: (*RSSServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "rss_service.proto",
}
