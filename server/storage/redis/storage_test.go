package redis

import (
	"context"
	"errors"
	"fmt"
	"testing"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/ngalaiko/miniboard/server/storage"
	"github.com/ngalaiko/miniboard/server/storage/resource"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

func Test_DB(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s, err := miniredis.Run()
	assert.NoError(t, err)
	defer s.Close()

	host := s.Addr()

	t.Run("With redis", func(t *testing.T) {
		db, err := New(ctx, host)
		if err != nil {
			t.Fatalf("failed to create database: %s", err)
		}

		assert.NoError(t, db.db.FlushAll().Err())

		t.Run("When data doesn't exist", func(t *testing.T) {
			t.Run("it should not be found", func(t *testing.T) {
				loaded, err := db.Load(ctx, resource.NewName("undefined", ksuid.New().String()))
				assert.Empty(t, loaded)
				assert.True(t, errors.Is(err, storage.ErrNotFound))
			})
		})

		t.Run("When a few elements exist", func(t *testing.T) {
			for i := 0; i < 10; i++ {
				name := resource.NewName("test", fmt.Sprint(i))
				data := []byte(fmt.Sprintf("data %d", i))
				assert.NoError(t, db.Store(ctx, name, data))
			}

			t.Run("When loading all", func(t *testing.T) {
				name := resource.NewName("test", "*")

				t.Run("Should load all", func(t *testing.T) {
					dd, err := db.LoadAll(ctx, name)
					assert.NoError(t, err)
					assert.Equal(t, 10, len(dd))
				})
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

						assert.Equal(t, fmt.Sprintf("data %d", 11-c), string(r.Data))
						return true, nil
					})

					assert.NoError(t, err)
					assert.Equal(t, 11, c)
				})
			})

			t.Run("When removing one element from the middle", func(t *testing.T) {
				name := resource.NewName("test", "*")

				c := 0
				err = db.ForEach(ctx, name, nil, func(r *resource.Resource) (bool, error) {
					c++
					if c == 5 {
						assert.NoError(t, db.Delete(ctx, r.Name))
						return false, nil
					}
					return true, nil
				})

				t.Run("Should be deleted", func(t *testing.T) {
					c := 0
					err := db.ForEach(ctx, name, nil, func(r *resource.Resource) (bool, error) {
						c++
						return true, nil
					})
					assert.NoError(t, err)
					assert.Equal(t, 9, c)
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

				t.Run("Data should not be found", func(t *testing.T) {
					_, err := db.Load(ctx, name)
					assert.Equal(t, storage.ErrNotFound, err)
				})
			})

			t.Run("When it's updated", func(t *testing.T) {
				assert.NoError(t, db.Store(ctx, name, []byte("updated")))

				t.Run("Data should be updated", func(t *testing.T) {
					d, err := db.Load(ctx, name)
					assert.NoError(t, err)
					assert.Equal(t, []byte("updated"), d)
				})

				t.Run("Data should not be duplicated", func(t *testing.T) {
					seen := map[string]bool{}
					err := db.ForEach(ctx, resource.NewName("noroot", "*"), nil, func(r *resource.Resource) (bool, error) {
						assert.False(t, seen[r.Name.String()])
						seen[r.Name.String()] = true
						return true, nil
					})
					assert.NoError(t, err)
				})
			})
		})
	})
}
