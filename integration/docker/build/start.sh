#!/bin/bash

echo "Building binaries..."

set -e
set -x

GOBIN=/usr/local/go/bin/go

$GOBIN build -o ./build/linux/script-eth-rpc-adaptor ./cmd/script-eth-rpc-adaptor

set +x 

echo "Done."



