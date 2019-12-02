package api

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
	articlesservice "miniboard.app/articles"
	codesservice "miniboard.app/codes"
	"miniboard.app/email"
	"miniboard.app/jwt"
	"miniboard.app/proto/codes/v1"
	"miniboard.app/proto/tokens/v1"
	"miniboard.app/proto/users/articles/v1"
	"miniboard.app/proto/users/v1"
	"miniboard.app/storage"
	tokensservice "miniboard.app/tokens"
	usersservice "miniboard.app/users"
)

func grpcServer(db storage.Storage, emailClient email.Client, jwtService *jwt.Service, domain string) *grpc.Server {
	logrusEntry := log("grpc")
	grpc_logrus.ReplaceGrpcLogger(logrusEntry)

	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.UnaryServerInterceptor(logrusEntry),
		),
	)

	articles.RegisterArticlesServiceServer(grpcServer, articlesservice.New(db))
	users.RegisterUsersServiceServer(grpcServer, usersservice.New(db))
	codes.RegisterCodesServiceServer(grpcServer, codesservice.New(domain, emailClient, jwtService))
	tokens.RegisterTokensServiceServer(grpcServer, tokensservice.New(jwtService))

	return grpcServer
}
