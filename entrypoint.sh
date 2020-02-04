#!/usr/bin/env sh

set -eou pipefail

./miniboard --static-path ./dist "$@"
