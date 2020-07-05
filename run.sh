#!/bin/sh

WORK_DIR="/home/e2e"

echo "Building Simple-E2E binary..."
cd "$FRAMEWORK_PATH" || return
go build -o ./simple-e2e
mv ./simple-e2e $WORK_DIR
cd "$WORK_DIR" || return

echo "Welcome to Simple-E2E container!"
/bin/sh