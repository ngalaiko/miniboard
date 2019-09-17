#!/bin/bash

set -eo pipefail

if [ -z ${VERSION} ]; then
    VERSION=$(git rev-parse --short HEAD)_$(uname -s)_$(uname -m)
fi

echo "MINIBOARD_VERSION ${VERSION}"
