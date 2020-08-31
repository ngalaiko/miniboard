package jwt

import (
	"context"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/ngalaiko/miniboard/server/db"
	"github.com/ngalaiko/miniboard/server/storage/resource"
	"github.com/stretchr/testify/assert"
)

func Test_Service(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tmpFile, err := ioutil.TempFile(os.TempDir(), "testdb-")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	sqlite, err := db.NewSQLite(tmpFile.Name())
	assert.NoError(t, err)
	assert.NoError(t, db.Migrate(ctx, sqlite))

	service := NewService(ctx, sqlite)

	t.Run("When creating a token", func(t *testing.T) {
		testSubject := resource.NewName("test", "test subject")
		token, err := service.NewToken(testSubject, time.Hour, "token")

		t.Run("It should return a token", func(t *testing.T) {
			assert.NoError(t, err)
			assert.NotEmpty(t, token)
		})

		t.Run("When parsing the token", func(t *testing.T) {
			parsedSubject, err := service.Validate(ctx, token, "token")

			t.Run("It should no error", func(t *testing.T) {
				assert.NoError(t, err)
				assert.Equal(t, testSubject, parsedSubject)
			})
		})

		t.Run("When the token is rotated", func(t *testing.T) {
			assert.NoError(t, service.newSigner(ctx))
		})
	})
}
