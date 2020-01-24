package articles

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"miniboard.app/api/actor"
	"miniboard.app/images"
	"miniboard.app/proto/users/articles/v1"
	"miniboard.app/reader"
	"miniboard.app/storage"
	"miniboard.app/storage/resource"
)

// Service controls articles resource.
type Service struct {
	storage storage.Storage

	client reader.GetClient
	images *images.Service
}

// New returns a new articles service instance.
func New(storage storage.Storage, images *images.Service) *Service {
	return &Service{
		storage: storage,
		client:  &http.Client{},
		images:  images,
	}
}

// ListArticles returns a list of articles.
func (s *Service) ListArticles(ctx context.Context, request *articles.ListArticlesRequest) (*articles.ListArticlesResponse, error) {
	lookFor := resource.ParseName(request.Parent).Child("articles", "*")

	if !actor.Owns(ctx, lookFor) {
		return nil, grpc.Errorf(codes.PermissionDenied, "forbidden")
	}

	var from *resource.Name
	if request.PageToken != "" {
		decoded, err := base64.StdEncoding.DecodeString(request.PageToken)
		if err != nil {
			return nil, grpc.Errorf(codes.InvalidArgument, "invalid page token")
		}
		from = resource.ParseName(string(decoded))
	}

	aa := []*articles.Article{}
	err := s.storage.ForEach(ctx, lookFor, from, func(r *resource.Resource) (bool, error) {
		if int64(len(aa)) == request.PageSize+1 {
			return false, nil
		}

		a := &articles.Article{}
		if err := proto.Unmarshal(r.Data, a); err != nil {
			return false, grpc.Errorf(codes.Internal, "failed to unmarshal article")
		}

		if request.IsRead != nil && a.IsRead != request.IsRead.GetValue() {
			return true, nil
		}

		if request.IsFavorite != nil && a.IsFavorite != request.IsFavorite.GetValue() {
			return true, nil
		}

		if request.Title != nil && !strings.Contains(
			strings.ToLower(a.Title),
			strings.ToLower(request.Title.GetValue()),
		) {
			return true, nil
		}

		if request.Url != nil && !strings.Contains(
			strings.ToLower(a.Url),
			strings.ToLower(request.Url.GetValue()),
		) {
			return true, nil
		}

		aa = append(aa, a)

		return true, nil
	})

	var nextPageToken string
	if len(aa) == int(request.PageSize+1) {
		nextPageToken = base64.StdEncoding.EncodeToString([]byte(aa[len(aa)-1].Name))
		aa = aa[:request.PageSize]
	}

	switch err {
	case nil, storage.ErrNotFound:
		return &articles.ListArticlesResponse{
			Articles:      aa,
			NextPageToken: nextPageToken,
		}, nil
	default:
		log("articles.list").Error(err)
		return nil, grpc.Errorf(codes.Internal, "failed to list articles")
	}
}

// CreateArticle creates a new article.
func (s *Service) CreateArticle(ctx context.Context, request *articles.CreateArticleRequest) (*articles.Article, error) {
	now := time.Now()

	if request.Article.Url == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "url is empty")
	}

	articleURL, err := url.ParseRequestURI(request.Article.Url)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "url is invalid")
	}

	name := resource.ParseName(request.Parent).Child("articles", ksuid.New().String())

	if !actor.Owns(ctx, name) {
		return nil, grpc.Errorf(codes.PermissionDenied, "forbidden")
	}

	r, err := reader.New(ctx, s.client, name, s.images, articleURL)
	var content []byte
	switch err {
	case nil:
		request.Article.Title = r.Title()
		request.Article.SiteName = r.SiteName()
		request.Article.IconUrl = r.IconURL()

		content = r.Content()
		if content != nil {
			if err := s.storage.Store(ctx, resource.NewName("content", name.ID()), content); err != nil {
				return nil, grpc.Errorf(codes.Internal, "failed to store the article content")
			}
		}
	}

	request.Article.Name = name.String()
	request.Article.CreateTime = &timestamp.Timestamp{
		Seconds: now.In(time.UTC).Unix(),
	}

	rawArticle, err := proto.Marshal(request.Article)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to marshal the article")
	}

	if err := s.storage.Store(ctx, name, rawArticle); err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to store the article")
	}

	request.Article.Content = content
	return request.Article, nil
}

// UpdateArticle updates the article.
func (s *Service) UpdateArticle(ctx context.Context, request *articles.UpdateArticleRequest) (*articles.Article, error) {
	name := resource.ParseName(request.Article.Name)

	if !actor.Owns(ctx, name) {
		return nil, grpc.Errorf(codes.PermissionDenied, "forbidden")
	}

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
		case "is_favorite":
			if article.IsFavorite == request.Article.IsFavorite {
				continue
			}
			article.IsFavorite = request.Article.IsFavorite
			updated = true
		default:
		}
	}

	if !updated {
		return article, nil
	}

	rawArticle, err := proto.Marshal(article)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to marshal the article")
	}

	if err := s.storage.Update(ctx, name, rawArticle); err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to store the article")
	}

	return article, nil
}

// GetArticle returns an article.
func (s *Service) GetArticle(ctx context.Context, request *articles.GetArticleRequest) (*articles.Article, error) {
	name := resource.ParseName(request.Name)

	if !actor.Owns(ctx, name) {
		return nil, grpc.Errorf(codes.PermissionDenied, "forbidden")
	}

	article, err := s.getArticle(ctx, name)
	if err != nil {
		return nil, err
	}
	if request.View != articles.ArticleView_ARTICLE_VIEW_FULL {
		return article, nil
	}

	article.Content, err = s.storage.Load(ctx, resource.NewName("content", name.ID()))
	switch errors.Cause(err) {
	case nil:
	case storage.ErrNotFound:
		return article, nil
	default:
		return nil, grpc.Errorf(codes.Internal, "failed to load the article's content")
	}

	return article, nil
}

func (s *Service) getArticle(ctx context.Context, name *resource.Name) (*articles.Article, error) {
	raw, err := s.storage.Load(ctx, name)
	switch errors.Cause(err) {
	case nil:
	case storage.ErrNotFound:
		return nil, grpc.Errorf(codes.NotFound, "not found")
	default:
		return nil, grpc.Errorf(codes.Internal, "failed to load the article")
	}

	article := &articles.Article{}
	if err := proto.Unmarshal(raw, article); err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to unmarshal the article")
	}

	return article, nil
}

// DeleteArticle removes an article.
func (s *Service) DeleteArticle(ctx context.Context, request *articles.DeleteArticleRequest) (*empty.Empty, error) {
	name := resource.ParseName(request.Name)

	if !actor.Owns(ctx, name) {
		return nil, grpc.Errorf(codes.PermissionDenied, "forbidden")
	}

	if err := s.storage.Delete(ctx, name); err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to delete the article")
	}

	return &empty.Empty{}, nil
}

func log(src string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": src,
	})
}
