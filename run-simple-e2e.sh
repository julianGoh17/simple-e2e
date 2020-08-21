#!/bin/sh

WORKDIR="$(pwd)"
SIMPLE_E2E_PATH="$WORKDIR/framework"
TEST_FILES_PATH="$WORKDIR/tests"
set -e

echo "Building Simple-E2E container"
docker build ./ -t simple-e2e

docker run -it \
    -v "$SIMPLE_E2E_PATH":/go/src/github.com/julianGoh17/simple-e2e/framework \
    -v "$TEST_FILES_PATH":/home/e2e/tests \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -e FRAMEWORK_PATH=/go/src/github.com/julianGoh17/simple-e2e/framework \
    simple-e2e 
