package main

import (
	"context"
	"flag"
	"net"

	"github.com/sirupsen/logrus"
	"miniboard.app/api"
	"miniboard.app/storage/bolt"
)

var (
	boltPath = flag.String("bolt-path", "./bolt.db", "Path to the bolt storage.")
	addr     = flag.String("addr", ":8080", "Address to listen for connections.")
)

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := bolt.New(ctx, *boltPath)
	if err != nil {
		logrus.Panicf("failed to create bolt database: %s", err)
	}

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		logrus.Panicf("failed to open a connection: %s", err)
	}

	server := api.NewServer(ctx, db)
	if err := server.Serve(ctx, lis); err != nil {
		logrus.Panicf("failed to start the server: %s", err)
	}
}
