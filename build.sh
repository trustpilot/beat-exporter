#!/bin/bash

set -e

if [[ -z "$GITHUB_WORKSPACE" ]]; then
  GITHUB_WORKSPACE=$(pwd)
  echo "Setting up GITHUB_WORKSPACE to current directory: ${GITHUB_WORKSPACE}"
fi

if [[ -z "$GITHUB_ACTOR" ]]; then
  GITHUB_ACTOR=$(whoami)
  echo "Setting up GITHUB_ACTOR to current user: ${GITHUB_ACTOR}"
fi

GITVERSION=$(git describe --tags --always)
GITBRANCH=$(git branch | grep \* | cut -d ' ' -f2)
GITREVISION=$(git log -1 --oneline | cut -d ' ' -f1)
TIME=$(date +%FT%T%z)
LDFLAGS="-s -X github.com/prometheus/common/version.Version=${GITVERSION} \
-X github.com/prometheus/common/version.Revision=${GITREVISION} \
-X github.com/prometheus/common/version.Branch=master \
-X github.com/prometheus/common/version.BuildUser=${GITHUB_ACTOR} \
-X github.com/prometheus/common/version.BuildDate=${TIME}"

for OS in "darwin" "linux" "windows"; do
    for ARCH in "amd64" "386"; do 
        echo "Building ${OS}/${ARCH} with version: ${GITVERSION}, revision: ${GITREVISION}, buildUser: ${GITHUB_ACTOR}"
        if [[ $OS == "windows" ]]; then
            GO111MODULE=on CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build -ldflags "${LDFLAGS}" -tags 'netgo static_build' -a -o ".build/${OS}-${ARCH}/beat-exporter.exe"
        else 
            GO111MODULE=on CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build -ldflags "${LDFLAGS}" -tags 'netgo static_build' -a -o ".build/${OS}-${ARCH}/beat-exporter"
        fi
    done
done
