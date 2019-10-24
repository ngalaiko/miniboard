package users

import (
	"context"

	"google.golang.org/grpc"
	"miniboard.app/proto/users/v1"
)

var _ users.UsersServiceClient = &Proxy{}

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

// GetUser calls GetUser method on the service instance.
func (p *Proxy) GetUser(ctx context.Context, in *users.GetUserRequest, opts ...grpc.CallOption) (*users.User, error) {
	return p.service.GetUser(ctx, in)
}
