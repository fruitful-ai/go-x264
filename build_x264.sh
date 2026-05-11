#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
LIBDIR="$SCRIPT_DIR/internal/x264/lib"
INCDIR="$SCRIPT_DIR/internal/x264/include"

# Already built locally?
if [ -f "$LIBDIR/libx264.a" ]; then
    echo "x264 already built"
    exit 0
fi

# Try to download a pre-built release asset
ARCH=$(uname -m)
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$ARCH" in
    x86_64)  ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
esac
PLATFORM="${OS}_${ARCH}"
REPO="${X264_REPO:-fruitfulch/go-x264}"
ASSET_URL="https://github.com/$REPO/releases/latest/download/x264-${PLATFORM}.tar.gz"

echo "Attempting to download x264 for $PLATFORM ..."
mkdir -p "$LIBDIR"
if curl -fsL "$ASSET_URL" | tar -xz -C "$SCRIPT_DIR/internal/x264" 2>/dev/null; then
    echo "x264 downloaded for $PLATFORM"
    exit 0
fi
rm -rf "$SCRIPT_DIR/internal/x264"
echo "Download failed, building from source ..."

# Build from source
echo "Downloading x264 source ..."
TMPDIR=$(mktemp -d)
trap "rm -rf $TMPDIR" EXIT

curl -sL "https://code.videolan.org/videolan/x264/-/archive/master/x264-master.tar.bz2" | tar -xj -C "$TMPDIR"

cd "$TMPDIR/x264-master"

echo "Configuring x264 ..."
./configure \
    --enable-static \
    --enable-pic \
    --disable-asm \
    --disable-cli \
    --prefix="$SCRIPT_DIR/internal/x264"

echo "Building x264 ..."
make -j"$(nproc)" install

echo "x264 built successfully"
