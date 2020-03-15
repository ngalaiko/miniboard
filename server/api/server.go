package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
	articlesservice "miniboard.app/articles"
	codes "miniboard.app/codes"
	"miniboard.app/email"
	feedsservice "miniboard.app/feeds"
	"miniboard.app/fetch"
	"miniboard.app/images"
	"miniboard.app/jwt"
	tokens "miniboard.app/proto/tokens/v1"
	articles "miniboard.app/proto/users/articles/v1"
	sources "miniboard.app/proto/users/sources/v1"
	users "miniboard.app/proto/users/v1"
	sourcesservice "miniboard.app/sources"
	"miniboard.app/storage"
	tokensservice "miniboard.app/tokens"
	usersservice "miniboard.app/users"
	"miniboard.app/web"
)

// todo: make it shorter
const authDuration = 28 * 24 * time.Hour

// Server is the api server.
type Server struct {
	httpServer *http.Server
}

// NewServer creates new api server.
func NewServer(
	ctx context.Context,
	db storage.Storage,
	emailClient email.Client,
	filePath string,
	domain string,
) (*Server, error) {
	log("server").Infof("using domain: %s", domain)

	fetcher := fetch.New()

	imagesService := images.New(db)
	jwtService := jwt.NewService(ctx, db)
	articlesService := articlesservice.New(db, imagesService, fetcher)
	feedsService := feedsservice.New(ctx, db, articlesService)
	usersService := usersservice.New()
	codesService := codes.New(domain, emailClient, jwtService)
	tokensService := tokensservice.New(jwtService)
	sourcesService := sourcesservice.New(articlesService, feedsService, fetcher)

	gwMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			OrigName:     false,
			EmitDefaults: true,
		}),
		runtime.WithForwardResponseOption(func(ctx context.Context, rw http.ResponseWriter, msg proto.Message) error {
			if token, ok := msg.(*tokens.Token); ok {
				http.SetCookie(rw, &http.Cookie{
					Name:     authCookie,
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
		return nil, fmt.Errorf("failed to register tokens http handler: %w", err)
	}

	if err := users.RegisterUsersServiceHandlerServer(ctx, gwMux, usersService); err != nil {
		return nil, fmt.Errorf("failed to register tokens http handler: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/logout", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     authCookie,
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})
	}))

	mux.Handle("/api/", authorize(gwMux, jwtService))

	imagesHandler := imagesService.Handler()
	webHandler := web.Handler(filePath)

	imageRegExp := regexp.MustCompile("images/.+")

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if imageRegExp.MatchString(r.RequestURI) {
			imagesHandler.ServeHTTP(w, r)
			return
		}

		webHandler.ServeHTTP(w, r)
	}))

	handler := http.Handler(mux)
	handler = withAccessLogs(handler)
	handler = withCompression(handler)
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
