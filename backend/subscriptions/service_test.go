package subscriptions

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/ngalaiko/miniboard/backend/items"
	"github.com/ngalaiko/miniboard/backend/tags"
)

func Test__Create_subscription_failed_to_parse_subscription(t *testing.T) {
	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	cw := &testCrawler{}
	cw.With("https://example.org", []byte("wrong"))
	service := NewService(sqldb, cw, &testLogger{}, nil, &testItemsService{})

	_, err := service.Create(ctx, "user id", mustParseURL("https://example.org"), []string{})
	if err != errFailedToParseSubscription {
		t.Fatalf("expected %s, got %s", errFailedToParseSubscription, err)
	}
}

func Test__Create_subscription_not_found(t *testing.T) {
	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	cw := &testCrawler{}
	service := NewService(sqldb, cw, &testLogger{}, nil, &testItemsService{})

	_, err := service.Create(ctx, "user id", mustParseURL("https://example.org"), []string{})
	if err != errFailedToDownloadSubscription {
		t.Fatalf("expected %s, got %s", errFailedToDownloadSubscription, err)
	}
}

func Test__Create(t *testing.T) {
	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	cw := &testCrawler{}
	service := NewService(sqldb, cw.With("https://example.org", subscriptionData), &testLogger{}, nil, &testItemsService{})

	subscription, err := service.Create(ctx, "user id", mustParseURL("https://example.org"), []string{})
	if err != nil {
		t.Fatalf("failed to create a subscription %s", err)
	}

	if subscription.UserID != "user id" {
		t.Errorf("user id expected: %s, got %s", "user id", subscription.UserID)
	}

	if subscription.Title != "Sample Subscription" {
		t.Errorf("title expected: %s, got %s", "Sample Subscription", subscription.Title)
	}

	if subscription.Updated != nil {
		t.Errorf("updated expected to be nil, got %+v", subscription.Updated)
	}

	if subscription.ID == "" {
		t.Errorf("id expected to not be empty")
	}

	if subscription.URL != "https://example.org" {
		t.Errorf("url should be https://example.org, got %s", subscription.URL)
	}
}

func Test__Create_twice(t *testing.T) {
	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	cw := &testCrawler{}
	cw = cw.With("https://example.org", subscriptionData)
	service := NewService(sqldb, cw, &testLogger{}, nil, &testItemsService{})

	_, err := service.Create(ctx, "user id", mustParseURL("https://example.org"), []string{})
	if err != nil {
		t.Fatalf("failed to create a subscription %s", err)
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
	service := NewService(sqldb, cw.With("https://example.org", subscriptionData), &testLogger{}, nil, &testItemsService{})

	tagService := tags.NewService(sqldb)
	tag, err := tagService.Create(ctx, "user id", "tag")
	if err != nil {
		t.Fatal(err)
	}

	subscription, err := service.Create(ctx, "user id", mustParseURL("https://example.org"), []string{tag.ID})
	if err != nil {
		t.Fatalf("failed to create a subscription: %s", err)
	}

	from, err := service.Get(ctx, "user id", subscription.ID)
	if err != nil {
		t.Fatalf("failed to get a subscription: %s", err)
	}

	if !cmp.Equal(subscription, from) {
		t.Error(cmp.Diff(subscription, from))
	}
}

func Test__Get_not_found(t *testing.T) {
	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	cw := &testCrawler{}
	service := NewService(sqldb, cw.With("https://example.org", subscriptionData), &testLogger{}, nil, &testItemsService{})

	_, err := service.Get(ctx, "user id", "id")
	if err != errNotFound {
		t.Errorf("expected %s, got %s", errNotFound, err)
	}
}

var subscriptionData = []byte(`<rss version="2.0">
<channel>
<title>Sample Subscription</title>
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

type testItemsService struct{}

func (*testItemsService) Create(context.Context, string, string, string, *time.Time, *string) (*items.Item, error) {
	return nil, nil
}
