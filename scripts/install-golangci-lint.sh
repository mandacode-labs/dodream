#!/usr/bin/env bash
set -euo pipefail

VERSION="v2.12.2"
LOCALBIN="${1:-$(pwd)/bin}"
TARGET="$LOCALBIN/golangci-lint-$VERSION"

if [ -f "$TARGET" ]; then
  echo "golangci-lint $VERSION already installed at $TARGET" >&2
  exit 0
fi

# Check system binary
if command -v golangci-lint &>/dev/null; then
  SYSTEM_VER=$(golangci-lint version --format=short 2>/dev/null || echo "unknown")
  if [ "$SYSTEM_VER" = "$VERSION" ]; then
    echo "Using system golangci-lint $VERSION" >&2
    exit 0
  fi
fi

# Detect platform
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  aarch64) ARCH="arm64" ;;
  armv*)   ARCH="armv7" ;;
esac

VER_NOV="${VERSION#v}"
URL="https://github.com/golangci/golangci-lint/releases/download/$VERSION/golangci-lint-$VER_NOV-$OS-$ARCH.tar.gz"

echo "Downloading golangci-lint $VERSION ($OS/$ARCH)..." >&2
mkdir -p "$LOCALBIN"
curl -sSfL "$URL" -o /tmp/golangci-lint.tar.gz
tar xzf /tmp/golangci-lint.tar.gz -C /tmp/
mv "/tmp/golangci-lint-$VER_NOV-$OS-$ARCH/golangci-lint" "$TARGET"
rm -f /tmp/golangci-lint.tar.gz
echo "Installed golangci-lint $VERSION at $TARGET" >&2
