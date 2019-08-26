#! /usr/bin/env bash

set -euo pipefail

go mod tidy
bazel run //:gazelle -- \
    update-repos \
        -from_file=go.mod \
        -to_macro=repositories.bzl%go_repositories
bazel run //:gazelle
