package authorizations // import "miniboard.app/services/authorizations"

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"miniboard.app/jwt"
	"miniboard.app/passwords"
	"miniboard.app/proto/users/authorizations/v1"
	"miniboard.app/storage"
	"miniboard.app/storage/bolt"
)

func Test_AuthorizationsService(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := testDB(ctx, t)
	jwt := jwt.NewService(db)
	passwords := passwords.NewService(db)

	t.Run("With new service", func(t *testing.T) {
		service := New(jwt, passwords)
		t.Run("When creating authorization for non existing user", func(t *testing.T) {
			auth, err := service.CreateAuthorization(ctx, &authorizations.CreateAuthorizationRequest{
				Parent:   "users/name",
				Password: "a passsword",
			})
			t.Run("Then error should be NotFound", func(t *testing.T) {
				status, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Nil(t, auth)
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
