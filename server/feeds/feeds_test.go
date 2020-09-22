package feeds

import (
	"bytes"
	"context"
	"database/sql"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/ngalaiko/miniboard/server/actor"
	"github.com/ngalaiko/miniboard/server/db"
	articles "github.com/ngalaiko/miniboard/server/genproto/articles/v1"
	"github.com/ngalaiko/miniboard/server/parsers"
	"github.com/stretchr/testify/assert"
)

func Test_feeds_Create_with_items(t *testing.T) {
	ctx := testContext()

	now := time.Now()
	parsedFeed := &parsers.Feed{
		Title: "feed",
		Items: []*parsers.Item{
			{
				UpdatedParsed: &now,
			},
		},
	}
	testArticles := &testArticles{}
	service := NewService(ctx, testDB(t), &testClient{}, testArticles, &testParser{feed: parsedFeed})

	url, _ := url.Parse("http://localhost")

	feed, err := service.CreateFeed(ctx, nil, url)
	if assert.NoError(t, err) {
		assert.NotEmpty(t, feed.Id)
		assert.NotEmpty(t, feed.LastFetched)
		assert.NotEmpty(t, feed.Url)
		assert.NotEmpty(t, feed.Title)
		assert.Equal(t, 1, testArticles.count)
	}
}

type testParser struct {
	feed *parsers.Feed
}

func (tp *testParser) Parse(feed io.Reader) (*parsers.Feed, error) {
	return tp.feed, nil
}

type testArticles struct {
	count int
}

func (ta *testArticles) CreateArticle(context.Context, io.Reader, *url.URL, *time.Time, *string) (*articles.Article, error) {
	ta.count++
	return nil, nil
}

type testClient struct{}

func (tc *testClient) Fetch(ctx context.Context, url string) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(&bytes.Buffer{}),
	}, nil
}

func testContext() context.Context {
	return actor.NewContext(context.Background(), "user_id")
}

func testDB(t *testing.T) *sql.DB {
	ctx := testContext()

	tmpFile, err := ioutil.TempFile(os.TempDir(), "testdb-")
	assert.NoError(t, err)

	t.Cleanup(func() {
		os.Remove(tmpFile.Name())
	})

	sqlite, err := db.NewSQLite(tmpFile.Name())
	assert.NoError(t, err)
	assert.NoError(t, db.Migrate(ctx, sqlite))

	return sqlite
}
