package api // import "miniboard.app/api"

import "context" // Server is the api server.

// Server is the api server.
type Server struct{}

// NewServer creates new api server.
func NewServer(ctx context.Context) *Server {
	return &Server{}
}
