package articles

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
	"github.com/spaolacci/murmur3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"miniboard.app/api/actor"
	"miniboard.app/fetch"
	"miniboard.app/reader"
	"miniboard.app/storage"
	"miniboard.app/storage/resource"
)

var (
	ErrAlreadyExists = errors.New("article already exists")
)

// Service controls articles resource.
type Service struct {
	storage storage.Storage

	client fetch.Fetcher
}

// NewService returns a new articles service instance.
func NewService(storage storage.Storage, fetcher fetch.Fetcher) *Service {
	return &Service{
		storage: storage,
		client:  fetcher,
	}
}

// ListArticles returns a list of articles.
func (s *Service) ListArticles(ctx context.Context, request *ListArticlesRequest) (*ListArticlesResponse, error) {
	actor, _ := actor.FromContext(ctx)
	lookFor := actor.Child("articles", "*")

	var from *resource.Name
	if request.PageToken != "" {
		decoded, err := base64.StdEncoding.DecodeString(request.PageToken)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid page token")
		}
		from = resource.ParseName(string(decoded))
	}

	aa := []*Article{}
	err := s.storage.ForEach(ctx, lookFor, from, func(r *resource.Resource) (bool, error) {
		a := &Article{}
		if err := proto.Unmarshal(r.Data, a); err != nil {
			return false, status.Errorf(codes.Internal, "failed to unmarshal article")
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

		aa = append(aa, a)

		if len(aa) == int(request.PageSize+1) {
			return false, nil
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
		return &ListArticlesResponse{
			Articles:      aa,
			NextPageToken: nextPageToken,
		}, nil
	default:
		log().Error(err)
		return nil, status.Errorf(codes.Internal, "failed to list articles")
	}
}

// CreateArticle creates a new article.
func (s *Service) CreateArticle(ctx context.Context, body io.Reader, articleURL *url.URL, published *time.Time) (*Article, error) {
	// before that date ksuid is no longer lexicographicaly sortable
	// https://github.com/segmentio/ksuid#how-do-they-work
	var timeLimit = time.Unix(1400000000, 0)

	createTime := time.Now()
	if published != nil {
		createTime = *published
	}

	if createTime.Before(timeLimit) {
		return nil, status.Errorf(codes.InvalidArgument, "articles written before %s not supported, sorry", timeLimit)
	}

	article := &Article{
		Url: articleURL.String(),
	}
	var err error
	article.CreateTime, err = ptypes.TimestampProto(createTime)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert timestamp")
	}

	var content []byte

	r, err := reader.NewFromReader(ctx, s.client, body, articleURL)
	if err == nil {
		article.Title = r.Title()
		article.SiteName = r.SiteName()
		article.IconUrl = r.IconURL()

		content = r.Content()
	}

	urlHash := murmur3.New128()
	_, _ = urlHash.Write([]byte(article.Url))

	// timestamp order == lexicographical order
	id, err := ksuid.FromParts(createTime, urlHash.Sum(nil))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate id")
	}

	actor, _ := actor.FromContext(ctx)
	name := actor.Child("articles", id.String())
	article.Name = name.String()

	contentHash := sha256.New()
	_, _ = contentHash.Write(content)
	article.ContentSha256Sum = fmt.Sprintf("%x", contentHash.Sum(nil))

	// if content exists
	if cc, err := s.storage.LoadAll(ctx, actor.Child("articles", "*").Child("content", fmt.Sprintf("%x", urlHash.Sum(nil)))); err == nil && len(cc) == 1 {
		// compare content
		if existingArticle, err := s.getArticle(ctx, cc[0].Name.Parent()); err == nil {
			if existingArticle.ContentSha256Sum == article.ContentSha256Sum {
				return nil, ErrAlreadyExists
			}
			existingArticle.ContentSha256Sum = article.ContentSha256Sum
			article = existingArticle
			name = resource.ParseName(existingArticle.Name)
			log().Infof("updating %s from '%s'", existingArticle.Name, existingArticle.Url)
		} else {
			log().Infof("creating %s from '%s'", name, existingArticle.Url)
		}
	}

	if err := s.storage.Store(ctx, name.Child("content", fmt.Sprintf("%x", urlHash.Sum(nil))), content); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to store the article content")
	}

	rawArticle, err := proto.Marshal(article)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to marshal the article")
	}

	if err := s.storage.Store(ctx, name, rawArticle); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to store the article")
	}

	article.Content = content
	return article, nil
}

// UpdateArticle updates the article.
func (s *Service) UpdateArticle(ctx context.Context, request *UpdateArticleRequest) (*Article, error) {
	name := resource.ParseName(request.Article.Name)

	if !actor.Owns(ctx, name) {
		return nil, status.Errorf(codes.PermissionDenied, "forbidden")
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
		return nil, status.Errorf(codes.Internal, "failed to marshal the article")
	}

	if err := s.storage.Store(ctx, name, rawArticle); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to store the article")
	}

	return article, nil
}

// GetArticle returns an article.
func (s *Service) GetArticle(ctx context.Context, request *GetArticleRequest) (*Article, error) {
	name := resource.ParseName(request.Name)

	if !actor.Owns(ctx, name) {
		return nil, status.Errorf(codes.PermissionDenied, "forbidden")
	}

	article, err := s.getArticle(ctx, name)
	if err != nil {
		return nil, err
	}
	if request.View != ArticleView_ARTICLE_VIEW_FULL {
		return article, nil
	}

	urlHash := murmur3.New128()
	_, _ = urlHash.Write([]byte(article.Url))

	article.Content, err = s.storage.Load(ctx, name.Child("content", fmt.Sprintf("%x", urlHash.Sum(nil))))
	switch {
	case err == nil:
	case errors.Is(err, storage.ErrNotFound):
		return article, nil
	default:
		return nil, status.Errorf(codes.Internal, "failed to load the article's content")
	}

	return article, nil
}

func (s *Service) getArticle(ctx context.Context, name *resource.Name) (*Article, error) {
	raw, err := s.storage.Load(ctx, name)
	switch {
	case err == nil:
	case errors.Is(err, storage.ErrNotFound):
		return nil, status.Errorf(codes.NotFound, "not found")
	default:
		return nil, status.Errorf(codes.Internal, "failed to load the article")
	}

	article := &Article{}
	if err := proto.Unmarshal(raw, article); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to unmarshal the article")
	}

	return article, nil
}

// DeleteArticle removes an article.
func (s *Service) DeleteArticle(ctx context.Context, request *DeleteArticleRequest) (*empty.Empty, error) {
	name := resource.ParseName(request.Name)

	if !actor.Owns(ctx, name) {
		return nil, status.Errorf(codes.PermissionDenied, "forbidden")
	}

	if err := s.storage.Delete(ctx, name); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete the article")
	}

	return &empty.Empty{}, nil
}

func log() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": "articles",
	})
}
