#!/bin/bash

. ./hack/config.sh

if ! podman volume exists ${VOLUME_NAME}; then
podman volume create ${VOLUME_NAME}
fi

podman run --rm \
    --name replay-installer \
    --publish 0.0.0.0:2345:2345 \
    --cap-add=SYS_PTRACE \
    --cap-add=PERFMON \
    --privileged \
    --security-opt seccomp=unconfined \
    --interactive \
    --tty \
    --env _RR_TRACE_DIR=/trace \
    --volume ${VOLUME_NAME}:/trace:Z \
    --entrypoint /root/go/bin/dlv \
    installer-debug:${BRANCH} \
    replay \
    --disable-aslr \
    --backend=rr \
    --headless \
    --listen=:2345 \
    --api-version=2 \
    --accept-multiclient /trace/openshift-install-0

# Not extensively tested, add back if needed
#    --cap-add=SYS_ADMIN \
