package vendor

//go:generate protoc -I=../../proto/third_party --go_out=plugins=grpc:. --grpc-gateway_out=:. ../../proto/third_party/google/longrunning/operations.proto
