package feeds

//go:generate protoc -I=../../proto/users/feeds/v1 -I=../../proto/third_party ../../proto/users/feeds/v1/feeds_service.proto --go_out=plugins=grpc:. --grpc-gateway_out=:.
