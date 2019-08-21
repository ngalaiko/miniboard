package jwt // import "miniboard.app/jwt"

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Service(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("With new service", func(t *testing.T) {
		service, err := New(ctx, testStorage(ctx, t))
		assert.NoError(t, err)

		t.Run("When creating a token", func(t *testing.T) {
			testSubject := "test subject"
			token, err := service.NewToken(testSubject, time.Hour)

			t.Run("It should return a token", func(t *testing.T) {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			})

			t.Run("When parsing the token", func(t *testing.T) {
				parsedSubject, err := service.Validate(token)

				t.Run("It should no error", func(t *testing.T) {
					assert.NoError(t, err)
					assert.Equal(t, parsedSubject, testSubject)
				})
			})
		})
	})
}
