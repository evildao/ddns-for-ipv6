#!/bin/bash

APP_NAME="ddns-ipv6"
VERSION="1.0.0"
COMMIT_HASH=$(git rev-parse --short HEAD)
BUILD_DATE=$(date +'%Y-%m-%d_%H:%M:%S')

PLATFORMS=(
    "linux/amd64"
    "windows/amd64"
    "darwin/amd64"
)

for platform in "${PLATFORMS[@]}"; do
    GOOS=${platform%/*}
    GOARCH=${platform#*/}
    output_name="${APP_NAME}_${GOOS}_${VERSION}"

    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    env GOOS=$GOOS GOARCH=$GOARCH go build \
        -ldflags "\
		-X main.version=$VERSION \
		-X main.commit=$COMMIT_HASH \
		-X main.date=$BUILD_DATE" \
        -o bin/$GOOS-$GOARCH/$output_name

    echo "Built $GOOS/$GOARCH"
done
