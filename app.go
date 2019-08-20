package application // import "miniflux.app/application"

import (
	"context"

	"miniboard.app/application/api"
)

// Application is the main application object.
type Application struct {
	Server *api.Server
}

// New creates new application.
func New(ctx context.Context) *Application {
	return &Application{
		Server: api.NewServer(ctx),
	}
}
