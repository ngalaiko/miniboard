package images

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"testing"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"miniboard.app/storage/redis"
)

func Test_Save__should_save_jpeg_image(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	h, err := miniredis.Run()
	assert.NoError(t, err)
	defer h.Close()

	db, err := redis.New(ctx, h.Addr())
	assert.NoError(t, err)

	s := NewService(db)

	file, err := os.Open("./testdata/image.jpeg")
	assert.NoError(t, err)

	_, err = s.Save(ctx, file)
	assert.NoError(t, err)
}

func Test_Save__should_save_png_image(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	h, err := miniredis.Run()
	assert.NoError(t, err)
	defer h.Close()

	db, err := redis.New(ctx, h.Addr())
	assert.NoError(t, err)

	s := NewService(db)

	file, err := os.Open("./testdata/image.png")
	assert.NoError(t, err)

	_, err = s.Save(ctx, file)
	assert.NoError(t, err)
}

func Test_Save__should_not_save_the_same_image_twice(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	h, err := miniredis.Run()
	assert.NoError(t, err)
	defer h.Close()

	db, err := redis.New(ctx, h.Addr())
	assert.NoError(t, err)

	s := NewService(db)

	file, err := ioutil.ReadFile("./testdata/image.jpeg")
	assert.NoError(t, err)

	i1, err := s.Save(ctx, bytes.NewBuffer(file))
	assert.NoError(t, err)

	i2, err := s.Save(ctx, bytes.NewBuffer(file))
	assert.NoError(t, err)

	assert.Equal(t, i1.String(), i2.String())
}
