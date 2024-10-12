#!/bin/sh
ROOT_DIR=./$(realpath -s --relative-to="." "$(dirname "$0")/..")
. "$ROOT_DIR/scripts/utils.sh"

required "APP"
required "VERSION"
required "GOOS"
required "GOARCH"
required "CGO_ENABLED"

APP_NAME=$APP$APP_SUFFIX
SOURCE_DIR=$ROOT_DIR/cmd/$APP
BUILD_DIR=$ROOT_DIR/.build
TARGET_DIR_NAME=$GOOS"_"$GOARCH
TARGET_DIR_PATH=$BUILD_DIR/$TARGET_DIR_NAME
TARGET_FILE_NAME=$APP_NAME$APP_EXT
TARGET_FILE_PATH=$TARGET_DIR_PATH/$TARGET_FILE_NAME

[ -z "$APP_NAME" ] || append "LDFLAGS" "-X 'main.appname=$APP_NAME'"
[ -z "$VERSION" ] || append "LDFLAGS" "-X 'main.version=$VERSION'"
[ -z "$COPYRIGHT" ] || append "LDFLAGS" "-X 'main.copyright=$COPYRIGHT'"
[ -z "$LICENSE" ] || append "LDFLAGS" "-X 'main.license=$LICENSE'"
[ -z "$URL" ] || append "LDFLAGS" "-X 'main.url=$URL'"

info "Building: $TARGET_FILE_NAME version $VERSION for $GOOS/$GOARCH (CGO_ENABLED=$CGO_ENABLED ldflags=$LDFLAGS)"

[ -d "$SOURCE_DIR" ] || error "Source directory '$APP' does not exist."
mkdir -p "$BUILD_DIR" || error "Failed to create directory '.build'."
mkdir -p "$TARGET_DIR_PATH" || error "Failed to create directory '$TARGET_DIR_NAME'."

GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=$CGO_ENABLED go build -ldflags="$LDFLAGS" -o "$TARGET_FILE_PATH" "$SOURCE_DIR" || error "Failed to build '$TARGET_FILE_NAME'."

info "Building: DONE"
