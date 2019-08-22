package users // import "miniboard.app/api/users"

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

		t.Run("When creating a new user empty name", func(t *testing.T) {
			user, err := service.CreateUser(ctx, &users.CreateUserRequest{
				Password: "test password",
			})
			t.Run("Then an error should be returned", func(t *testing.T) {
				assert.Nil(t, user)
				assert.Error(t, err)

				status, _ := status.FromError(err)
				assert.Equal(t, status.Code(), codes.InvalidArgument)
			})
		})

		t.Run("When creating a new user empty password", func(t *testing.T) {
			user, err := service.CreateUser(ctx, &users.CreateUserRequest{
				Name: "test name",
			})
			t.Run("Then an error should be returned", func(t *testing.T) {
				assert.Nil(t, user)
				assert.Error(t, err)

				status, _ := status.FromError(err)
				assert.Equal(t, status.Code(), codes.InvalidArgument)
			})
		})

		t.Run("When getting non existing user", func(t *testing.T) {
			user, err := service.GetUser(ctx, &users.GetUserRequest{
				Name: "test name",
			})
			t.Run("Then an error should be returned", func(t *testing.T) {
				assert.Nil(t, user)
				assert.Error(t, err)

				status, _ := status.FromError(err)
				assert.Equal(t, status.Code(), codes.NotFound)
			})
		})

		t.Run("When creating a new user with password", func(t *testing.T) {
			testName := "test name"
			createdUser, err := service.CreateUser(ctx, &users.CreateUserRequest{
				Name:     testName,
				Password: "test password",
			})
			t.Run("Then a user should be created", func(t *testing.T) {
				assert.NoError(t, err)
				assert.NotNil(t, createdUser)
			})

			t.Run("When getting the user", func(t *testing.T) {
				user, err := service.GetUser(ctx, &users.GetUserRequest{
					Name: testName,
				})
				t.Run("Then the user should be returned", func(t *testing.T) {
					assert.NoError(t, err)
					assert.Equal(t, user.Name, createdUser.Name)
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
