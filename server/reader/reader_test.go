package reader

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"miniboard.app/images"
	"miniboard.app/storage"
	"miniboard.app/storage/bolt"
	"miniboard.app/storage/resource"
)

func Test(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	url, _ := url.Parse("http://example.com")

	r, err := NewFromReader(ctx, &http.Client{}, resource.NewName("articles", "1"), images.New(testDB(ctx, t)), testData(t), url)
	assert.NoError(t, err)

	title := r.Title()
	assert.Equal(t, "Building a peer to peer messenger", title)
	assert.Equal(t, "http://example.com/apple-touch-icon.png", r.IconURL())

	content := r.Content()
	assert.NotEmpty(t, content)
}

func testDB(ctx context.Context, t *testing.T) storage.Storage {
	tmpfile, err := ioutil.TempFile("", "bolt")
	if err != nil {
		t.Fatalf("failed to create database: %s", err)
	}
	go func() {
		<-ctx.Done()
		defer os.Remove(tmpfile.Name()) // clean up
	}()

	db, err := bolt.New(ctx, tmpfile.Name())
	if err != nil {
		t.Fatalf("failed to create database: %s", err)
	}
	return db
}

func testData(t *testing.T) *os.File {
	file, err := os.Open("./testdata/test.html")
	assert.NoError(t, err)

	return file
}
