package codes

//go:generate protoc -I=../../proto/codes/v1 -I=../../proto/third_party ../../proto/codes/v1/codes_service.proto --go_out=plugins=grpc:. --grpc-gateway_out=:.
