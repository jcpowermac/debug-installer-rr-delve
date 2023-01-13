FROM quay.io/fedora/fedora:latest AS builder

ARG BRANCH=master
ARG REPO=https://github.com/openshift/installer
ENV WORKDIR_PATH="/go/src/github.com/openshift/installer/"

RUN mkdir -p $WORKDIR_PATH
WORKDIR $WORKDIR_PATH
COPY debug.patch ..
RUN dnf install golang git zip -y && \
    git clone --single-branch --branch $BRANCH $REPO . && \
    git apply ../debug.patch && \
    hack/build.sh

FROM quay.io/fedora/fedora:latest AS run
COPY --from=builder /go/src/github.com/openshift/installer/bin/openshift-install /bin/openshift-install

RUN dnf install golang rr -y && \
    go install github.com/go-delve/delve/cmd/dlv@latest

ENTRYPOINT ["/usr/bin/rr"]
CMD ["--help"]
