package api

import (
	"context"
	"net"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"miniboard.app/email"
	"miniboard.app/images"
	"miniboard.app/jwt"
	"miniboard.app/storage"
	"miniboard.app/web"
)

// Server is the api server.
type Server struct {
	imagesService *images.Service
	jwtService    *jwt.Service
	db            storage.Storage
	domain        string
	filePath      string
	emailClient   email.Client
}

// NewServer creates new api server.
func NewServer(
	ctx context.Context,
	db storage.Storage,
	emailClient email.Client,
	filePath string,
	domain string,
) *Server {
	log("server").Infof("using domain: %s", domain)

	srv := &Server{
		db:            db,
		domain:        domain,
		emailClient:   emailClient,
		filePath:      filePath,
		jwtService:    jwt.NewService(ctx, db),
		imagesService: images.New(db),
	}
	return srv
}

// Serve starts the server.
func (s *Server) Serve(ctx context.Context, lis net.Listener, tlsConfig *TLSConfig) error {
	m := cmux.New(lis)

	grpcL := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpL := m.Match(cmux.Any())

	go s.serveGRPC(ctx, grpcL, tlsConfig)
	go s.serveHTTP(ctx, httpL, tlsConfig)

	<-ctx.Done()

	return nil
}

func (s *Server) serveGRPC(ctx context.Context, lis net.Listener, tlsConfig *TLSConfig) error {
	log("grpc").Infof("starting server on %s", lis.Addr())

	options := []grpc.ServerOption{}
	if tlsConfig != nil && tlsConfig.valid() {
		creds, err := credentials.NewServerTLSFromFile(tlsConfig.CertPath, tlsConfig.KeyPath)
		if err != nil {
			return errors.Wrap(err, "failed to read certificates")
		}

		log("grpc").Infof("tls cert: %s", tlsConfig.CertPath)
		log("grpc").Infof("tls key: %s", tlsConfig.KeyPath)
		options = append(options,
			grpc.Creds(creds),
		)
	}

	srv := grpcServer(s.db, s.emailClient, s.jwtService, s.imagesService, s.domain, options...)

	idleConnsClosed := make(chan struct{})
	go func() {
		<-ctx.Done()
		log("grpc").Info("stopping server")
		srv.GracefulStop()
		close(idleConnsClosed)
	}()

	if err := srv.Serve(lis); err != nil {
		return errors.Wrap(err, "error serving grpc")
	}
	return nil
}

func (s *Server) serveHTTP(ctx context.Context, lis net.Listener, tlsConfig *TLSConfig) error {
	log("http").Infof("starting server on %s", lis.Addr())

	handler := httpHandler(web.Handler(s.filePath), s.jwtService, s.imagesService)
	handler = http.Handler(handler)
	handler = withCompression(handler)
	httpServer := &http.Server{
		Handler: handler,
	}
	if err := http2.ConfigureServer(httpServer, nil); err != nil {
		return errors.Wrapf(err, "can't configure http")
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		<-ctx.Done()
		log("http").Infof("stopping server")
		if err := httpServer.Shutdown(context.Background()); err != nil {
			log("http").Errorf("error stopping server: %s", err)
		}
		close(idleConnsClosed)
	}()

	switch tlsConfig != nil && tlsConfig.valid() {
	case true:
		log("http").Infof("tls cert: %s", tlsConfig.CertPath)
		log("http").Infof("tls key: %s", tlsConfig.KeyPath)
		if err := httpServer.ServeTLS(lis, tlsConfig.CertPath, tlsConfig.KeyPath); err != nil {
			return errors.Wrap(err, "failed to start tls http server")
		}
	case false:
		if err := httpServer.Serve(lis); err != nil {
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
