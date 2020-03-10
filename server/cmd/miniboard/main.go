package main

import (
	"context"
	"flag"
	"net"
	"os"

	"github.com/sirupsen/logrus"
	"miniboard.app/api"
	"miniboard.app/email"
	"miniboard.app/email/disabled"
	"miniboard.app/email/smtp"
	"miniboard.app/storage/redis"
)

func main() {
	domain := flag.String("domain", "http://localhost:8080", "Service domain.")
	addr := flag.String("addr", ":8080", "Address to listen for connections.")

	redisURI := flag.String("redis-uri", "", "Redis URI to connect to.")

	sslCert := flag.String("ssl-cert", "", "Path to ssl certificate.")
	sslKey := flag.String("ssl-key", "", "Path to ssl key.")

	smtpHost := flag.String("smtp-host", "", "SMTP server host.")
	smtpPort := flag.String("smtp-port", "", "SMTP server port.")
	smtpSender := flag.String("smtp-sender", "", "SMTP sender.")

	filePath := flag.String("static-path", "", "Filepath to static files.")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := redis.New(ctx, *redisURI)
	if err != nil {
		logrus.Fatalf("failed to connect to database: %s", err)
	}

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		logrus.Fatalf("failed to open a connection: %s", err)
	}

	server, err := api.NewServer(ctx, db, emailClient(*smtpHost, *smtpPort, *smtpSender), *filePath, *domain)
	if err != nil {
		logrus.Fatalf("failed to create server: %s", err)
	}

	if err := server.Serve(ctx, lis, &api.TLSConfig{
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
