package articles

import (
	"bytes"
	"context"
	"encoding/base64"
	"net/url"
	"strings"
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

// Service controls articles resource.
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

// SearchArticles searches through the articles.
func (s *Service) SearchArticles(ctx context.Context, request *articles.SearchArticlesRequest) (*articles.SearchArticlesResponse, error) {
	lookFor := resource.ParseName(request.Parent).Child("articles", "*")

	if len(request.Query) == 0 {
		return &articles.SearchArticlesResponse{}, nil
	}

	var from *resource.Name
	if request.PageToken != "" {
		decoded, err := base64.StdEncoding.DecodeString(request.PageToken)
		if err != nil {
			return nil, status.New(codes.InvalidArgument, "invalid page token").Err()
		}
		from = resource.ParseName(string(decoded))
	}

	aa := make([]*articles.Article, 0, request.PageSize+1)
	err := s.storage.ForEach(lookFor, from, func(r *resource.Resource) (bool, error) {
		if int64(len(aa)) == request.PageSize {
			return false, nil
		}

		a := &articles.Article{}
		if err := proto.Unmarshal(r.Data, a); err != nil {
			return false, status.New(codes.Internal, "failed to unmarshal article").Err()
		}

		if bytes.Contains(
			bytes.ToLower([]byte(a.Title)),
			bytes.ToLower([]byte(request.Query)),
		) {
			aa = append(aa, a)
			return true, nil
		}

		if strings.Contains(a.Url, request.Query) {
			aa = append(aa, a)
			return true, nil
		}

		return true, nil
	})

	var nextPageToken string
	if len(aa) == int(request.PageSize+1) {
		nextPageToken = base64.StdEncoding.EncodeToString([]byte(aa[len(aa)-1].Name))
		aa = aa[:request.PageSize]
	}

	switch err {
	case nil, storage.ErrNotFound:
		return &articles.SearchArticlesResponse{
			Articles:      aa,
			NextPageToken: nextPageToken,
		}, nil
	default:
		return nil, status.New(codes.Internal, "failed to search articles").Err()
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
	switch errors.Cause(err) {
	case storage.ErrNotFound, nil:
	default:
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

	name := resource.ParseName(request.Parent).Child("articles", ksuid.New().String())

	r, err := s.newReader(articleURL)
	var content []byte
	switch err {
	case nil:
		request.Article.Title = r.Title()
		request.Article.IconUrl = r.IconURL()

		content = r.Content()
		if content != nil {
			if err := s.storage.Store(resource.NewName("content", name.ID()), content); err != nil {
				return nil, status.New(codes.Internal, "failed to store the article content").Err()
			}
		}
	}

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

	request.Article.Content = content
	return request.Article, nil
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
		case "is_read":
			if article.IsRead == request.Article.IsRead {
				continue
			}
			article.IsRead = request.Article.IsRead
			updated = true
		default:
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
	article, err := s.getArticle(ctx, name)
	if err != nil {
		return nil, err
	}
	if request.View != articles.ArticleView_ARTICLE_VIEW_FULL {
		return article, nil
	}

	article.Content, err = s.storage.Load(resource.NewName("content", name.ID()))
	switch errors.Cause(err) {
	case nil:
	case storage.ErrNotFound:
		return article, nil
	default:
		return nil, status.New(codes.Internal, "failed to load the article's content").Err()
	}

	return article, nil
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
