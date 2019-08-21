package proto // import "miniflux.app/proto"

//go:generate protoc -I./third_party/googleapis -I. --go_out=plugins=grpc:. ./users/v1/users_service.proto
//go:generate protoc -I./third_party/googleapis -I. --go_out=plugins=grpc:. ./auth/v1/auth_service.proto
//go:generate protoc -I./third_party/googleapis -I. --go_out=plugins=grpc:. ./users/articles/v1/articles_service.proto
