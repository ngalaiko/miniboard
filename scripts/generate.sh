#! /usr/bin/env bash

set -euo pipefail

go mod tidy
bazel run //server:gazelle -- \
    update-repos \
        -from_file=./server/go.mod \
        -to_macro=./server/repositories.bzl%go_repositories
bazel run //server:gazelle
