#!/data/data/com.termux/files/usr/bin/bash

set -e

echo "Building TLS Fingerprint Requester for Termux..."

# Download dependencies
echo "Running go mod tidy..."
go mod tidy

# Create bin directory
mkdir -p bin

# Build flags for size optimization
LDFLAGS="-s -w"

# Android ARM64
echo "Building for Android ARM64..."
GOOS=android GOARCH=arm64 go build -ldflags="$LDFLAGS" -o bin/fingerprint_android_arm64 ./cmd/tlsRequester

echo "Build completed successfully!"
echo "Binary: bin/fingerprint_android_arm64"
echo ""
echo "Usage: ./bin/fingerprint_android_arm64 < request.json"
