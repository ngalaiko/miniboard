package genproto

//go:generate protoc -I=../../api --go_out=paths=source_relative:. --go-grpc_out=. --go-grpc_opt=paths=source_relative --grpc-gateway_out=paths=source_relative:. ../../api/articles/v1/articles_service.proto
//go:generate protoc -I=../../api --go_out=paths=source_relative:. --go-grpc_out=. --go-grpc_opt=paths=source_relative --grpc-gateway_out=paths=source_relative:. ../../api/codes/v1/codes_service.proto
//go:generate protoc -I=../../api --go_out=paths=source_relative:. --go-grpc_out=. --go-grpc_opt=paths=source_relative --grpc-gateway_out=paths=source_relative:. ../../api/feeds/v1/feeds_service.proto
//go:generate protoc -I=../../api                                                                                      --grpc-gateway_out=paths=source_relative:. ../../api/google/longrunning/operations.proto
//go:generate protoc -I=../../api --go_out=paths=source_relative:. --go-grpc_out=. --go-grpc_opt=paths=source_relative --grpc-gateway_out=paths=source_relative:. ../../api/sources/v1/sources_service.proto
//go:generate protoc -I=../../api --go_out=paths=source_relative:. --go-grpc_out=. --go-grpc_opt=paths=source_relative --grpc-gateway_out=paths=source_relative:. ../../api/tokens/v1/tokens_service.proto
