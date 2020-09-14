package sources

//go:generate protoc -I=../../proto/sources/v1 -I=../../proto/third_party ../../proto/sources/v1/sources_service.proto --go_out=plugins=grpc:. --grpc-gateway_out=:.
