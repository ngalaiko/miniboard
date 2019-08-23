package jwt // import "miniboard.app/jwt"

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/uuid"
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
			key, err := service.Get(uuid.New())
			assert.Nil(t, key)
			assert.Equal(t, errors.Cause(err), storage.ErrNotFound)
		})
	})

	t.Run("When a key exists", func(t *testing.T) {
		key, err := service.Create()
		assert.NoError(t, err)

		t.Run("Then it should be returned", func(t *testing.T) {
			fromStorage, err := service.Get(key.ID)
			assert.NoError(t, err)

			assert.Equal(t, fromStorage, key)
		})

	})

	key, err := service.Create()
	assert.NoError(t, err)
	assert.NotNil(t, key)
	assert.NotNil(t, key.Private)
	assert.NotNil(t, key.Public)
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
