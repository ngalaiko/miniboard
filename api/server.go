package api // import "miniboard.app/api"

import (
	"context"
	"net"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	usersservice "miniboard.app/api/users"
	"miniboard.app/proto/users/v1"
	"miniboard.app/storage"
)

// Server is the api server.
type Server struct {
	grpcServer *grpc.Server
}

// NewServer creates new api server.
func NewServer(ctx context.Context, db storage.DB) *Server {
	grpcServer := grpc.NewServer()

	users.RegisterUsersServiceServer(grpcServer, usersservice.New(db))

	reflection.Register(grpcServer)

	return &Server{
		grpcServer: grpcServer,
	}
}

// Serve starts the server.
func (s *Server) Serve(ctx context.Context, lis net.Listener) error {
	logrus.Infof("[gRPC] starting server on %s", lis.Addr())

	if err := s.grpcServer.Serve(lis); err != nil {
		return errors.Wrap(err, "failed to start gRPC server")
	}
	go func() {
		<-ctx.Done()
		logrus.Infof("[gRPC] stopping server")
		s.grpcServer.GracefulStop()
	}()

	return nil
}
