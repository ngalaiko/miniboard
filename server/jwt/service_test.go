package jwt

import (
	"context"
	"testing"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/ngalaiko/miniboard/server/storage/redis"
	"github.com/ngalaiko/miniboard/server/storage/resource"
	"github.com/stretchr/testify/assert"
)

func Test_Service(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("With new service", func(t *testing.T) {
		s, err := miniredis.Run()
		assert.NoError(t, err)
		defer s.Close()

		db, err := redis.New(ctx, s.Addr())
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
