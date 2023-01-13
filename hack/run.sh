#!/bin/bash

podman run -it -v installer-debug:latest record openshift-install create cluster --log-level debug --dir
