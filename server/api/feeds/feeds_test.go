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
	"github.com/ngalaiko/miniboard/server/api/feeds/parsers"
	"github.com/ngalaiko/miniboard/server/db"
	articles "github.com/ngalaiko/miniboard/server/genproto/articles/v1"
	"github.com/stretchr/testify/assert"
)

func Test_Create_no_items(t *testing.T) {
	ctx := testContext()

	parsedFeed := &parsers.Feed{
		Title: "feed",
	}
	testArticles := &testArticles{}
	service := NewService(ctx, &testLogger{}, testDB(t), &testClient{}, testArticles, &testParser{feed: parsedFeed})

	now := time.Now().UTC()
	service.nowFunc = func() time.Time { return now }

	url, _ := url.Parse("http://localhost")

	feed, err := service.CreateFeed(ctx, nil, url)
	if assert.NoError(t, err) {
		assert.NotEmpty(t, feed.Id)
		assert.Equal(t, now, feed.LastFetched.AsTime())
		assert.NotEmpty(t, feed.Url)
		assert.NotEmpty(t, feed.Title)
		assert.Equal(t, 0, testArticles.count)
	}
}

func Test_Create_with_items(t *testing.T) {
	ctx := testContext()

	now := time.Now().UTC()
	parsedFeed := &parsers.Feed{
		Title: "feed",
		Items: []*parsers.Item{
			{
				PublishedParsed: &now,
			},
		},
	}
	testArticles := &testArticles{}
	service := NewService(ctx, &testLogger{}, testDB(t), &testClient{}, testArticles, &testParser{feed: parsedFeed})

	service.nowFunc = func() time.Time { return now }

	url, _ := url.Parse("http://localhost")

	feed, err := service.CreateFeed(ctx, nil, url)
	if assert.NoError(t, err) {
		assert.NotEmpty(t, feed.Id)
		assert.Equal(t, now, feed.LastFetched.AsTime())
		assert.NotEmpty(t, feed.Url)
		assert.NotEmpty(t, feed.Title)
		assert.Equal(t, 1, testArticles.count)
	}
}

func Test_update_nothing_updated(t *testing.T) {
	ctx := testContext()

	now := time.Now().UTC()
	hourAgo := now.Add(-1 * time.Hour)
	dayAgo := now.Add(-1 * 24 * time.Hour)
	parsedFeed := &parsers.Feed{
		Title: "feed",
		Items: []*parsers.Item{
			{
				PublishedParsed: &dayAgo,
			},
		},
	}

	testArticles := &testArticles{}
	testParser := &testParser{feed: parsedFeed}
	service := NewService(ctx, &testLogger{}, testDB(t), &testClient{}, testArticles, testParser)

	service.nowFunc = func() time.Time { return hourAgo }
	service.updateLeeway = time.Duration(0)

	url, _ := url.Parse("http://localhost")

	feed, err := service.CreateFeed(ctx, nil, url)
	if assert.NoError(t, err) {
		assert.Equal(t, 1, testArticles.count)
		assert.Equal(t, hourAgo, feed.LastFetched.AsTime())
	}

	service.nowFunc = func() time.Time { return now }
	if assert.NoError(t, service.updateFeed(ctx, feed)) {
		assert.Equal(t, 1, testArticles.count)
		assert.Equal(t, now, feed.LastFetched.AsTime())
	}
}

func Test_update_item_added(t *testing.T) {
	ctx := testContext()

	now := time.Now().UTC()
	dayAgo := now.Add(-1 * 24 * time.Hour)
	threeHoursAgo := now.Add(-3 * time.Hour)
	twoHoursAgo := now.Add(-2 * time.Hour)

	parsedFeed := &parsers.Feed{
		Title: "feed",
		Items: []*parsers.Item{
			{
				PublishedParsed: &dayAgo,
			},
		},
	}

	testArticles := &testArticles{}
	testParser := &testParser{feed: parsedFeed}
	service := NewService(ctx, &testLogger{}, testDB(t), &testClient{}, testArticles, testParser)
	service.nowFunc = func() time.Time {
		return threeHoursAgo
	}

	url, _ := url.Parse("http://localhost")

	feed, err := service.CreateFeed(ctx, nil, url)
	if assert.NoError(t, err) {
		assert.Equal(t, 1, testArticles.count)
	}

	service.nowFunc = func() time.Time {
		return now
	}

	testParser.feed.Items = append(testParser.feed.Items, &parsers.Item{
		PublishedParsed: &twoHoursAgo,
	})

	if assert.NoError(t, service.updateFeed(ctx, feed)) {
		assert.Equal(t, 3, testArticles.count)
	}
}

func Test_update_item_added_with_no_time(t *testing.T) {
	ctx := testContext()

	now := time.Now()
	threeHoursAgo := now.Add(-3 * time.Hour)

	parsedFeed := &parsers.Feed{
		Title: "feed",
		Items: []*parsers.Item{
			{},
		},
	}

	testArticles := &testArticles{}
	testParser := &testParser{feed: parsedFeed}
	service := NewService(ctx, &testLogger{}, testDB(t), &testClient{}, testArticles, testParser)
	service.updateLeeway = time.Duration(0)
	service.nowFunc = func() time.Time {
		return threeHoursAgo
	}

	url, _ := url.Parse("http://localhost")

	feed, err := service.CreateFeed(ctx, nil, url)
	if assert.NoError(t, err) {
		assert.Equal(t, 1, testArticles.count)
	}

	service.nowFunc = func() time.Time {
		return now
	}

	testParser.feed.Items = append(testParser.feed.Items, &parsers.Item{})

	if assert.NoError(t, service.updateFeed(ctx, feed)) {
		assert.Equal(t, 3, testArticles.count)
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

	sqlite, err := db.New(ctx, &db.Config{
		Driver: "sqlite3",
		Addr:   tmpFile.Name(),
	}, &testLogger{})
	assert.NoError(t, err)

	return sqlite
}

type testLogger struct{}

func (l *testLogger) Error(string, ...interface{}) {}

func (l *testLogger) Info(string, ...interface{}) {}
