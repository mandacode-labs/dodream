#!/usr/bin/env bash
set -euo pipefail

echo "Installing protoc-gen-go..." >&2
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

echo "Installing protoc-gen-go-grpc..." >&2
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
