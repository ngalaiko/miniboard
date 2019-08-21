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
			token, err := service.NewToken("test subject", time.Hour)

			t.Run("It should return a token", func(t *testing.T) {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			})
		})
	})
}
