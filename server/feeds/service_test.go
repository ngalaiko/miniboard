package feeds

import (
	"context"
	"io"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"miniboard.app/api/actor"
	"miniboard.app/articles"
	"miniboard.app/storage/redis"
	"miniboard.app/storage/resource"
)

func Test_feeds(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = actor.NewContext(ctx, resource.NewName("users", "test"))

	s, err := miniredis.Run()
	assert.NoError(t, err)
	defer s.Close()

	db, err := redis.New(ctx, s.Addr())
	assert.NoError(t, err)

	testFeed, err := os.Open("./testdata/feed.xml")
	assert.NoError(t, err)

	testURL, err := url.Parse("https://meduza.io/rss/all")
	assert.NoError(t, err)

	t.Run("With feeds service", func(t *testing.T) {
		articlesService := &mockArticles{}
		service := NewService(ctx, db, articlesService)

		t.Run("When adding a feed", func(t *testing.T) {
			feed, err := service.CreateFeed(ctx, testFeed, testURL)
			t.Run("It should be added", func(t *testing.T) {
				if assert.NoError(t, err) {
					assert.NotEmpty(t, feed.Name)
					assert.Empty(t, feed.LastFetched)
					assert.NotEmpty(t, feed.Url)
				}
			})

			t.Run("Eventually, articles must be fetched", func(t *testing.T) {
				for {
					select {
					case <-time.Tick(time.Second):
						assert.Len(t, articlesService.articles, 30)
						return
					default:
						articlesService.RLock()
						if len(articlesService.articles) == 30 {
							articlesService.RUnlock()
							return
						}
						articlesService.RUnlock()
					}
				}
			})

			t.Run("When adding the same feed again", func(t *testing.T) {
				feed2, err := service.CreateFeed(ctx, testFeed, testURL)

				t.Run("It must not be duplicated", func(t *testing.T) {
					if assert.NoError(t, err) {
						assert.Equal(t, feed.Name, feed2.Name)
					}
				})
			})
		})
	})
}

type mockArticles struct {
	sync.RWMutex

	articles []*articles.Article
}

func (s *mockArticles) CreateArticle(ctx context.Context, body io.Reader, url *url.URL, _ *time.Time) (*articles.Article, error) {
	s.Lock()
	defer s.Unlock()

	s.articles = append(s.articles, &articles.Article{
		Url: url.String(),
	})
	return &articles.Article{}, nil
}