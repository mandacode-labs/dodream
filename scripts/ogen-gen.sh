#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."

rm -rf pkg/oas

go tool ogen --clean api/openapi/v1/openapi.yaml

mkdir -p pkg/oas
mv api/openapi/v1/oas_*.go pkg/oas/
sed -i 's/package api/package oas/g' pkg/oas/*.go
