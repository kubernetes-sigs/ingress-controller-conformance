#! /usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

BINDATA=${BINDATA:-go run github.com/go-bindata/go-bindata/go-bindata}

exec $BINDATA "$@"

