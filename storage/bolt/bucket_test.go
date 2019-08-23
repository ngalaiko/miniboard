package bolt // import "miniboard.app/storage/bolt_test"

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"miniboard.app/storage"
	"miniboard.app/storage/resource"
)

func Test_DB(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("Given a bucket", func(t *testing.T) {
		db := testBucket(ctx, t)

		t.Run("When data doesn't exist", func(t *testing.T) {
			t.Run("it should not be found", func(t *testing.T) {
				loaded, err := db.Load(resource.NewName("test", uuid.New().String()))
				assert.Empty(t, loaded)
				assert.Equal(t, errors.Cause(err), storage.ErrNotFound)
			})
		})

		t.Run("When root exists", func(t *testing.T) {
			name := resource.NewName("test", uuid.New().String())
			data := []byte("data")
			assert.NoError(t, db.Store(name, data))

			t.Run("It should be found", func(t *testing.T) {
				loaded, err := db.Load(name)
				if assert.NoError(t, err) {
					assert.Equal(t, loaded, data)
				}
			})

			t.Run("When child exists", func(t *testing.T) {
				name := name.Child("child", "id")
				data := []byte("data")
				assert.NoError(t, db.Store(name, data))

				t.Run("It should be found", func(t *testing.T) {
					loaded, err := db.Load(name)
					if assert.NoError(t, err) {
						assert.Equal(t, loaded, data)
					}
				})
			})
		})
	})
}

func testBucket(ctx context.Context, t *testing.T) storage.Storage {
	tmpfile, err := ioutil.TempFile("", "bolt")
	if err != nil {
		t.Fatalf("failed to create database: %s", err)
	}
	go func() {
		<-ctx.Done()
		defer os.Remove(tmpfile.Name()) // clean up
	}()

	db, err := New(ctx, tmpfile.Name())
	if err != nil {
		t.Fatalf("failed to create database: %s", err)
	}

	return db
}
