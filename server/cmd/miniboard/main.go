package main

import (
	"context"
	"flag"
	"net"
	"os"

	"github.com/ngalaiko/miniboard/server"
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

	filePath := flag.String("static-path", "", "Filepath to static files.")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sqldb, err := db.New(ctx, *dbType, *dbAddr)
	if err != nil {
		logrus.Fatalf("%s", err)
	}

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		logrus.Fatalf("failed to open a connection: %s", err)
	}

	srv, err := server.New(ctx, sqldb, emailClient(*smtpHost, *smtpPort, *smtpSender), *filePath, *domain)
	if err != nil {
		logrus.Fatalf("failed to create server: %s", err)
	}

	if err := srv.Serve(ctx, lis, &server.TLSConfig{
		CertPath: *sslCert,
		KeyPath:  *sslKey,
	}); err != nil {
		logrus.Fatalf("failed to start the server: %s", err)
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
