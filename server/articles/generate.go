package articles

//go:generate protoc -I=../../proto/articles/v1 -I=../../proto/third_party ../../proto/articles/v1/articles_service.proto --go_out=plugins=grpc:. --grpc-gateway_out=:.
