#!/bin/sh


SRC="${SUB_PROJECT_DIR}/src"
DIST="${SUB_PROJECT_DIR}/build"

echo "INFO: Installing Go dependencies..."
go -C "${SRC}" mod tidy
go -C "${SRC}" mod download

echo "INFO: Cleaning build directory..."
rm -rf "$DIST"
mkdir -p "$DIST"

echo "INFO: Building for Linux amd64..."
GOOS=linux GOARCH=amd64 go -C "${SRC}" build -o "$DIST/ops" .
chmod +x "${DIST}/ops"

echo "INFO: Building for Windows amd64..."
GOOS=windows GOARCH=amd64 go -C "${SRC}" build -o "$DIST/ops.exe" .

echo "INFO: Creating distribution archive..."
zip -j "${DIST}/${CI_PROJECT_NAME}-v${IMAGE_VERSION}.zip" "${DIST}/ops" "${DIST}/ops.exe"
