#!/bin/sh
ROOT_DIR=./$(realpath -s --relative-to="." "$(dirname "$0")")
. "$ROOT_DIR/scripts/utils.sh"

export APP=stdhttp
export COPYRIGHT="2024, MainDen"
export LICENSE="BSD-3-Clause License"
export URL="https://github.com/MainDen/stdhttp"
export CGO_ENABLED=0

[ -z "$VERSION" ] && export VERSION=$(git describe --tags --always --dirty)
required "VERSION"

BUILD=$ROOT_DIR/scripts/build.sh
PACKAGE=$ROOT_DIR/scripts/package.sh

# linux
GOOS="linux" GOARCH="amd64" LDFLAGS="" APP_SUFFIX="" APP_EXT="" $BUILD || error "Failed to build '$APP'."
GOOS="linux" GOARCH="amd64" $PACKAGE || error "Failed to package '$APP'."

# windows
GOOS="windows" GOARCH="amd64" LDFLAGS="" APP_SUFFIX="" APP_EXT=".exe" $BUILD || error "Failed to build '$APP'."
GOOS="windows" GOARCH="amd64" LDFLAGS="-H windowsgui" APP_SUFFIX="d" APP_EXT=".exe" $BUILD || error "Failed to build '$APP'."
GOOS="windows" GOARCH="amd64" $PACKAGE || error "Failed to package '$APP'."
