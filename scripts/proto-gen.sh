#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

rm -rf pkg/proto/dodream
PATH="$(go env GOPATH)/bin:$PATH" buf generate
