package actor

import (
	"context"
)

type actorCtxKeyType struct{}

// Actor contains id of a user performing a request.
type Actor struct {
	ID string
}

// NewContext returns context with actor name set.
func NewContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, actorCtxKeyType{}, id)
}

// FromContext returns actor from context.
func FromContext(ctx context.Context) (*Actor, bool) {
	id, ok := ctx.Value(actorCtxKeyType{}).(string)
	return &Actor{
		ID: id,
	}, ok
}
