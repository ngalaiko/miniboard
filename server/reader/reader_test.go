package reader

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testClient struct{}

func (tc *testClient) Fetch(ctx context.Context, url string) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBuffer([]byte{})),
	}, nil
}

func Test(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	url, _ := url.Parse("http://example.com")

	r, err := NewFromReader(ctx, &testClient{}, testData(t), url)
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
