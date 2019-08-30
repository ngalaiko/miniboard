package articles // "miniboard.app/api/articles"

import (
	"context"
	"encoding/base64"
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
	lookFor := resource.ParseName(request.Parent).Child("articles", "*")

	var from *resource.Name
	if request.PageToken != "" {
		decoded, err := base64.StdEncoding.DecodeString(request.PageToken)
		if err != nil {
			return nil, status.New(codes.InvalidArgument, "invalid page token").Err()
		}
		from = resource.ParseName(string(decoded))
	}

	dd, err := s.storage.LoadChildren(lookFor, from, int(request.PageSize+1))
	if err != nil {
		return nil, status.New(codes.Internal, "failed to load articles").Err()
	}

	aa := make([]*articles.Article, 0, len(dd))
	for _, d := range dd {
		a := &articles.Article{}
		if err := proto.Unmarshal(d.Data, a); err != nil {
			return nil, status.New(codes.Internal, "failed to unmarshal article").Err()
		}
		aa = append(aa, a)
	}

	var nextPageToken string
	if len(aa) == int(request.PageSize+1) {
		nextPageToken = base64.StdEncoding.EncodeToString([]byte(aa[len(aa)-1].Name))
		aa = aa[:request.PageSize]
	}

	return &articles.ListArticlesResponse{
		Articles:      aa,
		NextPageToken: nextPageToken,
	}, nil
}

// CreateArticle creates a new article.
func (s *Service) CreateArticle(ctx context.Context, request *articles.CreateArticleRequest) (*articles.Article, error) {
	if request.Article.Url == "" {
		return nil, status.New(codes.InvalidArgument, "url is empty").Err()
	}

	if _, err := url.ParseRequestURI(request.Article.Url); err != nil {
		return nil, status.New(codes.InvalidArgument, "url is invalid").Err()
	}

	name := resource.ParseName(request.Parent).Child("articles", ksuid.New().String())

	request.Article.Name = name.String()

	rawArticle, err := proto.Marshal(request.Article)
	if err != nil {
		return nil, status.New(codes.Internal, "failed to marshal the article").Err()
	}

	if err := s.storage.Store(name, rawArticle); err != nil {
		return nil, status.New(codes.Internal, "failed to store the article").Err()
	}

	return request.Article, nil
}

// GetArticle returns an article.
func (s *Service) GetArticle(ctx context.Context, request *articles.GetArticleRequest) (*articles.Article, error) {
	name := resource.ParseName(request.Name)

	raw, err := s.storage.Load(name)
	if err != nil {
		return nil, status.New(codes.Internal, "failed to load the article").Err()
	}

	article := &articles.Article{}
	if err := proto.Unmarshal(raw, article); err != nil {
		return nil, status.New(codes.Internal, "failed to unmarshal the article").Err()
	}

	return article, nil
}
