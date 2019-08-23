package articles // "miniboard.app/api/articles"

import (
	"context"
	"net/url"

	"github.com/gogo/protobuf/proto"
	"github.com/segmentio/ksuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"miniboard.app/proto/users/articles/v1"
	"miniboard.app/storage"
	"miniboard.app/storage/resource"
)

// Service controlls articles resource.
type Service struct {
	storage storage.Storage
}

// New returns a new articles service instance.
func New(storage storage.Storage) *Service {
	return &Service{
		storage: storage,
	}
}

// ListArticles returns a list of articles.
func (s *Service) ListArticles(ctx context.Context, request *articles.ListArticlesRequest) (*articles.ListArticlesResponse, error) {
	return nil, nil
}

// CreateArticle creates a new article.
func (s *Service) CreateArticle(ctx context.Context, request *articles.CreateArticleRequest) (*articles.Article, error) {
	if _, err := url.Parse(request.Article.Url); err != nil {
		return nil, status.New(codes.InvalidArgument, "url is invalid").Err()
	}

	name := resource.ParseName(request.Parent).Child("articles", ksuid.New().String())

	article := request.Article
	article.Name = name.String()

	rawArticle, err := proto.Marshal(article)
	if err != nil {
		return nil, status.New(codes.Internal, "failed to marshal the article").Err()
	}

	if err := s.storage.Store(name, rawArticle); err != nil {
		return nil, status.New(codes.Internal, "failed to store the article").Err()
	}

	return article, nil
}
