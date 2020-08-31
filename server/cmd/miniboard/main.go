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
	"github.com/ngalaiko/miniboard/server/storage/redis"
	"github.com/sirupsen/logrus"
)

func main() {
	domain := flag.String("domain", "http://localhost:8080", "Service domain.")
	addr := flag.String("addr", ":8080", "Address to listen for connections.")

	redisURI := flag.String("redis-uri", "", "Redis URI to connect to.")

	psqlURI := flag.String("psql-uri", "", "Postgres URI to connect to.")

	sslCert := flag.String("ssl-cert", "", "Path to ssl certificate.")
	sslKey := flag.String("ssl-key", "", "Path to ssl key.")

	smtpHost := flag.String("smtp-host", "", "SMTP server host.")
	smtpPort := flag.String("smtp-port", "", "SMTP server port.")
	smtpSender := flag.String("smtp-sender", "", "SMTP sender.")

	filePath := flag.String("static-path", "", "Filepath to static files.")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	postgres, err := db.NewPostgres(*psqlURI)
	if err != nil {
		logrus.Fatalf("failed to connect to postgres: %s", err)
	}

	if err := db.Migrate(ctx, postgres); err != nil {
		logrus.Fatalf("failed to migrate postgres: %s", err)
	}

	redis, err := redis.New(ctx, *redisURI)
	if err != nil {
		logrus.Fatalf("failed to connect to database: %s", err)
	}

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		logrus.Fatalf("failed to open a connection: %s", err)
	}

	srv, err := server.New(ctx, redis, postgres, emailClient(*smtpHost, *smtpPort, *smtpSender), *filePath, *domain)
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
