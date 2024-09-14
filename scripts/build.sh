#!/bin/sh
ROOT_DIR=$(dirname "$0")/..
source "$ROOT_DIR/scripts/utils.sh"

required "APP"
required "GOOS"
required "GOARCH"
required "CGO_ENABLED"

[ -z "$APP" ] || append "LDFLAGS" "-X 'main.appname=$APP$APP_SUFFIX'"
[ -z "$VERSION" ] && VERSION="$(cat "$ROOT_DIR/VERSION")"
[ -z "$VERSION" ] || append "LDFLAGS" "-X 'main.version=$VERSION'"
[ -z "$COPYRIGHT" ] || append "LDFLAGS" "-X 'main.copyright=$COPYRIGHT'"
[ -z "$LICENSE" ] || append "LDFLAGS" "-X 'main.license=$LICENSE'"
[ -z "$URL" ] || append "LDFLAGS" "-X 'main.url=$URL'"

echo "Building: $APP$APP_SUFFIX for $GOOS/$GOARCH (CGO_ENABLED=$CGO_ENABLED ldflags=$LDFLAGS)"
GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=$CGO_ENABLED go build -ldflags="$LDFLAGS" -o "$ROOT_DIR/.build/$GOOS"_"$GOARCH/$APP$APP_SUFFIX$APP_EXT" "$ROOT_DIR/cmd/$APP" || exit 1
echo "Building: DONE"
