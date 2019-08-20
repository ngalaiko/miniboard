package bolt // import "miniboard.app/storage/bolt_test"

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"miniboard.app/application/storage"
)

func Test_DB_Store(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := testBucket(ctx, t)

	err := db.Store([]byte("id"), []byte("data"))
	assert.NoError(t, err)
}

func Test_DB_Load(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("given a bucket", func(t *testing.T) {
		db := testBucket(ctx, t)

		t.Run("when data doesn't exist", func(t *testing.T) {
			t.Run("it should not be found", func(t *testing.T) {
				loaded, err := db.Load([]byte("any id"))
				assert.Empty(t, loaded)
				assert.Equal(t, errors.Cause(err), storage.ErrNotFound)
			})
		})

		t.Run("when data exists", func(t *testing.T) {
			id := []byte("id")
			data := []byte("data")
			assert.NoError(t, db.Store(id, data))

			t.Run("it should be found", func(t *testing.T) {
				loaded, err := db.Load(id)
				assert.Equal(t, loaded, data)
				assert.NoError(t, err)
			})
		})
	})
}

func Test_DB_LoadPrefix(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("given a bucket", func(t *testing.T) {

		db := testBucket(ctx, t)

		t.Run("when data doesn't exist", func(t *testing.T) {
			t.Run("it should not be found", func(t *testing.T) {
				loaded, err := db.LoadPrefix([]byte("any id"))
				assert.Empty(t, loaded)
				assert.NoError(t, err)
			})
		})

		t.Run("when data exists", func(t *testing.T) {
			prefix := []byte("prefix")

			id1 := append(prefix, []byte("id1")...)
			data1 := []byte("data 1")
			assert.NoError(t, db.Store(id1, data1))

			id2 := append(prefix, []byte("id2")...)
			data2 := []byte("data 2")
			assert.NoError(t, db.Store(id2, data2))

			t.Run("it should be found", func(t *testing.T) {
				loaded, err := db.LoadPrefix(prefix)
				assert.Len(t, loaded, 2)
				assert.Contains(t, loaded, data1)
				assert.Contains(t, loaded, data2)
				assert.NoError(t, err)
			})
		})
	})
}

func testBucket(ctx context.Context, t *testing.T) *Bucket {
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

	bucket, err := db.Bucket("test")
	if err != nil {
		t.Fatalf("failed to create database: %s", err)
	}
	return bucket
}
