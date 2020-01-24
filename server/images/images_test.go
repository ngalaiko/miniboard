package images

import (
	"context"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"miniboard.app/storage"
	"miniboard.app/storage/bolt"
	"miniboard.app/storage/resource"
)

func Test_Save__should_save_jpeg_image(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := New(testDB(ctx, t))

	file, err := os.Open("./testdata/image.jpeg")
	assert.NoError(t, err)

	name, err := s.Save(ctx, resource.NewName("article", "1"), file)
	assert.NoError(t, err)
	assert.True(t, strings.HasPrefix(name.String(), "article/1"))
}

func Test_Save__should_save_png_image(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := New(testDB(ctx, t))

	file, err := os.Open("./testdata/image.png")
	assert.NoError(t, err)

	name, err := s.Save(ctx, resource.NewName("article", "1"), file)
	assert.NoError(t, err)
	assert.True(t, strings.HasPrefix(name.String(), "article/1"))
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
