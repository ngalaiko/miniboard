package images

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"miniboard.app/storage/redis"
)

func Test_Save__should_save_jpeg_image(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	host := os.Getenv("REDIS_HOST")
	if host == "" {
		t.Skip("REDIS_HOST is not set")
	}

	db, err := redis.New(ctx, host)
	assert.NoError(t, err)

	s := New(db)

	file, err := os.Open("./testdata/image.jpeg")
	assert.NoError(t, err)

	_, err = s.Save(ctx, file)
	assert.NoError(t, err)
}

func Test_Save__should_save_png_image(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	host := os.Getenv("REDIS_HOST")
	if host == "" {
		t.Skip("REDIS_HOST is not set")
	}

	db, err := redis.New(ctx, host)
	assert.NoError(t, err)

	s := New(db)

	file, err := os.Open("./testdata/image.png")
	assert.NoError(t, err)

	_, err = s.Save(ctx, file)
	assert.NoError(t, err)
}
