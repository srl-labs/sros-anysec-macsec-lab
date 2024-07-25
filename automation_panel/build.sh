#!/bin/bash

set -e

# build script for frontend and backend containers

BASEDIR=$(realpath $(dirname "$0"))
PROXY_BUILD_ARGS="--build-arg http_proxy=${http_proxy} --build-arg HTTP_PROXY=${HTTP_PROXY} --build-arg https_proxy=${https_proxy} --build-arg HTTPS_PROXY=${HTTPS_PROXY}"

# build the frontend container
docker build ${PROXY_BUILD_ARGS} -t ghcr.io/srl-labs/sros-anysec-macsec-lab/frontend ${BASEDIR}/frontend

# build the backend container
docker build ${PROXY_BUILD_ARGS} -t ghcr.io/srl-labs/sros-anysec-macsec-lab/backend ${BASEDIR}/backend
