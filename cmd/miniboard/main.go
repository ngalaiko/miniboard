package main // import "miniflux.app/cmd"

import (
	"context"
	"flag"
	"net"
	"os"

	"github.com/sirupsen/logrus"
	"miniboard.app/application/api"
	"miniboard.app/application/storage/bolt"
)

var (
	boltPath = flag.String("bolt-path", "./bolt.db", "Path to the bolt storage.")
	addr     = flag.String("addr", ":8080", "Address to listen for connections.")
)

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	file, err := os.OpenFile(*boltPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		logrus.Panicf("failed to open bolt file: %s", err)
	}

	db, err := bolt.New(ctx, file.Name())
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
