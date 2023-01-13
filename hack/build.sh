#!/bin/bash

. ./hack/config.sh

podman build \
    --build-arg BRANCH=${BRANCH} \
    --buld-arg REPO=${REPO} \
    --file images/Dockerfile \
    --tag installer-debug:${BRANCH} .
