package users

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"miniboard.app/proto/users/v1"
	"miniboard.app/storage"
	"miniboard.app/storage/bolt"
)

func Test_UsersService(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := testDB(ctx, t)

	t.Run("With new service", func(t *testing.T) {
		service := New(db)
		t.Run("When getting non existing user", func(t *testing.T) {
			user, err := service.GetUser(ctx, &users.GetUserRequest{
				Name: "users/name",
			})
			t.Run("Then an error should be returned", func(t *testing.T) {
				assert.Nil(t, user)
				assert.Error(t, err)

				status, _ := status.FromError(err)
				assert.Equal(t, status.Code(), codes.NotFound)
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
