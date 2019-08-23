.PHONY: build

linkify:
	./scripts/linkify.sh

generate:
	go mod tidy
	bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro=repositories.bzl%go_repositories
	bazel run //:gazelle
 
build:
	bazel build //cmd/miniboard:miniboard

run:
	bazel run //cmd/miniboard:miniboard

test:
	bazel test --test_output=errors //...
