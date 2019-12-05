package users

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"miniboard.app/api/actor"
	"miniboard.app/proto/users/v1"
	"miniboard.app/storage"
	"miniboard.app/storage/bolt"
	"miniboard.app/storage/resource"
)

func Test_UsersService(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := testDB(ctx, t)

	t.Run("With new service", func(t *testing.T) {
		service := New(db)
		t.Run("When getting user", func(t *testing.T) {
			ctx = actor.NewContext(ctx, resource.NewName("users", "name"))
			user, err := service.GetMe(ctx, &users.GetMeRequest{})
			t.Run("Should response with the user struct", func(t *testing.T) {
				assert.NoError(t, err)
				user.Name = "users/name"
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
