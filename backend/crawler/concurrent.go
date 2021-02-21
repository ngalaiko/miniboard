package crawler

import (
	"context"
	"fmt"
	"net/url"
)

// ConcurrentCrawler limits number of concurrent requests.
type ConcurrentCrawler struct {
	crawler *Crawler

	sem chan struct{}
}

// WithConcurrencyLimit adds concurrency limit to the crawler.
func WithConcurrencyLimit(crawler *Crawler, limit int) *ConcurrentCrawler {
	return &ConcurrentCrawler{
		crawler: crawler,
		sem:     make(chan struct{}, limit),
	}
}

// Crawl downloads a page.
func (c *ConcurrentCrawler) Crawl(ctx context.Context, u *url.URL) ([]byte, error) {
	results := make(chan []byte)
	errors := make(chan error)

	c.sem <- struct{}{}
	go func(ctx context.Context, u *url.URL) {
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
