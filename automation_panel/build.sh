#!/bin/bash

# for multiarch builds check that you have buildx setup:
# docker buildx create --name multiarch --platform linux/amd64,linux/arm64 --use
# docker buildx inspect --bootstrap
# docker buildx ls

set -e

# build script for frontend and backend containers
TAG=$1
if [ -z "$TAG" ]; then

    echo "Usage: bash build.sh <version>"
    echo "Example: bash build.sh v0.5.0"
    exit 1
fi


BASEDIR=$(realpath $(dirname "$0"))
PROXY_BUILD_ARGS="--build-arg http_proxy=${http_proxy} --build-arg HTTP_PROXY=${HTTP_PROXY} --build-arg https_proxy=${https_proxy} --build-arg HTTPS_PROXY=${HTTPS_PROXY}"

# build the automation panel container
docker buildx build --push --platform linux/amd64,linux/arm64 ${PROXY_BUILD_ARGS} -t ghcr.io/srl-labs/sros-anysec-macsec-lab/panel:${TAG} ${BASEDIR}

# for local testing without pushing to registry, use this instead of the previous command:
# docker build ${PROXY_BUILD_ARGS} -t -t ghcr.io/srl-labs/sros-anysec-macsec-lab/panel:${TAG} ${BASEDIR}
