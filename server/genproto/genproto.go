package genproto

//go:generate protoc -I=../../proto --go_out=paths=source_relative:. --go-grpc_out=. --go-grpc_opt=paths=source_relative --grpc-gateway_out=paths=source_relative:. ../../proto/articles/v1/articles_service.proto
//go:generate protoc -I=../../proto --go_out=paths=source_relative:. --go-grpc_out=. --go-grpc_opt=paths=source_relative --grpc-gateway_out=paths=source_relative:. ../../proto/codes/v1/codes_service.proto
//go:generate protoc -I=../../proto --go_out=paths=source_relative:. --go-grpc_out=. --go-grpc_opt=paths=source_relative --grpc-gateway_out=paths=source_relative:. ../../proto/feeds/v1/feeds_service.proto
//go:generate protoc -I=../../proto --grpc-gateway_out=paths=source_relative:. ../../proto/google/longrunning/operations.proto
//go:generate protoc -I=../../proto --go_out=paths=source_relative:. --go-grpc_out=. --go-grpc_opt=paths=source_relative --grpc-gateway_out=paths=source_relative:. ../../proto/sources/v1/sources_service.proto
//go:generate protoc -I=../../proto --go_out=paths=source_relative:. --go-grpc_out=. --go-grpc_opt=paths=source_relative --grpc-gateway_out=paths=source_relative:. ../../proto/tokens/v1/tokens_service.proto
