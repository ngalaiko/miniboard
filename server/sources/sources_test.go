package sources

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	articles "miniboard.app/proto/users/articles/v1"
	feeds "miniboard.app/proto/users/feeds/v1"
	sources "miniboard.app/proto/users/sources/v1"
)

type testClient struct {
	typ string
}

func (tc *testClient) Get(url string) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(&bytes.Buffer{}),
		Header: map[string][]string{
			"Content-Type": {tc.typ},
		},
	}, nil
}

func Test_sources(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("With sources service", func(t *testing.T) {
		articles := &mockArticles{}
		feeds := &mockFeeds{}
		service := New(articles, feeds)

		t.Run("When creating a source from html page", func(t *testing.T) {
			service.client = &testClient{typ: "text/html"}
			source, err := service.CreateSource(ctx, &sources.CreateSourceRequest{
				Source: &sources.Source{
					Url: "http://example.com",
				},
			})
			assert.NoError(t, err)
			assert.Equal(t, "http://example.com", source.Url)

			t.Run("Should create an article", func(t *testing.T) {
				assert.Equal(t, len(articles.articles), 1)
			})
		})

		t.Run("When creating a source from rss page", func(t *testing.T) {
			service.client = &testClient{typ: "application/rss+xml"}
			source, err := service.CreateSource(ctx, &sources.CreateSourceRequest{
				Source: &sources.Source{
					Url: "http://example.com",
				},
			})

			if assert.NoError(t, err) {
				assert.Equal(t, "http://example.com", source.Url)
			}

			t.Run("Should create a feed", func(t *testing.T) {
				assert.Equal(t, len(feeds.feeds), 1)
			})
		})

		t.Run("When creating a source from unknown page", func(t *testing.T) {
			service.client = &testClient{typ: "something else"}
			_, err := service.CreateSource(ctx, &sources.CreateSourceRequest{
				Source: &sources.Source{
					Url: "http://example.com",
				},
			})
			assert.Error(t, err)
		})
	})
}

type mockFeeds struct {
	feeds []*feeds.Feed
}

func (s *mockFeeds) CreateFeed(context.Context, io.Reader, *url.URL) (*feeds.Feed, error) {
	feed := &feeds.Feed{}
	s.feeds = append(s.feeds, feed)
	return feed, nil
}

type mockArticles struct {
	articles []*articles.Article
}

func (s *mockArticles) CreateArticle(ctx context.Context, body io.Reader, url *url.URL, _ *time.Time) (*articles.Article, error) {
	s.articles = append(s.articles, &articles.Article{
		Url: url.String(),
	})
	return &articles.Article{}, nil
}

func (s *mockArticles) ListArticles(ctx context.Context, request *articles.ListArticlesRequest) (*articles.ListArticlesResponse, error) {
	return nil, nil
}

func (s *mockArticles) UpdateArticle(ctx context.Context, request *articles.UpdateArticleRequest) (*articles.Article, error) {
	return nil, nil
}

func (s *mockArticles) GetArticle(ctx context.Context, request *articles.GetArticleRequest) (*articles.Article, error) {
	return nil, nil
}

func (s *mockArticles) DeleteArticle(ctx context.Context, request *articles.DeleteArticleRequest) (*empty.Empty, error) {
	return nil, nil
}
