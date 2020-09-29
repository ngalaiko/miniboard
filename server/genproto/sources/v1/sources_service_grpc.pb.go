// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package sources

import (
	context "context"
	longrunning "google.golang.org/genproto/googleapis/longrunning"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// SourcesServiceClient is the client API for SourcesService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SourcesServiceClient interface {
	// Create source
	//
	// Allows to create a new source for a user.
	CreateSource(ctx context.Context, in *CreateSourceRequest, opts ...grpc.CallOption) (*longrunning.Operation, error)
}

type sourcesServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSourcesServiceClient(cc grpc.ClientConnInterface) SourcesServiceClient {
	return &sourcesServiceClient{cc}
}

func (c *sourcesServiceClient) CreateSource(ctx context.Context, in *CreateSourceRequest, opts ...grpc.CallOption) (*longrunning.Operation, error) {
	out := new(longrunning.Operation)
	err := c.cc.Invoke(ctx, "/miniboard.sources.v1.SourcesService/CreateSource", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SourcesServiceServer is the server API for SourcesService service.
type SourcesServiceServer interface {
	// Create source
	//
	// Allows to create a new source for a user.
	CreateSource(context.Context, *CreateSourceRequest) (*longrunning.Operation, error)
}

// UnimplementedSourcesServiceServer can be embedded to have forward compatible implementations.
type UnimplementedSourcesServiceServer struct {
}

func (*UnimplementedSourcesServiceServer) CreateSource(context.Context, *CreateSourceRequest) (*longrunning.Operation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSource not implemented")
}

func RegisterSourcesServiceServer(s *grpc.Server, srv SourcesServiceServer) {
	s.RegisterService(&_SourcesService_serviceDesc, srv)
}

func _SourcesService_CreateSource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SourcesServiceServer).CreateSource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/miniboard.sources.v1.SourcesService/CreateSource",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SourcesServiceServer).CreateSource(ctx, req.(*CreateSourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _SourcesService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "miniboard.sources.v1.SourcesService",
	HandlerType: (*SourcesServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateSource",
			Handler:    _SourcesService_CreateSource_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sources/v1/sources_service.proto",
}
