package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ngalaiko/miniboard/server"
	"github.com/ngalaiko/miniboard/server/api"
	"github.com/ngalaiko/miniboard/server/db"
	"github.com/ngalaiko/miniboard/server/email"
	"github.com/ngalaiko/miniboard/server/email/disabled"
	"github.com/ngalaiko/miniboard/server/email/smtp"
	"github.com/sirupsen/logrus"
)

func main() {
	domain := flag.String("domain", "http://localhost:8080", "Service domain.")
	addr := flag.String("addr", ":8080", "Address to listen for connections.")

	dbType := flag.String("db-type", "sqlite3", "Database type (sqlite3, postgres).")
	dbAddr := flag.String("db-addr", "db.sqlite", "Database URI to connect to.")

	sslCert := flag.String("ssl-cert", "", "Path to ssl certificate.")
	sslKey := flag.String("ssl-key", "", "Path to ssl key.")

	smtpHost := flag.String("smtp-host", "", "SMTP server host.")
	smtpPort := flag.String("smtp-port", "", "SMTP server port.")
	smtpSender := flag.String("smtp-sender", "", "SMTP sender.")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv, err := server.New(ctx, &server.Config{
		HTTP: &api.HTTPConfig{
			Addr: *addr,
			TLS: &api.TLSConfig{
				CertPath: *sslCert,
				KeyPath:  *sslKey,
			},
			Domain: *domain,
		},
		DB: &db.Config{
			Driver: *dbType,
			Addr:   *dbAddr,
		},
	}, emailClient(*smtpHost, *smtpPort, *smtpSender))
	if err != nil {
		logrus.Fatalf("failed to create server: %s", err)
	}

	errCh := make(chan error)
	go func() {
		shutdownCh := make(chan os.Signal)
		signal.Notify(shutdownCh, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)
		<-shutdownCh

		shutdownTimeout := 30 * time.Second
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		errCh <- srv.Shutdown(shutdownCtx)
	}()

	if err := srv.Serve(ctx); err != nil {
		logrus.Fatalf("failed to start the server: %s", err)
	}

	if err := <-errCh; err != nil {
		logrus.Fatalf("error during shutdown: %s", err)
	}
}

func emailClient(host, port, sender string) email.Client {
	if host == "" {
		return disabled.New()
	}
	return smtp.New(
		host,
		port,
		sender,
		os.Getenv("SMTP_USERNAME"),
		os.Getenv("SMTP_PASSWORD"),
	)
}
