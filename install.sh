#!/usr/bin/env bash

set -e

REPO="moxus/secrets-cli"
BINARY_NAME="secrets-cli"

# Detect OS
OS="$(uname | tr '[:upper:]' '[:lower:]')"
case "$OS" in
  linux)   OS="linux" ;;
  darwin)  OS="darwin" ;;
  *)       echo "Unsupported OS: $OS"; exit 1 ;;
esac

# Detect ARCH
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *)            echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

echo "https://api.github.com/repos/$REPO/releases/latest"

# Get latest release tag from GitHub API
TAG=$(curl -sSL "https://github.com/$REPO/releases/latest" \
  | sed -n 's|.*href="/'"$REPO"'/releases/tag/\([^"]*\)".*|\1|p' | head -n1)

if [ -z "$TAG" ]; then
  echo "Could not find latest release tag."
  exit 1
fi

# Construct download URL
BIN="${BINARY_NAME}-${OS}-${ARCH}"
ASSET="${BIN}.zip"
URL="https://github.com/$REPO/releases/download/$TAG/$ASSET"

echo "Downloading $URL..."

TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

# Download and unzip
curl -sSL -o "$ASSET" "$URL"
unzip "$ASSET"
mv "$BIN" "$BINARY_NAME"

# Move binary to /usr/local/bin
chmod +x "$BINARY_NAME"
sudo mv "$BINARY_NAME" /usr/local/bin/

echo "Installed $BINARY_NAME to /usr/local/bin"

# Cleanup
cd /
rm -rf "$TMP_DIR"
