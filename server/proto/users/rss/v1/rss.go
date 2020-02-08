package rss

//go:generate protoc -I=../../../../../proto/users/rss/v1 -I=../../../../../proto/third_party ../../../../../proto/users/rss/v1/rss_service.proto --go_out=plugins=grpc:. --grpc-gateway_out=:.
