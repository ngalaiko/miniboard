package labels

import (
	"context"

	"google.golang.org/grpc"
	"miniboard.app/proto/users/labels/v1"
)

var _ labels.LabelsServiceClient = &Proxy{}

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

// CreateLabel calls CreateLabel method on the service instance.
func (p *Proxy) CreateLabel(ctx context.Context, in *labels.CreateLabelRequest, opts ...grpc.CallOption) (*labels.Label, error) {
	return p.service.CreateLabel(ctx, in)
}

// GetLabel calls GetLabel method on the service instance.
func (p *Proxy) GetLabel(ctx context.Context, in *labels.GetLabelRequest, opts ...grpc.CallOption) (*labels.Label, error) {
	return p.service.GetLabel(ctx, in)
}
