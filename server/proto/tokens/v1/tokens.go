package tokens

//go:generate protoc -I=../../../../proto/tokens/v1 -I=../../../../proto/third_party ../../../../proto/tokens/v1/tokens_service.proto --go_out=plugins=grpc:. --grpc-gateway_out=:.
