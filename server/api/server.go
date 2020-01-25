package api

import (
	"context"
	"net"
	"net/http"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
	"miniboard.app/email"
	"miniboard.app/images"
	"miniboard.app/jwt"
	"miniboard.app/storage"
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
	log("server").Infof("using domain: %s", domain)

	jwtService := jwt.NewService(ctx, db)
	images := images.New(db)

	handler := httpHandler(web.Handler(), jwtService, images)
	grpcServer := grpcServer(db, emailClient, jwtService, images, domain)
	grpcWebServer := grpcweb.WrapServer(grpcServer, grpcweb.WithAllowedRequestHeaders([]string{"Set-Cookie"}))

	grpcWebHandler := http.Handler(grpcWebServer)
	grpcWebHandler = withCompression(grpcWebHandler)

	grpcWebProxyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !grpcWebServer.IsGrpcWebRequest(r) {
			handler.ServeHTTP(w, r)
			return
		}

		grpcWebHandler.ServeHTTP(w, r)
	})

	srv := &Server{
		httpServer: &http.Server{
			Handler: grpcWebProxyHandler,
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
