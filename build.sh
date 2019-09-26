#!/bin/bash

set -e

if [[ -z "$GITHUB_WORKSPACE" ]]; then
  echo "Set the GITHUB_WORKSPACE env variable."
  exit 1
fi

GITVERSION=$(git describe --tags --always --long --dirty)
GITBRANCH=$(git branch | grep \* | cut -d ' ' -f2)
GITREVISION=$(git log -1 --oneline | cut -d ' ' -f1)
TIME=$(date +%FT%T%z)
LDFLAGS="-s -X github.com/prometheus/common/version.Version=${GITVERSION} \
-X github.com/prometheus/common/version.Revision=${GITREVISION} \
-X github.com/prometheus/common/version.Branch=${GITBRANCH} \
-X github.com/prometheus/common/version.BuildUser=root@localhost \
-X github.com/prometheus/common/version.BuildDate=${TIME}"

for OS in "darwin" "linux" "windows"
do
    for ARCH in "amd64" "386"
    do 
        echo "Building ${OS}/${ARCH}"
        if [[ $OS == "windows" ]]; then
            GO111MODULE=on CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build -ldflags "${LDFLAGS}" -tags 'netgo static_build' -a -o ".build/${OS}-${ARCH}/beat-exporter.exe"
        else 
            GO111MODULE=on CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build -ldflags "${LDFLAGS}" -tags 'netgo static_build' -a -o ".build/${OS}-${ARCH}/beat-exporter"
        fi
       
    done
done

