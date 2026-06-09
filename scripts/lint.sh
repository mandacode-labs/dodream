#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
if [ $# -eq 0 ]; then
  set -- run
fi
go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint "$@"
