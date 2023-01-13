#!/bin/bash

. ./hack/config.sh

podman volume rm ${VOLUME_NAME}
