package users

//go:generate protoc -I=../../../proto/users/v1 -I=../../../proto/third_party ../../../proto/users/v1/users_service.proto --go_out=plugins=grpc:. --grpc-gateway_out=:.
