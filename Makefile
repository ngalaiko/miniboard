linkify:
	./scripts/linkify.sh

update:
	bazel run //:gazelle

deps:
	go mod tidy
	bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro=repositories.bzl%go_repositories

build: linkify
	bazel build //cmd/miniboard:miniboard

run:
	bazel run //cmd/miniboard:miniboard

test:
	bazel test //...
