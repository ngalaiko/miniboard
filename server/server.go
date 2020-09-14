package server

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/ngalaiko/miniboard/server/articles"
	"github.com/ngalaiko/miniboard/server/codes"
	"github.com/ngalaiko/miniboard/server/email"
	"github.com/ngalaiko/miniboard/server/feeds"
	"github.com/ngalaiko/miniboard/server/fetch"
	"github.com/ngalaiko/miniboard/server/jwt"
	"github.com/ngalaiko/miniboard/server/middleware"
	"github.com/ngalaiko/miniboard/server/sources"
	"github.com/ngalaiko/miniboard/server/tokens"
	"github.com/ngalaiko/miniboard/server/web"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
)

// todo: make it shorter
const authDuration = 28 * 24 * time.Hour

// TLSConfig contains ssl certificates.
type TLSConfig struct {
	CertPath string
	KeyPath  string
}

func (cfg *TLSConfig) valid() bool {
	return cfg.CertPath != "" || cfg.KeyPath != ""
}

// Server is the api server.
type Server struct {
	httpServer *http.Server
}

// New creates new api server.
func New(
	ctx context.Context,
	sqldb *sql.DB,
	emailClient email.Client,
	filePath string,
	domain string,
) (*Server, error) {
	log("server").Infof("using domain: %s", domain)

	fetcher := fetch.New()

	jwtService := jwt.NewService(ctx, sqldb)
	articlesService := articles.NewService(sqldb, fetcher)
	feedsService := feeds.NewService(ctx, sqldb, fetcher, articlesService)
	codesService := codes.NewService(domain, emailClient, jwtService)
	tokensService := tokens.NewService(jwtService)
	sourcesService := sources.NewService(articlesService, feedsService, fetcher)

	gwMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			OrigName:     false,
			EmitDefaults: true,
		}),
		runtime.WithForwardResponseOption(func(ctx context.Context, rw http.ResponseWriter, msg proto.Message) error {
			if token, ok := msg.(*tokens.Token); ok {
				http.SetCookie(rw, &http.Cookie{
					Name:     middleware.AuthCookie,
					Value:    token.Token,
					Path:     "/",
					Expires:  time.Now().Add(authDuration),
					HttpOnly: true,
				})
			}
			return nil
		}),
	)

	if err := articles.RegisterArticlesServiceHandlerServer(ctx, gwMux, articlesService); err != nil {
		return nil, fmt.Errorf("failed to register articles http handler: %w", err)
	}

	if err := tokens.RegisterTokensServiceHandlerServer(ctx, gwMux, tokensService); err != nil {
		return nil, fmt.Errorf("failed to register tokens http handler: %w", err)
	}

	if err := codes.RegisterCodesServiceHandlerServer(ctx, gwMux, codesService); err != nil {
		return nil, fmt.Errorf("failed to register codes http handler: %w", err)
	}

	if err := sources.RegisterSourcesServiceHandlerServer(ctx, gwMux, sourcesService); err != nil {
		return nil, fmt.Errorf("failed to register sources http handler: %w", err)
	}

	if err := feeds.RegisterFeedsServiceHandlerServer(ctx, gwMux, feedsService); err != nil {
		return nil, fmt.Errorf("failed to register feeds http handler: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/logout", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     middleware.AuthCookie,
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})
	}))

	mux.Handle("/api/v1/tokens", gwMux)
	mux.Handle("/api/v1/codes", gwMux)
	mux.Handle("/api/", middleware.Authorized(gwMux, jwtService))
	mux.Handle("/", web.Handler(filePath))

	handler := http.Handler(mux)
	handler = middleware.WithAccessLogs(handler)
	handler = middleware.WithCompression(handler)
	httpServer := &http.Server{
		Handler: handler,
	}
	if err := http2.ConfigureServer(httpServer, nil); err != nil {
		return nil, fmt.Errorf("can't configure http: %w", err)
	}

	return &Server{
		httpServer: httpServer,
	}, nil
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
			return fmt.Errorf("failed to start tls http server: %w", err)
		}
	case false:
		if err := s.httpServer.Serve(lis); err != nil {
			return fmt.Errorf("failed to start http server: %w", err)
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
