#!/bin/sh
ROOT_DIR=$(dirname "$0")

export COPYRIGHT="2024, MainDen"
export LICENSE="BSD-3-Clause License"
export URL="https://github.com/MainDen/stdhttp"

export APP=stdhttp
export CGO_ENABLED=0

mkdir -p "$ROOT_DIR/.build"

# linux
export GOOS=linux
export GOARCH=amd64
tar -cf "$ROOT_DIR/.build/$APP-$GOOS-$GOARCH.tar" -C "$ROOT_DIR" "LICENSE.md" || exit 1
tar -rf "$ROOT_DIR/.build/$APP-$GOOS-$GOARCH.tar" -C "$ROOT_DIR" "README.md" || exit 1

# linux console
export LDFLAGS=""
export APP_SUFFIX=""
export APP_EXT=""
"$ROOT_DIR/scripts/build.sh" || exit 1
tar -rf "$ROOT_DIR/.build/$APP-$GOOS-$GOARCH.tar" -C "$ROOT_DIR/.build/$GOOS"_"$GOARCH" "$APP$APP_SUFFIX$APP_EXT" || exit 1
gzip -f "$ROOT_DIR/.build/$APP-$GOOS-$GOARCH.tar" || exit 1

# windows
export GOOS=windows
export GOARCH=amd64
tar -cf "$ROOT_DIR/.build/$APP-$GOOS-$GOARCH.tar" -C "$ROOT_DIR" "LICENSE.md" || exit 1
tar -rf "$ROOT_DIR/.build/$APP-$GOOS-$GOARCH.tar" -C "$ROOT_DIR" "README.md" || exit 1

# windows console
export LDFLAGS=""
export APP_SUFFIX=""
export APP_EXT=".exe"
"$ROOT_DIR/scripts/build.sh" || exit 1
tar -rf "$ROOT_DIR/.build/$APP-$GOOS-$GOARCH.tar" -C "$ROOT_DIR/.build/$GOOS"_"$GOARCH" "$APP$APP_SUFFIX$APP_EXT" || exit 1

# widows daemon
export LDFLAGS="-H windowsgui"
export APP_SUFFIX="d"
export APP_EXT=".exe"
"$ROOT_DIR/scripts/build.sh" || exit 1
tar -rf "$ROOT_DIR/.build/$APP-$GOOS-$GOARCH.tar" -C "$ROOT_DIR/.build/$GOOS"_"$GOARCH" "$APP$APP_SUFFIX$APP_EXT" || exit 1
gzip -f "$ROOT_DIR/.build/$APP-$GOOS-$GOARCH.tar" || exit 1
