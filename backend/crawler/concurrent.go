package crawler

import (
	"context"
	"fmt"
	"net/url"
)

type crawler interface {
	Crawl(context.Context, *url.URL) ([]byte, error)
}

// ConcurrentCrawler limits number of concurrent requests.
type ConcurrentCrawler struct {
	crawler crawler

	sem chan struct{}
}

// WithConcurrencyLimit adds concurrency limit to the crawler.
func WithConcurrencyLimit(crawler crawler, limit int) *ConcurrentCrawler {
	return &ConcurrentCrawler{
		crawler: crawler,
		sem:     make(chan struct{}, limit),
	}
}

// Crawl downloads a page.
func (c *ConcurrentCrawler) Crawl(ctx context.Context, u *url.URL) ([]byte, error) {
	results := make(chan []byte)
	errors := make(chan error)

	go func(ctx context.Context, u *url.URL) {
		c.sem <- struct{}{}
		defer func() { <-c.sem }()
		res, err := c.crawler.Crawl(ctx, u)
		if err != nil {
			errors <- err
		} else {
			results <- res
		}
	}(ctx, u)

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("ctx cancelled")
	case err := <-errors:
		return nil, err
	case res := <-results:
		return res, nil
	}
}
