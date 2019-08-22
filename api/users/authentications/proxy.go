package authentications // import "miniboard.app/api/users/authentications"

import (
	"context"

	"google.golang.org/grpc"
	"miniboard.app/proto/users/authentications/v1"
)

var _ authentications.AuthenticationsServiceClient = &Proxy{}

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

// CreateAuthentication calls CreateAuthentication method on the service instance.
func (p *Proxy) CreateAuthentication(ctx context.Context, in *authentications.CreateAuthenticationRequest, opts ...grpc.CallOption) (*authentications.Authentication, error) {
	return p.service.CreateAuthentication(ctx, in)
}
