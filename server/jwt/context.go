package jwt

import "context"

type tokenContextKey struct{}

// NewContext creates a new context with the token value.
func NewContext(ctx context.Context, token *Token) context.Context {
	return context.WithValue(ctx, tokenContextKey{}, token)
}

// FromContext returns a token from the given context.
func FromContext(ctx context.Context) (*Token, bool) {
	token, ok := ctx.Value(tokenContextKey{}).(*Token)
	return token, ok
}
