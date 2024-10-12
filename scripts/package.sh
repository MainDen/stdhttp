#!/bin/sh
ROOT_DIR=./$(realpath -s --relative-to="." "$(dirname "$0")/..")
. "$ROOT_DIR/scripts/utils.sh"

required "APP"
required "VERSION"
required "GOOS"
required "GOARCH"

BUILD_DIR=$ROOT_DIR/.build
TARGET_DIR_NAME=$GOOS"_"$GOARCH
TARGET_DIR_PATH=$BUILD_DIR/$TARGET_DIR_NAME
PACKAGE_NAME=$APP-$VERSION-$GOOS-$GOARCH
PACKAGE_FILE_NAME=$PACKAGE_NAME.tar
PACKAGE_FILE_PATH=$BUILD_DIR/$PACKAGE_FILE_NAME

info "Packaging: $PACKAGE_NAME"

[ -d "$BUILD_DIR" ] || error "Directory '.build' does not exist."
[ -d "$TARGET_DIR_PATH" ] || error "Directory '$TARGET_DIR_NAME' does not exist."
[ -f "$ROOT_DIR/LICENSE.md" ] || warn "File 'LICENSE.md' does not exist."
[ -f "$ROOT_DIR/README.md" ] || warn "File 'README.md' does not exist."
rm -f "$PACKAGE_FILE_PATH"

tar -cf "$PACKAGE_FILE_PATH" -C "$TARGET_DIR_PATH" "." || error "Failed to add build artifacts to file '$PACKAGE_FILE_NAME'."
[ -f "$ROOT_DIR/LICENSE.md" ] && tar -rf "$PACKAGE_FILE_PATH" -C "$ROOT_DIR" "LICENSE.md" || error "Failed to add 'LICENSE.md' to file '$PACKAGE_FILE_NAME'."
[ -f "$ROOT_DIR/README.md" ] && tar -rf "$PACKAGE_FILE_PATH" -C "$ROOT_DIR" "README.md" || error "Failed to add 'README.md' to file '$PACKAGE_FILE_NAME'."
gzip -f "$PACKAGE_FILE_PATH" || error "Failed to compress file '$PACKAGE_FILE_NAME'."

info "Packaging: DONE"
