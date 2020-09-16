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

	"github.com/ngalaiko/miniboard/server/articles"
	"github.com/ngalaiko/miniboard/server/feeds"
)

func Test_Create(t *testing.T) {

}

type testClient struct {
	typ string
}

func (tc *testClient) Fetch(ctx context.Context, url string) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(&bytes.Buffer{}),
		Header: map[string][]string{
			"Content-Type": {tc.typ},
		},
	}, nil
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

func (s *mockArticles) CreateArticle(ctx context.Context, body io.Reader, url *url.URL, _ *time.Time, _ *string) (*articles.Article, error) {
	s.articles = append(s.articles, &articles.Article{
		Url: url.String(),
	})
	return &articles.Article{}, nil
}
