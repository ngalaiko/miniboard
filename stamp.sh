#!/bin/bash
set -eo pipefail

if [ ! -z "${API_URL}" ]; then
    echo "API_URL $API_URL"
fi

echo "VERSION $(git rev-parse HEAD)"
