#!/bin/bash

set -e

echo "Building TLS Fingerprint Requester..."

# Create bin directory
mkdir -p bin

# Build flags for size optimization
LDFLAGS="-s -w"

# Linux AMD64
echo "Building for Linux AMD64..."
GOOS=linux GOARCH=amd64 go build -ldflags="$LDFLAGS" -o bin/fingerprint_linux_amd64

# Windows AMD64
echo "Building for Windows AMD64..."
GOOS=windows GOARCH=amd64 go build -ldflags="$LDFLAGS" -o bin/fingerprint_windows_amd64.exe

# Android ARM64
echo "Building for Android ARM64..."
GOOS=android GOARCH=arm64 go build -ldflags="$LDFLAGS" -o bin/fingerprint_android_arm64

echo "Build completed successfully!"
echo "Binaries are in the bin/ directory"
