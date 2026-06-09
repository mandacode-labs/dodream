#!/usr/bin/env bash
set -euo pipefail

LOCALBIN="${1:-$(pwd)/bin}"
shift 2>/dev/null || true
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

"$SCRIPT_DIR/install-golangci-lint.sh" "$LOCALBIN"

VERSION="v2.12.2"

# Prefer system binary if version matches
LINT_BIN=""
if command -v golangci-lint &>/dev/null; then
  SYSTEM_VER=$(golangci-lint version --format=short 2>/dev/null || echo "unknown")
  if [ "$SYSTEM_VER" = "$VERSION" ]; then
    LINT_BIN="golangci-lint"
  fi
fi
if [ -z "$LINT_BIN" ]; then
  LINT_BIN="$LOCALBIN/golangci-lint-$VERSION"
fi

if [ $# -eq 0 ]; then
  set -- run
fi
exec "$LINT_BIN" "$@"
