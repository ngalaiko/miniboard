package api // import "miniboard.app/api"

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	usersservice "miniboard.app/api/users"
	authenticatationsservice "miniboard.app/api/users/authentications"
	"miniboard.app/jwt"
	"miniboard.app/passwords"
	"miniboard.app/proto/users/authentications/v1"
	"miniboard.app/proto/users/v1"
	"miniboard.app/storage"
)

// Server is the api server.
type Server struct {
	httpServer *http.Server
}

// NewServer creates new api server.
func NewServer(ctx context.Context, db storage.DB) *Server {
	passwordsService := passwords.NewService(db)
	usersService := usersservice.New(db, passwordsService)
	jwtService := jwt.NewService(db)
	authenticationsService := authenticatationsservice.New(jwtService, passwordsService)

	gwMux := runtime.NewServeMux()

	users.RegisterUsersServiceHandlerClient(
		ctx,
		gwMux,
		usersservice.NewProxyClient(usersService),
	)
	authentications.RegisterAuthenticationsServiceHandlerClient(
		ctx,
		gwMux,
		authenticatationsservice.NewProxyClient(authenticationsService),
	)

	return &Server{
		httpServer: &http.Server{
			Handler: withAccessLogs(gwMux),
		},
	}
}

func withAccessLogs(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer log("access").WithFields(logrus.Fields{
			"method":   r.Method,
			"path":     r.URL.String(),
			"ts":       start.Format(time.RFC3339),
			"duration": time.Since(start),
		}).Info()
		h.ServeHTTP(w, r)
	})
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
