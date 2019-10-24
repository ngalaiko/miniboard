package codes

import (
	"context"

	"google.golang.org/grpc"
	"miniboard.app/proto/codes/v1"
)

var _ codes.CodesServiceClient = &Proxy{}

// Proxy is a gRPC client that proxies all calls to the server.
type Proxy struct {
	service *Service
}

// NewProxyClient returns new proxy client to the service.
func NewProxyClient(service *Service) *Proxy {
	return &Proxy{
		service: service,
	}
}

// CreateCode creates new authorization code.
func (s *Proxy) CreateCode(ctx context.Context, request *codes.CreateCodeRequest, opts ...grpc.CallOption) (*codes.Code, error) {
	return s.service.CreateCode(ctx, request)
}
