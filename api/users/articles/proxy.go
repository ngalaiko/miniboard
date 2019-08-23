package articles // import "miniboard.app/api/articles"

import (
	"context"

	"google.golang.org/grpc"
	"miniboard.app/proto/users/articles/v1"
)

var _ articles.ArticlesServiceClient = &Proxy{}

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

// ListArticles calls method ListArticles on the service instance.
func (p *Proxy) ListArticles(ctx context.Context, in *articles.ListArticlesRequest, opts ...grpc.CallOption) (*articles.ListArticlesResponse, error) {
	return p.service.ListArticles(ctx, in)
}

// CreateArticle calls CreateArticle method on the service instance.
func (p *Proxy) CreateArticle(ctx context.Context, in *articles.CreateArticleRequest, opts ...grpc.CallOption) (*articles.Article, error) {
	return p.service.CreateArticle(ctx, in)
}
