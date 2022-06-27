#!/bin/sh

# force go modules
export GOPATH=""

# disable cgo
export CGO_ENABLED=0

set -e
set -x

# linux
GOOS=linux GOARCH=amd64 go build -o release/linux-amd64-plugin
GOOS=linux GOARCH=arm64 go build -o release/linux-arm64-plugin
# mac
GOOS=darwin GOARCH=amd64 go build -o release/darwin-amd64-plugin
GOOS=darwin GOARCH=arm64 go build -o release/darwin-arm64-plugin
# windows
GOOS=windows GOARCH=amd64 go build -o release/plugin.exe
