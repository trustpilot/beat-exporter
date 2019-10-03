#!/bin/bash

set -e

if [[ -z "$GITHUB_WORKSPACE" ]]; then
  echo "Set the GITHUB_WORKSPACE env variable."
  exit 1
fi

if [[ -z "$GITHUB_REPOSITORY" ]]; then
  echo "Set the GITHUB_REPOSITORY env variable."
  exit 1
fi

echo "GITHUB WORKSPACE: ${GITHUB_WORKSPACE}"
echo "GITHUB REPO: ${GITHUB_REPOSITORY}"

RELEASE_DIR=.release
VERSION=$(git describe --tags | cut -d '-' -f1 | cut -d 'v' -f2)
RELEASE_FILES=LICENSE

mkdir -p $RELEASE_DIR

for ARTIFACT in $(ls .build); do
    ARTIFACT_NAME="beat-exporter-${VERSION}-${ARTIFACT}.tar.gz"
    echo "Creating ${ARTIFACT_NAME}"
    for FILE in $RELEASE_FILES; do
      cp "${GITHUB_WORKSPACE}/${FILE}" "${GITHUB_WORKSPACE}/.build/${ARTIFACT}/${FILE}"
    done
    
    cd "${GITHUB_WORKSPACE}/.build/${ARTIFACT}" && tar -cvzf ${GITHUB_WORKSPACE}/${RELEASE_DIR}/${ARTIFACT_NAME} *
done
