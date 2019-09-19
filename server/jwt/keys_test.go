package jwt

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"miniboard.app/storage"
	"miniboard.app/storage/bolt"
)

func Test_keyStorage_Create__should_create_new_key(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := newKeyStorage(testDB(ctx, t))

	key, err := service.Create()
	assert.NoError(t, err)
	assert.NotNil(t, key)
	assert.NotNil(t, key.Private)
	assert.NotNil(t, key.Public)
}

func Test_keyStorage_Get(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := newKeyStorage(testDB(ctx, t))

	t.Run("When a key doesn't exist", func(t *testing.T) {
		t.Run("Then error should be not found", func(t *testing.T) {
			key, err := service.Get("random id")
			assert.Nil(t, key)
			assert.Equal(t, errors.Cause(err), storage.ErrNotFound)
		})
	})

	t.Run("When a creating a key", func(t *testing.T) {
		key, err := service.Create()
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
				fromStorage, err := service.Get(key.ID)
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
			_, err := service.Create()
			assert.NoError(t, err)
			t.Run("Both should be listed", func(t *testing.T) {
				keys, err := service.List()
				assert.NoError(t, err)
				assert.Len(t, keys, 2)

				assert.NotEqual(t, keys[0].ID, keys[1].ID)
			})
		})
	})
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
