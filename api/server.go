package api // import "miniboard.app/api"

import (
	"context"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	usersservice "miniboard.app/api/users"
	"miniboard.app/proto/users/v1"
	"miniboard.app/storage"
)

// Server is the api server.
type Server struct {
	httpServer *http.Server
}

// NewServer creates new api server.
func NewServer(ctx context.Context, db storage.DB) *Server {
	usersService := usersservice.New(db)

	gwMux := runtime.NewServeMux()
	users.RegisterUsersServiceHandlerClient(ctx, gwMux, usersservice.NewProxyClient(usersService))

	return &Server{
		httpServer: &http.Server{
			Handler: gwMux,
		},
	}
}

// Serve starts the server.
func (s *Server) Serve(ctx context.Context, lis net.Listener) error {
	logrus.Infof("[http] starting server on %s", lis.Addr())

	idleConnsClosed := make(chan struct{})
	go func() {
		<-ctx.Done()
		logrus.Infof("[http] stopping server")
		if err := s.httpServer.Shutdown(context.Background()); err != nil {
			logrus.Errorf("[http] error stopping server: %s", err)
		}
		close(idleConnsClosed)
	}()

	if err := s.httpServer.Serve(lis); err != nil {
		return errors.Wrap(err, "failed to start http server")
	}

	<-idleConnsClosed

	return nil
}
