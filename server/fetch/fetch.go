package fetch

import (
	"context"
	"net/http"
)

type Fetcher interface {
	Fetch(context.Context, string) (*http.Response, error)
}

type HTTPFetcher struct {
	client *http.Client
}

// New returns new http fetcher.
func New() *HTTPFetcher {
	return &HTTPFetcher{
		client: &http.Client{},
	}
}

func (hf *HTTPFetcher) Fetch(ctx context.Context, url string) (*http.Response, error) {
	return hf.client.Get(url)
}
