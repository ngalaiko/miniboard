package api

import (
	"context"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
	authenticatationsservice "miniboard.app/api/authorizations"
	usersservice "miniboard.app/api/users"
	articlesservice "miniboard.app/api/users/articles"
	"miniboard.app/jwt"
	"miniboard.app/passwords"
	"miniboard.app/proto/authorizations/v1"
	"miniboard.app/proto/users/articles/v1"
	"miniboard.app/proto/users/v1"
	"miniboard.app/storage"
	"miniboard.app/web"
)

// Server is the api server.
type Server struct {
	httpServer *http.Server
}

// NewServer creates new api server.
func NewServer(ctx context.Context, db storage.Storage) *Server {
	passwordsService := passwords.NewService(db)
	usersService := usersservice.New(db, passwordsService)
	jwtService := jwt.NewService(ctx, db)
	authorizationsService := authenticatationsservice.New(jwtService, passwordsService)
	articlesService := articlesservice.New(db)

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
	articles.RegisterArticlesServiceHandlerClient(
		ctx,
		gwMux,
		articlesservice.NewProxyClient(articlesService),
	)

	mux := http.NewServeMux()
	mux.Handle("/api/", withAuthorization(
		convertFormData(gwMux),
		jwtService),
	)
	mux.Handle("/", web.Handler())

	srv := &Server{
		httpServer: &http.Server{
			Handler: withGzip(withAccessLogs(mux)),
		},
	}
	if err := http2.ConfigureServer(srv.httpServer, nil); err != nil {
		log("http2").Errorf("can't configure http2: %s", err)
	}
	return srv
}

// Serve starts the server.
func (s *Server) Serve(ctx context.Context, lis net.Listener, tlsConfig *TLSConfig) error {
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

	switch tlsConfig != nil && tlsConfig.valid() {
	case true:
		log("http").Infof("tls cert: %s", tlsConfig.CertPath)
		log("http").Infof("tls key: %s", tlsConfig.KeyPath)
		if err := s.httpServer.ServeTLS(lis, tlsConfig.CertPath, tlsConfig.KeyPath); err != nil {
			return errors.Wrap(err, "failed to start tls http server")
		}
	case false:
		if err := s.httpServer.Serve(lis); err != nil {
			return errors.Wrap(err, "failed to start http server")
		}
	}

	<-idleConnsClosed

	return nil
}

func log(src string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": src,
	})
}
