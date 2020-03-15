package reader

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"

	"github.com/go-shiori/go-readability"
	"golang.org/x/net/html"
	"miniboard.app/fetch"
)

// Reader returns simplified HTML content.
type Reader struct {
	article *readability.Article
	url     *url.URL
	content []byte
}

// NewFromReader returns new reader from io.Reader.
func NewFromReader(ctx context.Context, client fetch.Fetcher, raw io.Reader, url *url.URL) (*Reader, error) {
	article, err := readability.FromReader(raw, url.String())
	if err != nil {
		return nil, fmt.Errorf("failed to parse document: %w", err)
	}

	b := &bytes.Buffer{}
	if err := html.Render(b, article.Node); err != nil {
		return nil, err
	}

	return &Reader{
		article: &article,
		url:     url,
		content: b.Bytes(),
	}, nil
}

// Title returns the page title.
func (r *Reader) Title() (title string) {
	return r.article.Title
}

// SiteName returns name of source website.
func (r *Reader) SiteName() string {
	return r.article.SiteName
}

// Content returns page content.
func (r *Reader) Content() []byte {
	return r.content
}

// IconURL returns a link to the first page favicon.
func (r *Reader) IconURL() string {
	return r.article.Favicon
}
