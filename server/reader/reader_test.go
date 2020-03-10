package reader

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"testing"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"miniboard.app/images"
	"miniboard.app/storage/redis"
)

func Test(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	url, _ := url.Parse("http://example.com")

	s, err := miniredis.Run()
	assert.NoError(t, err)
	defer s.Close()

	db, err := redis.New(ctx, s.Addr())
	assert.NoError(t, err)

	r, err := NewFromReader(ctx, &http.Client{}, images.New(db), testData(t), url)
	assert.NoError(t, err)

	title := r.Title()
	assert.Equal(t, "Building a peer to peer messenger", title)
	assert.Equal(t, "http://example.com/apple-touch-icon.png", r.IconURL())

	content := r.Content()
	assert.NotEmpty(t, content)
}

func testData(t *testing.T) *os.File {
	file, err := os.Open("./testdata/test.html")
	assert.NoError(t, err)

	return file
}
