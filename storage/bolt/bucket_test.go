package bolt // import "miniboard.app/storage/bolt_test"

import (
	"context"
	"fmt"
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

		t.Run("When a few elements exist", func(t *testing.T) {
			for i := 0; i < 10; i++ {
				name := resource.NewName("test", fmt.Sprint(i))
				data := []byte(fmt.Sprintf("data %d", i))
				assert.NoError(t, db.Store(name, data))
			}

			t.Run("When loading all elements", func(t *testing.T) {
				name := resource.NewName("test", "*")

				dd, err := db.LoadChildren(name, nil, 10)
				assert.NoError(t, err)

				assert.Len(t, dd, 10)
				for i, d := range dd {
					assert.Equal(t, d, []byte(fmt.Sprintf("data %d", i)))
				}
			})

			t.Run("When loading with limit", func(t *testing.T) {
				name := resource.NewName("test", "*")

				dd, err := db.LoadChildren(name, nil, 5)
				assert.NoError(t, err)

				assert.Len(t, dd, 5)
				for i, d := range dd {
					assert.Equal(t, d, []byte(fmt.Sprintf("data %d", i)))
				}
			})

			t.Run("When loading elements from", func(t *testing.T) {
				name := resource.NewName("test", "*")
				from := resource.NewName("test", "3")

				dd, err := db.LoadChildren(name, from, 10)
				assert.NoError(t, err)

				assert.Len(t, dd, 7)

				i := 3
				for _, d := range dd {
					assert.Equal(t, d, []byte(fmt.Sprintf("data %d", i)))
					i++
				}
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

			t.Run("When it's deleted", func(t *testing.T) {
				assert.NoError(t, db.Delete(name))

				t.Run("Data should not ne found", func(t *testing.T) {
					_, err := db.Load(name)
					assert.Equal(t, storage.ErrNotFound, err)
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
