package articles

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/golang/protobuf/ptypes"
	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/ngalaiko/miniboard/server/actor"
	"github.com/ngalaiko/miniboard/server/fetch"
	"github.com/ngalaiko/miniboard/server/reader"
	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
	"github.com/spaolacci/murmur3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// FromID returns page token converted to id.
func (lr *ListArticlesRequest) FromID() (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(lr.PageToken)
	if err != nil {
		return "", status.Errorf(codes.InvalidArgument, "invalid page token")
	}
	return string(decoded), nil
}

// Service controls articles resource.
type Service struct {
	storage *articlesDB
	client  fetch.Fetcher
}

// NewService returns a new articles service instance.
func NewService(sqldb *sql.DB, fetcher fetch.Fetcher) *Service {
	return &Service{
		storage: newDB(sqldb),
		client:  fetcher,
	}
}

// ListArticles returns a list of articles.
func (s *Service) ListArticles(ctx context.Context, request *ListArticlesRequest) (*ListArticlesResponse, error) {
	request.PageSize++
	aa, err := s.storage.List(ctx, request)

	var nextPageToken string
	if len(aa) == int(request.PageSize) {
		nextPageToken = base64.StdEncoding.EncodeToString([]byte(aa[len(aa)-1].Id))
		aa = aa[:request.PageSize-1]
	}

	switch err {
	case nil, sql.ErrNoRows:
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
func (s *Service) CreateArticle(
	ctx context.Context,
	body io.Reader,
	articleURL *url.URL,
	published *time.Time,
	sourceID string,
) (*Article, error) {
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

	actor, _ := actor.FromContext(ctx)

	article := &Article{
		Url:    articleURL.String(),
		UserId: actor.ID(),
	}

	article.FeedId = sourceID

	var err error
	article.CreateTime, err = ptypes.TimestampProto(createTime)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert timestamp")
	}

	var content []byte

	r, err := reader.NewFromReader(ctx, s.client, body, articleURL)
	if err == nil {
		article.Title = r.Title()
		content = r.Content()
	}

	urlHash := murmur3.New128()
	_, _ = urlHash.Write([]byte(article.Url))

	// timestamp order == lexicographical order
	id, err := ksuid.FromParts(createTime, urlHash.Sum(nil))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate id")
	}

	article.Id = id.String()

	contentHash := sha256.New()
	_, _ = contentHash.Write(content)
	article.ContentSha256 = fmt.Sprintf("%x", contentHash.Sum(nil))
	article.Content = content

	// if content exists
	if existing, err := s.storage.GetByUserIDUrl(ctx, article.UserId, article.Url); err == nil && existing != nil {
		if existing.ContentSha256 == article.ContentSha256 {
			return existing, nil
		}
		existing.ContentSha256 = article.ContentSha256
		existing.Content = article.Content

		if err := s.storage.Update(ctx, existing, actor.ID()); err != nil {
			log().Error(err)
			return nil, status.Errorf(codes.Internal, "failed to store the article")
		}

		return existing, nil
	}

	if err := s.storage.Create(ctx, article); err != nil {
		log().Error(err)
		return nil, status.Errorf(codes.Internal, "failed to store the article")
	}

	return article, nil
}

// UpdateArticle updates the article.
func (s *Service) UpdateArticle(ctx context.Context, request *UpdateArticleRequest) (*Article, error) {
	a, _ := actor.FromContext(ctx)

	article, err := s.getArticle(ctx, request.Article.Id, a.ID())
	if err != nil {
		return nil, err
	}

	var updated bool

	for _, path := range request.UpdateMask.GetPaths() {
		switch path {
		case "is_read_eq":
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

	if err := s.storage.Update(ctx, article, a.ID()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to store the article")
	}

	return article, nil
}

// GetArticle returns an article.
func (s *Service) GetArticle(ctx context.Context, request *GetArticleRequest) (*Article, error) {
	a, _ := actor.FromContext(ctx)

	article, err := s.getArticle(ctx, request.Id, a.ID())
	if err != nil {
		return nil, err
	}

	if request.View != ArticleView_ARTICLE_VIEW_FULL {
		article.Content = nil
	}

	return article, nil
}

func (s *Service) getArticle(ctx context.Context, id string, userID string) (*Article, error) {
	article, err := s.storage.Get(ctx, id, userID)
	switch {
	case err == nil:
		return article, nil
	case errors.Is(err, sql.ErrNoRows):
		return nil, status.Errorf(codes.NotFound, "not found")
	default:
		return nil, status.Errorf(codes.Internal, "failed to load the article")
	}
}

// DeleteArticle removes an article.
func (s *Service) DeleteArticle(ctx context.Context, request *DeleteArticleRequest) (*empty.Empty, error) {
	a, _ := actor.FromContext(ctx)
	if err := s.storage.Delete(ctx, request.Id, a.ID()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete the article")
	}

	return &empty.Empty{}, nil
}

func log() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": "articles",
	})
}
