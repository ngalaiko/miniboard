package jwt

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"miniboard.app/storage"
	"miniboard.app/storage/redis"
)

func Test_keyStorage(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	host := os.Getenv("REDIS_HOST")
	if host == "" {
		t.Skip("REDIS_HOST is not set")
	}

	db, err := redis.New(ctx, host)
	assert.NoError(t, err)

	service := newKeyStorage(db)

	t.Run("When a key doesn't exist", func(t *testing.T) {
		t.Run("Then error should be not found", func(t *testing.T) {
			key, err := service.Get(ctx, "random id")
			assert.Nil(t, key)
			assert.True(t, errors.Is(err, storage.ErrNotFound))
		})
	})

	t.Run("When a creating a key", func(t *testing.T) {
		key, err := service.Create(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, key)
		assert.NotNil(t, key.Private)
		assert.NotNil(t, key.Public)

		t.Run("It should be added to the cache", func(t *testing.T) {
			found, ok := service.cache.Get(key.ID)
			assert.True(t, ok)
			assert.Equal(t, key, found)
		})

		t.Run("If key is not cached", func(t *testing.T) {
			service.cache.store.Delete(key.ID)

			t.Run("Then it should be returned", func(t *testing.T) {
				fromStorage, err := service.Get(ctx, key.ID)
				assert.NoError(t, err)
				assert.Equal(t, fromStorage, key)
			})

			t.Run("And cached", func(t *testing.T) {
				found, ok := service.cache.Get(key.ID)
				assert.True(t, ok)
				assert.Equal(t, key, found)
			})
		})
		t.Run("When creating another key", func(t *testing.T) {
			_, err := service.Create(ctx)
			assert.NoError(t, err)
			t.Run("Both should be listed", func(t *testing.T) {
				keys, err := service.List(ctx)
				assert.NoError(t, err)
				assert.Len(t, keys, 2)

				assert.NotEqual(t, keys[0].ID, keys[1].ID)
			})
		})
	})
}
