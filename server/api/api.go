package api

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/mmcdole/gofeed"
	"github.com/sirupsen/logrus"

	articlesv1 "github.com/ngalaiko/miniboard/server/genproto/articles/v1"
	codesv1 "github.com/ngalaiko/miniboard/server/genproto/codes/v1"
	feedsv1 "github.com/ngalaiko/miniboard/server/genproto/feeds/v1"
	longrunning "github.com/ngalaiko/miniboard/server/genproto/google/longrunning"
	sourcesv1 "github.com/ngalaiko/miniboard/server/genproto/sources/v1"
	tokensv1 "github.com/ngalaiko/miniboard/server/genproto/tokens/v1"

	"github.com/ngalaiko/miniboard/server/api/articles"
	"github.com/ngalaiko/miniboard/server/api/codes"
	"github.com/ngalaiko/miniboard/server/api/feeds"
	"github.com/ngalaiko/miniboard/server/api/operations"
	"github.com/ngalaiko/miniboard/server/api/sources"
	"github.com/ngalaiko/miniboard/server/api/tokens"
	"github.com/ngalaiko/miniboard/server/jwt"
	"github.com/ngalaiko/miniboard/server/web"
)

// todo: make it shorter
const authDuration = 28 * 24 * time.Hour

type fetcher interface {
	Fetch(context.Context, string) (*http.Response, error)
}

type emailSender interface {
	Send(to string, subject string, payload string) error
}

// TLSConfig contains ssl certificates.
type TLSConfig struct {
	CertPath string
	KeyPath  string
}

func (cfg *TLSConfig) valid() bool {
	return cfg != nil && (cfg.CertPath != "" || cfg.KeyPath != "")
}

// HTTPConfig is the http api configuration.
type HTTPConfig struct {
	Domain string
	Addr   string
	TLS    *TLSConfig
}

// HTTP exposes http api.
type HTTP struct {
	cfg    *HTTPConfig
	server *http.Server
}

// NewHTTP returns new http api.
func NewHTTP(ctx context.Context, cfg *HTTPConfig, sqldb *sql.DB, fetcher fetcher, emailClient emailSender, jwtService *jwt.Service) (*HTTP, error) {
	if cfg == nil {
		cfg = &HTTPConfig{}
	}
	articlesService := articles.NewService(sqldb)
	feedsService := feeds.NewService(ctx, sqldb, fetcher, articlesService, gofeed.NewParser())
	codesService := codes.NewService(cfg.Domain, emailClient, jwtService)
	tokensService := tokens.NewService(jwtService)
	operationsService := operations.New(ctx, sqldb)
	sourcesService := sources.NewService(articlesService, feedsService, operationsService, fetcher)

	gwMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			OrigName:     false,
			EmitDefaults: true,
		}),
		runtime.WithForwardResponseOption(func(ctx context.Context, rw http.ResponseWriter, msg proto.Message) error {
			if token, ok := msg.(*tokensv1.Token); ok {
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

	if err := articlesv1.RegisterArticlesServiceHandlerServer(ctx, gwMux, articlesService); err != nil {
		return nil, fmt.Errorf("failed to register articles http handler: %w", err)
	}

	if err := tokensv1.RegisterTokensServiceHandlerServer(ctx, gwMux, tokensService); err != nil {
		return nil, fmt.Errorf("failed to register tokens http handler: %w", err)
	}

	if err := codesv1.RegisterCodesServiceHandlerServer(ctx, gwMux, codesService); err != nil {
		return nil, fmt.Errorf("failed to register codes http handler: %w", err)
	}

	if err := sourcesv1.RegisterSourcesServiceHandlerServer(ctx, gwMux, sourcesService); err != nil {
		return nil, fmt.Errorf("failed to register sources http handler: %w", err)
	}

	if err := feedsv1.RegisterFeedsServiceHandlerServer(ctx, gwMux, feedsService); err != nil {
		return nil, fmt.Errorf("failed to register feeds http handler: %w", err)
	}

	if err := longrunning.RegisterOperationsHandlerServer(ctx, gwMux, operationsService); err != nil {
		return nil, fmt.Errorf("failed to register operations http handler: %w", err)
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

	mux.Handle("/api/v1/tokens", gwMux)
	mux.Handle("/api/v1/codes", gwMux)
	mux.Handle("/api/", authorized(gwMux, jwtService))
	mux.Handle("/", web.Handler())

	handler := http.Handler(mux)
	handler = withAccessLogs(handler)
	handler = withCompression(handler)
	handler = withRecover(handler)

	addr := "0.0.0.0:8080"
	if cfg.Addr != "" {
		addr = cfg.Addr
	}
	return &HTTP{
		cfg: cfg,
		server: &http.Server{
			Addr:    addr,
			Handler: handler,
			// https://blog.cloudflare.com/exposing-go-on-the-internet/
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
			TLSConfig: &tls.Config{
				NextProtos:       []string{"h2", "http/1.1"},
				MinVersion:       tls.VersionTLS12,
				CurvePreferences: []tls.CurveID{tls.CurveP256, tls.X25519},
				CipherSuites: []uint16{
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
					tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				},
				PreferServerCipherSuites: true,
			},
		},
	}, nil
}

// ListenAndServe stars http server.
func (h *HTTP) ListenAndServe(ctx context.Context) error {
	log("http").Infof("starting server on %s", h.server.Addr)

	if h.cfg.TLS.valid() {
		if err := h.server.ListenAndServeTLS(h.cfg.TLS.CertPath, h.cfg.TLS.KeyPath); err != http.ErrServerClosed {
			return err
		}
	} else {
		if err := h.server.ListenAndServe(); err != http.ErrServerClosed {
			return err
		}
	}
	return nil
}

// Shutdown gracefully stops the api.
func (h *HTTP) Shutdown(ctx context.Context) error {
	log("http").Infof("stopping server")

	if err := h.server.Shutdown(ctx); err == context.DeadlineExceeded {
		return fmt.Errorf("timeout exceeded while waiting on HTTP shutdown")
	}

	return nil
}

func log(src string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": src,
	})
}
