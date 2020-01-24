package api

import (
	"context"
	"strings"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
	grpcCodes "google.golang.org/grpc/codes"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"miniboard.app/api/actor"
	articlesservice "miniboard.app/articles"
	codesservice "miniboard.app/codes"
	"miniboard.app/email"
	"miniboard.app/images"
	"miniboard.app/jwt"
	"miniboard.app/proto/codes/v1"
	"miniboard.app/proto/tokens/v1"
	"miniboard.app/proto/users/articles/v1"
	"miniboard.app/proto/users/v1"
	"miniboard.app/storage"
	tokensservice "miniboard.app/tokens"
	usersservice "miniboard.app/users"
)

const authCookie = "auth"

func grpcServer(db storage.Storage, emailClient email.Client, jwtService *jwt.Service, images *images.Service, domain string) *grpc.Server {
	logrusEntry := log("grpc")
	grpc_logrus.ReplaceGrpcLogger(logrusEntry)

	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.UnaryServerInterceptor(logrusEntry),
			authorize(jwtService),
		),
	)

	articles.RegisterArticlesServiceServer(grpcServer, articlesservice.New(db, images))
	users.RegisterUsersServiceServer(grpcServer, usersservice.New(db))
	codes.RegisterCodesServiceServer(grpcServer, codesservice.New(domain, emailClient, jwtService))
	tokens.RegisterTokensServiceServer(grpcServer, tokensservice.New(jwtService))

	return grpcServer
}

func authorize(jwtService *jwt.Service) grpc.UnaryServerInterceptor {
	whitelisted := map[string]bool{
		"/app.miniboard.codes.v1.CodesService/CreateCode":    true,
		"/app.miniboard.tokens.v1.TokensService/CreateToken": true,
	}
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if whitelisted[info.FullMethod] {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.New(grpcCodes.Unauthenticated, "auth cookie missing").Err()
		}

		cookies, ok := md["cookie"]
		if !ok {
			return nil, status.New(grpcCodes.Unauthenticated, "auth cookie missing").Err()
		}

		for _, cookie := range strings.Split(cookies[0], "; ") {
			parts := strings.Split(cookie, "=")

			if parts[0] != authCookie {
				continue
			}

			subject, err := jwtService.Validate(ctx, parts[1], "access_token")
			if err != nil {
				return nil, status.New(grpcCodes.PermissionDenied, "invalid auth token").Err()
			}

			return handler(actor.NewContext(ctx, subject), req)
		}

		return nil, status.New(grpcCodes.Unauthenticated, "auth cookie missing").Err()
	}
}
