package authorizations // import "miniboard.app/api/users/authorizations"

import (
	"context"

	"google.golang.org/grpc"
	"miniboard.app/proto/users/authorizations/v1"
)

var _ authorizations.AuthorizationsServiceClient = &Proxy{}

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

// CreateAuthorization calls CreateAuthorization method on the service instance.
func (p *Proxy) CreateAuthorization(ctx context.Context, in *authorizations.CreateAuthorizationRequest, opts ...grpc.CallOption) (*authorizations.Authorization, error) {
	return p.service.CreateAuthorization(ctx, in)
}
