#!/bin/bash

. ./hack/config.sh

podman build \
    --build-arg BRANCH=${BRANCH} \
    --file images/Dockerfile \
    --tag installer-debug:${BRANCH} .
