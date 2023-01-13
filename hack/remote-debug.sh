#!/bin/bash

# need a volume to put rr in.

podman run -it -p 2345:2345 --entrypoint dlv installer-debug \
    replay \
    --disable-aslr \
    --backend=rr \
    --headless \
    --listen=:2345 \
    --api-version=2 \
    --accept-multiclient ${HOME}/.local/share/rr/openshift-install-0

