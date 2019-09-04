package http

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"miniboard.app/reader"
)

// Reader returns simplified HTML content.
type Reader struct {
	doc *Document
	url *url.URL
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
	content, err := ioutil.ReadAll(raw)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read body")
	}
	doc, err := NewDocument(content, url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse document")
	}
	return &Reader{
		doc: doc,
		url: url,
	}, nil
}

// Title returns the page title.
func (r *Reader) Title() (title string) {
	return r.doc.Title()
}

// Content returns page content.
func (r *Reader) Content() []byte {
	return r.doc.Content()
}

// IconURL returns a link to the first page favicon.
func (r *Reader) IconURL() (iconURLs []string) {
	return r.doc.IconURL()
}
