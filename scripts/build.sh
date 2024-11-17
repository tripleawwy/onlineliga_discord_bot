#!/bin/bash

# get PLATFORMS from the first argument with a default value
PLATFORM=${1:-"linux/amd64"}
GOOS=${PLATFORM%/*}
GOARCH=${PLATFORM#*/}

BIN_FILE_NAME="bin/ol_discord_bot"

echo "Building ${GOOS}-${GOARCH}..."

GOOS="${GOOS}" GOARCH="${GOARCH}" \
go build -ldflags="-w -s" -o "${BIN_FILE_NAME}" \
|| FAILURES="${FAILURES} ${GOOS}-${GOARCH}"

if [ -n "${FAILURES}" ]; then
  echo Build Failures:
  for FAILURE in $FAILURES; do
    echo "${FAILURE}"
  done
  exit 1
fi
