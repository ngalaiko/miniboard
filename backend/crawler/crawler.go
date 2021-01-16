package crawler

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const userAgent = "miniboard-crawler"

// Crawler crawls web and downloads stuff.
type Crawler struct {
	httpClient *http.Client
}

// New returns a new crawler.
func New() *Crawler {
	return &Crawler{
		httpClient: &http.Client{},
	}
}

// Crawl downloads a page.
func (c *Crawler) Crawl(ctx context.Context, url *url.URL) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request")
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send a request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch page: status %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body")
	}

	return body, nil
}
