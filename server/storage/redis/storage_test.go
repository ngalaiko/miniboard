package redis

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"miniboard.app/storage"
	"miniboard.app/storage/resource"
)

func Test_DB(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("Given a bucket", func(t *testing.T) {
		host := os.Getenv("REDIS_HOST")
		if host == "" {
			t.Skip("no redis host provided")
		}

		db := testBucket(ctx, t, host)

		conn := db.(*Storage).db.Get()
		_, err := conn.Do("FLUSHALL")
		_ = conn.Close()

		assert.NoError(t, err)

		t.Run("When data doesn't exist", func(t *testing.T) {
			t.Run("it should not be found", func(t *testing.T) {
				loaded, err := db.Load(ctx, resource.NewName("undefined", ksuid.New().String()))
				assert.Empty(t, loaded)
				assert.Equal(t, errors.Cause(err), storage.ErrNotFound)
			})
		})

		t.Run("When a few elements exist", func(t *testing.T) {
			for i := 0; i < 10; i++ {
				name := resource.NewName("test", fmt.Sprint(i))
				data := []byte(fmt.Sprintf("data %d", i))
				assert.NoError(t, db.Store(ctx, name, data))
			}

			t.Run("When loading all elements", func(t *testing.T) {
				name := resource.NewName("test", "*")

				dd, err := db.LoadChildren(ctx, name, nil, 10)
				assert.NoError(t, err)

				assert.Len(t, dd, 10)
				for i, d := range dd {
					assert.Equal(t, d.Data, []byte(fmt.Sprintf("data %d", 9-i)))
				}
			})

			t.Run("When loading with limit", func(t *testing.T) {
				name := resource.NewName("test", "*")

				dd, err := db.LoadChildren(ctx, name, nil, 5)
				assert.NoError(t, err)

				assert.Len(t, dd, 5)
				for i, d := range dd {
					assert.Equal(t, d.Data, []byte(fmt.Sprintf("data %d", 9-i)))
				}
			})

			t.Run("When loading elements from", func(t *testing.T) {
				name := resource.NewName("test", "*")
				from := resource.NewName("test", "6")

				dd, err := db.LoadChildren(ctx, name, from, 10)
				assert.NoError(t, err)

				assert.Len(t, dd, 7)

				i := 6
				for _, d := range dd {
					assert.Equal(t, d.Data, []byte(fmt.Sprintf("data %d", i)))
					i--
				}
			})

			t.Run("When iterating through all elemetns", func(t *testing.T) {
				name := resource.NewName("test", "*")

				t.Run("Should iterate from end to start", func(t *testing.T) {
					c := 0
					err := db.ForEach(ctx, name, nil, func(r *resource.Resource) (bool, error) {
						c++
						assert.Equal(t, fmt.Sprintf("data %d", 10-c), string(r.Data))
						return true, nil
					})
					assert.NoError(t, err)
					assert.Equal(t, 10, c)
				})
			})

			t.Run("When iterating in batches", func(t *testing.T) {
				name := resource.NewName("test", "*")

				t.Run("Should iterate from end to start", func(t *testing.T) {
					c := 0
					var from *resource.Name
					err := db.ForEach(ctx, name, nil, func(r *resource.Resource) (bool, error) {
						c++

						assert.Equal(t, fmt.Sprintf("data %d", 10-c), string(r.Data))
						if c == 5 {
							from = r.Name
							return false, nil
						}
						return true, nil
					})
					assert.NoError(t, err)
					assert.Equal(t, 5, c)

					err = db.ForEach(ctx, name, from, func(r *resource.Resource) (bool, error) {
						c++

						assert.Equal(t, fmt.Sprintf("data %d", 10-c), string(r.Data))
						return true, nil
					})

					assert.NoError(t, err)
					assert.Equal(t, 10, c)
				})
			})
		})

		t.Run("When root exists", func(t *testing.T) {
			name := resource.NewName("noroot", ksuid.New().String())
			data := []byte("data")
			assert.NoError(t, db.Store(ctx, name, data))

			t.Run("It should be found", func(t *testing.T) {
				loaded, err := db.Load(ctx, name)
				if assert.NoError(t, err) {
					assert.Equal(t, loaded, data)
				}
			})

			t.Run("When child exists", func(t *testing.T) {
				name := name.Child("child", "id")
				data := []byte("data")
				assert.NoError(t, db.Store(ctx, name, data))

				t.Run("It should be found", func(t *testing.T) {
					loaded, err := db.Load(ctx, name)
					if assert.NoError(t, err) {
						assert.Equal(t, loaded, data)
					}
				})
			})

			t.Run("When it's deleted", func(t *testing.T) {
				assert.NoError(t, db.Delete(ctx, name))

				t.Run("Data should not ne found", func(t *testing.T) {
					_, err := db.Load(ctx, name)
					assert.Equal(t, storage.ErrNotFound, err)
				})
			})

			t.Run("When it's updated", func(t *testing.T) {
				assert.NoError(t, db.Update(ctx, name, []byte("updated")))

				t.Run("Data should not ne found", func(t *testing.T) {
					d, err := db.Load(ctx, name)
					assert.NoError(t, err)
					assert.Equal(t, []byte("updated"), d)
				})
			})
		})
	})
}

func testBucket(ctx context.Context, t *testing.T, host string) storage.Storage {
	s, err := New(ctx, "localhost:6379")
	if err != nil {
		t.Fatalf("failed to create database: %s", err)
	}
	return s
}
