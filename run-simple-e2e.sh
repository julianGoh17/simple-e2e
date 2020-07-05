#!/bin/sh

WORKDIR="$(pwd)"
SIMPLE_E2E_PATH="$WORKDIR/framework"
set -e

echo "Building Simple-E2E container"
docker build ./ -t simple-e2e

docker run -it \
    -v "$SIMPLE_E2E_PATH":/go/src/github.com/julianGoh17/simple-e2e/framework \
    simple-e2e 
