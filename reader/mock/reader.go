package mock

import (
	"net/url"

	"miniboard.app/reader"
)

type mock struct{}

// New returns mock reader.
func New(*url.URL) (reader.Reader, error) {
	return &mock{}, nil
}

// Title returns title.
func (m *mock) Title() string {
	return "Title"
}

// IconURL returns icon urls.
func (m *mock) IconURL() string {
	return "http://example.com/icon.png"
}

// Content returns page content.
func (m *mock) Content() []byte {
	return []byte("page!")
}
