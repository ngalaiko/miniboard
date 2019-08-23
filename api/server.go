package api // import "miniboard.app/api"

import (
	"context"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	usersservice "miniboard.app/api/users"
	authenticatationsservice "miniboard.app/api/users/authorizations"
	"miniboard.app/jwt"
	"miniboard.app/passwords"
	"miniboard.app/proto/users/authorizations/v1"
	"miniboard.app/proto/users/v1"
	"miniboard.app/storage"
)

// Server is the api server.
type Server struct {
	httpServer *http.Server
}

// NewServer creates new api server.
func NewServer(ctx context.Context, db storage.Storage) *Server {
	passwordsService := passwords.NewService(db)
	usersService := usersservice.New(db, passwordsService)
	jwtService := jwt.NewService(db)
	authorizationsService := authenticatationsservice.New(jwtService, passwordsService)

	gwMux := runtime.NewServeMux()

	users.RegisterUsersServiceHandlerClient(
		ctx,
		gwMux,
		usersservice.NewProxyClient(usersService),
	)
	authorizations.RegisterAuthorizationsServiceHandlerClient(
		ctx,
		gwMux,
		authenticatationsservice.NewProxyClient(authorizationsService),
	)

	return &Server{
		httpServer: &http.Server{
			Handler: withAccessLogs(
				withAuthorization(gwMux, jwtService),
			),
		},
	}
}

// Serve starts the server.
func (s *Server) Serve(ctx context.Context, lis net.Listener) error {
	log("http").Infof("starting server on %s", lis.Addr())

	idleConnsClosed := make(chan struct{})
	go func() {
		<-ctx.Done()
		log("http").Infof("stopping server")
		if err := s.httpServer.Shutdown(context.Background()); err != nil {
			log("http").Errorf("error stopping server: %s", err)
		}
		close(idleConnsClosed)
	}()

	if err := s.httpServer.Serve(lis); err != nil {
		return errors.Wrap(err, "failed to start http server")
	}

	<-idleConnsClosed

	return nil
}

func log(src string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": src,
	})
}
