package feeds

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ngalaiko/miniboard/backend/tags"
)

func Test__Create_feed_failed_to_parse_feed(t *testing.T) {
	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	cw := &testCrawler{}
	cw.With("https://example.org", []byte("wrong"))
	service := NewService(sqldb, cw, &testLogger{})

	_, err := service.Create(ctx, "user id", mustParseURL("https://example.org"), []string{})
	if err != errFailedToParseFeed {
		t.Fatalf("expected %s, got %s", errFailedToParseFeed, err)
	}
}

func Test__Create_feed_not_found(t *testing.T) {
	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	cw := &testCrawler{}
	service := NewService(sqldb, cw, &testLogger{})

	_, err := service.Create(ctx, "user id", mustParseURL("https://example.org"), []string{})
	if err != errFailedToDownloadFeed {
		t.Fatalf("expected %s, got %s", errFailedToDownloadFeed, err)
	}
}

func Test__Create(t *testing.T) {
	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	cw := &testCrawler{}
	service := NewService(sqldb, cw.With("https://example.org", feedData), &testLogger{})

	feed, err := service.Create(ctx, "user id", mustParseURL("https://example.org"), []string{})
	if err != nil {
		t.Fatalf("failed to create a feed %s", err)
	}

	if feed.UserID != "user id" {
		t.Errorf("user id expected: %s, got %s", "user id", feed.UserID)
	}

	if feed.Title != "Sample Feed" {
		t.Errorf("title expected: %s, got %s", "Sample Feed", feed.Title)
	}

	if feed.Updated != nil {
		t.Errorf("updated expected to be nil, got %+v", feed.Updated)
	}

	if feed.ID == "" {
		t.Errorf("id expected to not be empty")
	}

	if feed.URL != "https://example.org" {
		t.Errorf("url should be https://example.org, got %s", feed.URL)
	}
}

func Test__Create_twice(t *testing.T) {
	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	cw := &testCrawler{}
	cw = cw.With("https://example.org", feedData)
	service := NewService(sqldb, cw, &testLogger{})

	_, err := service.Create(ctx, "user id", mustParseURL("https://example.org"), []string{})
	if err != nil {
		t.Fatalf("failed to create a feed %s", err)
	}

	_, secondErr := service.Create(ctx, "user id", mustParseURL("https://example.org"), []string{})
	if secondErr != errAlreadyExists {
		t.Fatalf("expected %s, got %s", errAlreadyExists, secondErr)
	}
}

func Test__Get(t *testing.T) {
	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	cw := &testCrawler{}
	service := NewService(sqldb, cw.With("https://example.org", feedData), &testLogger{})

	tagService := tags.NewService(sqldb)
	tag, err := tagService.Create(ctx, "user id", "tag")
	if err != nil {
		t.Fatal(err)
	}

	feed, err := service.Create(ctx, "user id", mustParseURL("https://example.org"), []string{tag.ID})
	if err != nil {
		t.Fatalf("failed to create a feed: %s", err)
	}

	from, err := service.Get(ctx, "user id", feed.ID)
	if err != nil {
		t.Fatalf("failed to get a feed: %s", err)
	}

	if !cmp.Equal(feed, from) {
		t.Error(cmp.Diff(feed, from))
	}
}

func Test__Get_not_found(t *testing.T) {
	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	cw := &testCrawler{}
	service := NewService(sqldb, cw.With("https://example.org", feedData), &testLogger{})

	_, err := service.Get(ctx, "user id", "id")
	if err != errNotFound {
		t.Errorf("expected %s, got %s", errNotFound, err)
	}
}

var feedData = []byte(`<rss version="2.0">
<channel>
<title>Sample Feed</title>
</channel>
</rss>`)

type testCrawler struct {
	content map[string][]byte
}

// With adds a mock content.
func (tc *testCrawler) With(url string, content []byte) *testCrawler {
	if tc.content == nil {
		tc.content = map[string][]byte{}
	}
	tc.content[url] = content
	return tc
}

// Crawl inplements crawler.
func (tc *testCrawler) Crawl(_ context.Context, url *url.URL) ([]byte, error) {
	data, found := tc.content[url.String()]
	if !found {
		return nil, fmt.Errorf("not found")
	}

	return data, nil
}

func mustParseURL(raw string) *url.URL {
	url, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	return url
}
