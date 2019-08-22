package passwords // import "miniboard.app/passwords"

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

func Test_passwords(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := NewService(testDB(ctx, t))

	t.Run("With an empty service", func(t *testing.T) {
		t.Run("When validating non existing password", func(t *testing.T) {
			valid, err := service.Validate("user", "password")
			t.Run("Then it should not be valid", func(t *testing.T) {
				assert.Equal(t, errors.Cause(err), storage.ErrNotFound)
				assert.False(t, valid)
			})
		})
		t.Run("When adding a password for a user", func(t *testing.T) {
			username := "user"
			password := "user's password"
			err := service.Set(username, password)
			assert.NoError(t, err)

			t.Run("When validating correct password", func(t *testing.T) {
				valid, err := service.Validate(username, password)
				t.Run("Then it should be valid", func(t *testing.T) {
					assert.NoError(t, err)
					assert.True(t, valid)
				})
			})

			t.Run("When validating incorrect password", func(t *testing.T) {
				valid, err := service.Validate(username, "something else")
				t.Run("Then it should not be valid", func(t *testing.T) {
					assert.NoError(t, err)
					assert.False(t, valid)
				})
			})
		})
	})
}

func testDB(ctx context.Context, t *testing.T) storage.DB {
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
