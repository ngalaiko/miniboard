package articles

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
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

// UpdateArticle calls UpdateArticle method on the service instance.
func (p *Proxy) UpdateArticle(ctx context.Context, in *articles.UpdateArticleRequest, opts ...grpc.CallOption) (*articles.Article, error) {
	return p.service.UpdateArticle(ctx, in)
}

// GetArticle calls GetArticle method on the service instance.
func (p *Proxy) GetArticle(ctx context.Context, in *articles.GetArticleRequest, opts ...grpc.CallOption) (*articles.Article, error) {
	return p.service.GetArticle(ctx, in)
}

// DeleteArticle calls DeleteArticle method on the service instance.
func (p *Proxy) DeleteArticle(ctx context.Context, in *articles.DeleteArticleRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	return p.service.DeleteArticle(ctx, in)
}
