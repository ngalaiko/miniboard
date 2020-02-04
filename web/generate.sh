#!/usr/bin/env bash

set -eou pipefail

rm -rf ./src/clients/proto/*

protoc \
	-I ../server/proto/users/articles/v1 \
	../server/proto/users/articles/v1/articles_service.proto \
	--js_out=import_style=commonjs:./src/clients/proto \
	--grpc-web_out=import_style=commonjs,mode=grpcwebtext:./src/clients/proto

protoc \
	-I ../server/proto/codes/v1 \
	../server/proto/codes/v1/codes_service.proto \
	--js_out=import_style=commonjs:./src/clients/proto \
	--grpc-web_out=import_style=commonjs,mode=grpcwebtext:./src/clients/proto

protoc \
	-I ../server/proto/users/sources/v1 \
	../server/proto/users/sources/v1/sources_service.proto \
	--js_out=import_style=commonjs:./src/clients/proto \
	--grpc-web_out=import_style=commonjs,mode=grpcwebtext:./src/clients/proto

protoc \
	-I ../server/proto/tokens/v1 \
	../server/proto/tokens/v1/tokens_service.proto \
	--js_out=import_style=commonjs:./src/clients/proto \
	--grpc-web_out=import_style=commonjs,mode=grpcwebtext:./src/clients/proto

protoc \
	-I ../server/proto/users/v1 \
	../server/proto/users/v1/users_service.proto \
	--js_out=import_style=commonjs:./src/clients/proto \
	--grpc-web_out=import_style=commonjs,mode=grpcwebtext:./src/clients/proto
