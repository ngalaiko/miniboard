package api

import (
	"context"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
	articlesservice "miniboard.app/articles"
	codesservice "miniboard.app/codes"
	"miniboard.app/email"
	"miniboard.app/jwt"
	"miniboard.app/proto/codes/v1"
	"miniboard.app/proto/users/articles/v1"
	"miniboard.app/proto/users/v1"
	"miniboard.app/storage"
	usersservice "miniboard.app/users"
	"miniboard.app/web"
)

// Server is the api server.
type Server struct {
	httpServer *http.Server
}

// NewServer creates new api server.
func NewServer(
	ctx context.Context,
	db storage.Storage,
	emailClient email.Client,
	domain string,
) *Server {
	jwtService := jwt.NewService(ctx, db)
	usersService := usersservice.New(db)
	codesService := codesservice.New(domain, emailClient, jwtService)
	articlesService := articlesservice.New(db)

	gwMux := runtime.NewServeMux()

	users.RegisterUsersServiceHandlerClient(
		ctx,
		gwMux,
		usersservice.NewProxyClient(usersService),
	)
	codes.RegisterCodesServiceHandlerClient(
		ctx,
		gwMux,
		codesservice.NewProxyClient(codesService),
	)
	articles.RegisterArticlesServiceHandlerClient(
		ctx,
		gwMux,
		articlesservice.NewProxyClient(articlesService),
	)

	mux := http.NewServeMux()
	mux.Handle("/api/", gwMux)
	mux.Handle("/logout", removeCookie())
	mux.Handle("/", web.Handler())

	handler := http.Handler(mux)
	handler = authorize(handler, jwtService)
	handler = exchangeAuthCode(handler, jwtService)
	handler = withGzip(handler)
	handler = withAccessLogs(handler)

	srv := &Server{
		httpServer: &http.Server{
			Handler: handler,
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
