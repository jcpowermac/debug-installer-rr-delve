#!/bin/bash

. ./hack/config.sh

sudo sysctl kernel.perf_event_paranoid=-1
sudo sysctl kernel.kptr_restrict=0


if ! podman volume exists ${VOLUME_NAME}; then
podman volume create ${VOLUME_NAME}
fi

podman run --rm \
    --name record-installer \
    --cap-add=SYS_PTRACE \
    --cap-add=PERFMON \
    --privileged \
    --security-opt seccomp=unconfined \
    --interactive \
    --tty \
    --env _RR_TRACE_DIR=/trace \
    --volume ${VOLUME_NAME}:/trace:Z \
    --volume ./installer-dir/:/installer:Z \
    installer-debug:${BRANCH} \
    record \
    --disable-avx-512 \
    openshift-install create cluster --log-level debug --dir /installer

