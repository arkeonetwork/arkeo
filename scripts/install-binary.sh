#!/bin/bash

# GitHub repository information
REPO="arkeonetwork/arkeo"

# Fetch the latest release version from GitHub
LATEST_VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | jq -r '.tag_name')

if [ -z "$LATEST_VERSION" ]; then
    echo "Failed to fetch the latest version from GitHub."
    exit 1
fi

echo "Latest version is $LATEST_VERSION"

# Determine the platform
ARCH=$(uname -m)
OS=$(uname -s | tr '[:upper:]' '[:lower:]')

# Set the download URL based on the architecture
if [[ "$OS" == "darwin" && "$ARCH" == "arm64" ]]; then
    BINARY="arkeod_darwin_arm64-testnet"
elif [[ "$OS" == "linux" && "$ARCH" == "x86_64" ]]; then
    BINARY="arkeod_linux_amd64-testnet"
else
    echo "Unsupported platform: $OS $ARCH"
    exit 1
fi

# Construct the download URL
URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/${BINARY}"

# Download the binary
echo "Downloading $BINARY..."
curl -L -o arkeod $URL

# Copy the binary to /usr/local/bin
echo "Copying arkeod to /usr/local/bin..."
sudo cp arkeod /usr/local/bin/arkeod

# Set execute permissions
echo "Setting execute permissions..."
sudo chmod +x /usr/local/bin/arkeod

# Check the version of the installed binary
INSTALLED_VERSION=$(/usr/local/bin/arkeod version)

# Verify the installed version
if [[ "$INSTALLED_VERSION" == "$LATEST_VERSION" ]]; then
    echo "Version check passed. Installed version: $INSTALLED_VERSION"
else
    echo "Version check failed. Installed version: $INSTALLED_VERSION, but expected: $LATEST_VERSION"
    exit 1
fi

# Clean up downloaded binary
rm arkeod

echo "Done."