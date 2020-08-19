package users

import (
	"context"
	"testing"

	"github.com/ngalaiko/miniboard/server/api/actor"
	"github.com/ngalaiko/miniboard/server/storage/resource"
	"github.com/stretchr/testify/assert"
)

func Test_UsersService(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("With new service", func(t *testing.T) {
		service := NewService()
		t.Run("When getting user", func(t *testing.T) {
			ctx = actor.NewContext(ctx, resource.NewName("users", "name"))
			user, err := service.GetMe(ctx, &GetMeRequest{})
			t.Run("Should response with the user struct", func(t *testing.T) {
				assert.NoError(t, err)
				user.Name = "users/name"
			})
		})
	})
}