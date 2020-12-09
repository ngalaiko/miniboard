package http

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"
)

type tlsConfig struct {
	KeyPath  string
	CertPath string
}

// Config contains http server configuration values.
type Config struct {
	Addr string     `yaml:"addr"`
	TLS  *tlsConfig `yaml:"tls"`
}

type logger interface {
	Info(string, ...interface{})
}

// Server represents an HTTP server.
type Server struct {
	logger logger
	server *http.Server
	cfg    *Config
}

// NewServer creates a new HTTP server.
//
// The handler can be nil, in which case http.DefaultServeMux is used.
func NewServer(cfg *Config, logger logger, handler http.Handler) (*Server, error) {
	if cfg == nil {
		cfg = &Config{
			// Preserve the stdlib default.
			Addr: ":http",
		}
	}
	srv := &Server{
		logger: logger,
		cfg:    cfg,
		server: &http.Server{
			Addr:    cfg.Addr,
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
	}

	if cfg.TLS == nil {
		return srv, nil
	}

	cert, err := tls.LoadX509KeyPair(cfg.TLS.CertPath, cfg.TLS.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read x509 keypair: %w", err)
	}

	srv.server.TLSConfig.Certificates = []tls.Certificate{cert}

	return srv, nil
}

// isTLS returns whether TLS is enabled.
func (srv *Server) isTLS() bool {
	return len(srv.server.TLSConfig.Certificates) > 0 || srv.server.TLSConfig.GetCertificate != nil
}

// Start starts the server, handling incoming requests.
//
// Accepted connections are configured to enable TCP keep-alives.
func (srv *Server) Start() error {
	ln, err := net.Listen("tcp", srv.cfg.Addr)
	if err != nil {
		return err
	}

	if srv.isTLS() {
		ln = tls.NewListener(ln, srv.server.TLSConfig)
		srv.logger.Info("listening https on %s", ln.Addr().String())
	} else {
		srv.logger.Info("listening http on %s", ln.Addr().String())
	}

	if err := srv.server.Serve(ln); err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown gracefully shutdowns the server.
func (srv *Server) Shutdown(ctx context.Context) error {
	if err := srv.server.Shutdown(ctx); err == context.DeadlineExceeded {
		return fmt.Errorf("timeout exceeded while waiting on shutdown")
	}
	return nil
}
