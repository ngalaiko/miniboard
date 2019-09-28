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
	"miniboard.app/storage/bolt"
)

var (
	domain = flag.String("domain", "http://localhost:8080", "Service domain.")
	addr   = flag.String("addr", ":8080", "Address to listen for connections.")

	boltPath = flag.String("bolt-path", "./bolt.db", "Path to the bolt storage.")

	sslCert = flag.String("ssl-cert", "", "Path to ssl certificate.")
	sslKey  = flag.String("ssl-key", "", "Path to ssl key.")

	smtpHost   = flag.String("smtp-host", "", "SMTP server host.")
	smtpPort   = flag.String("smtp-port", "", "SMTP server port.")
	smtpSender = flag.String("smtp-sender", "", "SMTP sender.")
)

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := bolt.New(ctx, *boltPath)
	if err != nil {
		logrus.Fatalf("failed to create bolt database: %s", err)
	}

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		logrus.Fatalf("failed to open a connection: %s", err)
	}

	server := api.NewServer(ctx, db, emailClient(), *domain)
	if err := server.Serve(ctx, lis, &api.TLSConfig{
		CertPath: *sslCert,
		KeyPath:  *sslKey,
	}); err != nil {
		logrus.Fatalf("failed to start the server: %s", err)
	}
}

func emailClient() email.Client {
	if *smtpHost == "" {
		return disabled.New()
	}
	return smtp.New(
		*smtpHost,
		*smtpPort,
		*smtpSender,
		os.Getenv("SMTP_USERNAME"),
		os.Getenv("SMTP_PASSWORD"),
	)
}
