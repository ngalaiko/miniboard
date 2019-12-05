package actor

import (
	"context"
	"strings"

	"miniboard.app/storage/resource"
)

type actorCtxKeyType struct{}

// NewContext returns context with actor name set.
func NewContext(ctx context.Context, actor *resource.Name) context.Context {
	return context.WithValue(ctx, actorCtxKeyType{}, actor)
}

// FromContext returns actor from context.
func FromContext(ctx context.Context) (*resource.Name, bool) {
	actor, ok := ctx.Value(actorCtxKeyType{}).(*resource.Name)
	return actor, ok
}

// Owns returns true if actor from context owns the resource.
func Owns(ctx context.Context, res *resource.Name) bool {
	actor, ok := ctx.Value(actorCtxKeyType{}).(*resource.Name)
	if !ok {
		return false
	}
	return strings.HasPrefix(res.String(), actor.String())
}
