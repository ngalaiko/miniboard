package jwt

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"miniboard.app/storage/redis"
	"miniboard.app/storage/resource"
)

func Test_Service(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("With new service", func(t *testing.T) {
		host := os.Getenv("REDIS_HOST")
		if host == "" {
			t.Skip("REDIS_HOST is not set")
		}

		db, err := redis.New(ctx, host)
		assert.NoError(t, err)

		service := NewService(ctx, db)

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
	})
}
