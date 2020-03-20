package fetch

import (
	"context"
	"net/http"
)

// Fetcher downloads content from the internet.
type Fetcher interface {
	Fetch(context.Context, string) (*http.Response, error)
}

// HTTPFetcher downloads content from the internet.
type HTTPFetcher struct {
	client *http.Client
}

// New returns new http fetcher.
func New() *HTTPFetcher {
	return &HTTPFetcher{
		client: &http.Client{},
	}
}

// Fetch returns http response by url.
func (hf *HTTPFetcher) Fetch(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return hf.client.Do(req)
}
