package http

import (
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	url, _ := url.Parse("http://example.com")

	r, err := NewFromReader(testData(t), url)
	assert.NoError(t, err)

	title := r.Title()
	assert.Equal(t, "Building a peer to peer messenger", title)

	iconURL := r.IconURL()
	assert.Contains(t, iconURL, "http://example.com/favicon-16x16.png")
	assert.Contains(t, iconURL, "http://example.com/favicon-32x32.png")
	assert.Contains(t, iconURL, "http://example.com/apple-touch-icon.png")
	assert.Contains(t, iconURL, "http://example.com/favicon.ico")

	content := r.Content()
	assert.NotEmpty(t, content)
}

func testData(t *testing.T) *os.File {
	file, err := os.Open("./testdata/test.html")
	assert.NoError(t, err)

	return file
}
