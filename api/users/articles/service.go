package articles

import (
	"context"
	"encoding/base64"
	"net/url"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"miniboard.app/proto/users/articles/v1"
	"miniboard.app/reader"
	"miniboard.app/reader/http"
	"miniboard.app/storage"
	"miniboard.app/storage/resource"
)

// Service controlls articles resource.
type Service struct {
	storage storage.Storage

	newReader func(*url.URL) (reader.Reader, error)
}

// New returns a new articles service instance.
func New(storage storage.Storage) *Service {
	return &Service{
		storage:   storage,
		newReader: http.New,
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
	now := time.Now()

	if request.Article.Url == "" {
		return nil, status.New(codes.InvalidArgument, "url is empty").Err()
	}

	articleURL, err := url.ParseRequestURI(request.Article.Url)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, "url is invalid").Err()
	}

	r, err := s.newReader(articleURL)
	if err == nil {
		enrich(request.Article, r)
	}

	name := resource.ParseName(request.Parent).Child("articles", ksuid.New().String())

	request.Article.Name = name.String()
	request.Article.CreateTime = &timestamp.Timestamp{
		Seconds: now.In(time.UTC).Unix(),
	}

	rawArticle, err := proto.Marshal(request.Article)
	if err != nil {
		return nil, status.New(codes.Internal, "failed to marshal the article").Err()
	}

	if err := s.storage.Store(name, rawArticle); err != nil {
		return nil, status.New(codes.Internal, "failed to store the article").Err()
	}

	return request.Article, nil
}

func enrich(article *articles.Article, r reader.Reader) {
	article.Title = r.Title()
	article.IconUrl = r.IconURL()
	article.Content = r.Content()
}

// UpdateArticle updates the article.
func (s *Service) UpdateArticle(ctx context.Context, request *articles.UpdateArticleRequest) (*articles.Article, error) {
	name := resource.ParseName(request.Article.Name)

	article, err := s.getArticle(ctx, name)
	if err != nil {
		return nil, err
	}

	var updated bool

	for _, path := range request.UpdateMask.GetPaths() {
		switch path {
		case "label_ids":
			article.LabelIds = request.Article.LabelIds
			updated = true
		}
	}

	if !updated {
		return article, nil
	}

	rawArticle, err := proto.Marshal(article)
	if err != nil {
		return nil, status.New(codes.Internal, "failed to marshal the article").Err()
	}

	if err := s.storage.Update(name, rawArticle); err != nil {
		return nil, status.New(codes.Internal, "failed to store the article").Err()
	}

	return article, nil
}

// GetArticle returns an article.
func (s *Service) GetArticle(ctx context.Context, request *articles.GetArticleRequest) (*articles.Article, error) {
	name := resource.ParseName(request.Name)
	return s.getArticle(ctx, name)
}

func (s *Service) getArticle(ctx context.Context, name *resource.Name) (*articles.Article, error) {
	raw, err := s.storage.Load(name)
	switch errors.Cause(err) {
	case nil:
	case storage.ErrNotFound:
		return nil, status.New(codes.NotFound, "not found").Err()
	default:
		return nil, status.New(codes.Internal, "failed to load the article").Err()
	}

	article := &articles.Article{}
	if err := proto.Unmarshal(raw, article); err != nil {
		return nil, status.New(codes.Internal, "failed to unmarshal the article").Err()
	}

	return article, nil
}

// DeleteArticle removes an article.
func (s *Service) DeleteArticle(ctx context.Context, request *articles.DeleteArticleRequest) (*empty.Empty, error) {
	name := resource.ParseName(request.Name)

	if err := s.storage.Delete(name); err != nil {
		return nil, status.New(codes.Internal, "failed to delete the article").Err()
	}

	return &empty.Empty{}, nil
}
