#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
go run github.com/evilmartians/lefthook/v2 "$@"
