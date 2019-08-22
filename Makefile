linkify:
	./scripts/linkify.sh

build: linkify
	bazel build //cmd/miniboard:miniboard

run:
	bazel run //cmd/miniboard:miniboard

test:
	bazel test //...
