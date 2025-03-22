#!/bin/bash
set -e

# Configuration 
GITHUB_REPO="SticketInya/kredentials"
CLI_NAME="kredentials"
INSTALL_DIR="/usr/local/bin/"

# Detect platform
OS=$(uname -s)
ARCH=$(uname -m)

# Create temp directory
TMP_DIR=$(mktemp -d)
cd $TMP_DIR

# Download latest release
echo "Downloading latest $CLI_NAME..."
DOWNLOAD_URL="https://github.com/$GITHUB_REPO/releases/latest/download/${CLI_NAME}_${OS}_${ARCH}.tar.gz"
curl -q --progress-bar -sL "$DOWNLOAD_URL" -o "${CLI_NAME}.tar.gz"

# Extract and install
echo "Extracting archive..."
tar -xzf "${CLI_NAME}.tar.gz"

echo "Installing binary..."
chmod +x "$CLI_NAME"
sudo mv "$CLI_NAME" "$INSTALL_DIR/$CLI_NAME"

# Clean up
cd - > /dev/null
rm -rf "$TMP_DIR"

echo "$CLI_NAME installed successfully!"