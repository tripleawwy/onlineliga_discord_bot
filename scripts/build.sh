#!/bin/bash

for PLATFORM in "linux/amd64"; do
    GOOS=${PLATFORM%/*}
    GOARCH=${PLATFORM#*/}

    BIN_FILE_NAME="bin/ol_discord_bot-${GOOS}-${GOARCH}"

    echo Building ${GOOS}-${GOARCH} version ${OLS_VERSION}...

    GOOS="${GOOS}" GOARCH="${GOARCH}" \
    go build -o "${BIN_FILE_NAME}" \
    || FAILURES="${FAILURES} ${GOOS}-${GOARCH}"
done

if [ -n "${FAILURES}" ]; then
    echo Build Failures:
    for FAILURE in $FAILURES; do
        echo $FAILURE
    done
    exit 1
fi
