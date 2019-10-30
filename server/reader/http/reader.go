package http

import (
	"io"
	"net/http"
	"net/url"

	"github.com/go-shiori/go-readability"
	"github.com/pkg/errors"
	"miniboard.app/reader"
)

// Reader returns simplified HTML content.
type Reader struct {
	article *readability.Article
	url     *url.URL
}

// New erturns new reader from a url.
func New(url *url.URL) (reader.Reader, error) {
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch url")
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return NewFromReader(resp.Body, url)
}

// NewFromReader returns new reader from io.Reader.
// URL is needed to form complete links to images.
func NewFromReader(raw io.Reader, url *url.URL) (*Reader, error) {
	article, err := readability.FromReader(raw, url.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse document")
	}
	return &Reader{
		article: &article,
		url:     url,
	}, nil
}

// Title returns the page title.
func (r *Reader) Title() (title string) {
	return r.article.Title
}

// Content returns page content.
func (r *Reader) Content() []byte {
	return []byte(r.article.Content)
}

// IconURL returns a link to the first page favicon.
func (r *Reader) IconURL() string {
	return r.article.Favicon
}
