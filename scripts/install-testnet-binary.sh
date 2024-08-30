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
     CHECKSUM_URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/arkeod_${LATEST_VERSION}__testnet_cross_checksums.txt"
    BINARY="arkeod_darwin_arm64-testnet"
    
elif [[ "$OS" == "linux" && "$ARCH" == "x86_64" ]]; then
    CHECKSUM_URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/arkeod_${LATEST_VERSION}__testnet_checksums.txt"
    BINARY="arkeod_linux_amd64-testnet"
    
else
    echo "Unsupported platform: $OS $ARCH"
    exit 1
fi

# Construct the download URL
URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/${BINARY}"

# Download the binary
echo "Downloading $BINARY..."
curl -L -o $BINARY $URL


# Download the checksum file
echo "Downloading checksum file..."
curl -L -o checksums.txt $CHECKSUM_URL


# Calculate the downloaded binary's checksum
echo "Verifying checksum..."
if [[ "$OS" == "darwin" ]]; then
    DOWNLOAD_CHECKSUM=$(shasum -a 256 $BINARY | awk '{ print $1 }')
else
    DOWNLOAD_CHECKSUM=$(sha256sum $BINARY | awk '{ print $1 }')
fi

# Extract the expected checksum for the downloaded binary from the checksums.txt file
EXPECTED_CHECKSUM=$(grep "$BINARY" checksums.txt |  grep -v '\.zip' | awk '{ print $1 }')

if [ "$DOWNLOAD_CHECKSUM" != "$EXPECTED_CHECKSUM" ]; then
    echo "Checksum verification failed. The downloaded binary is corrupted."
    rm arkeod checksums.txt
    exit 1
fi

echo "Checksum verification passed."

# Copy the binary to /usr/local/bin
echo "Copying arkeod to ${HOME}/go/bin..."
cp $BINARY ${HOME}/go/bin/arkeod

# Set execute permissions
echo "Setting execute permissions..."
chmod +x ${HOME}/go/bin/arkeod


# Clean up downloaded binary
rm $BINARY checksums.txt

echo "Done."