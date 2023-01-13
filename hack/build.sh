#!/bin/bash

. ./hack/config.sh

podman build \
    --build-arg BRANCH=${BRANCH} \
    --build-arg REPO=${REPO} \
    --file images/Dockerfile \
    --tag installer-debug:${BRANCH} .
