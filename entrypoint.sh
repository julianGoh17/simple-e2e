#!/bin/sh

set -x 
set -v 

WORK_DIR="/home/e2e"
SIMPLE_E2E_PATH="$GOPATH/src/github.com/julianGoh17/simple-e2e/"

export TEST_DIR="$WORK_DIR/tests"

echo "Building Simple-E2E binary..."
cd "$SIMPLE_E2E_PATH/framework" || return
sudo go build -o ./simple-e2e
mv ./simple-e2e $WORK_DIR
cd "$WORK_DIR" || return

# For github action
if [ "$1" = "unit-test" ]; then 
    echo "*** Running unit tests***"
    cd "$SIMPLE_E2E_PATH/framework" || return 
    sudo go get 
    sudo go test -cover ./... 
    if [ $? -ne 0 ]; then 
        exit 1
    fi
else 
    echo "*** Welcome to Simple-E2E container ***"
    /bin/sh
fi