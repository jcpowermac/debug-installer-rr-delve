FROM quay.io/fedora/fedora:latest AS builder

ARG BRANCH=master
ARG REPO=https://github.com/openshift/installer
ARG DELVE_BRANCH=master
ARG DELVE_REPO=https://github.com/go-delve/delve
ENV WORKDIR_PATH="/go/src/github.com/openshift/installer/"
ENV DELVE_WORKDIR_PATH="/go/src/github.com/go-delve/delve/"

RUN mkdir -p $WORKDIR_PATH && \
    mkdir -p $DELVE_WORKDIR_PATH

WORKDIR $WORKDIR_PATH
COPY debug.patch ..
RUN dnf install golang git zip -y
RUN git clone --single-branch --branch $BRANCH $REPO . && \
    git apply ../debug.patch && \
    hack/build.sh

WORKDIR $DELVE_WORKDIR_PATH
COPY delve.patch ..
RUN git clone --single-branch --branch $DELVE_BRANCH $DELVE_REPO . && \
    git apply ../delve.patch && \
    make build

FROM quay.io/fedora/fedora:latest AS run
COPY --from=builder /go/src/github.com/openshift/installer/bin/openshift-install /bin/openshift-install
COPY --from=builder /go/src/github.com/go-delve/delve/dlv /bin/dlv

RUN dnf install golang rr -y

ENTRYPOINT ["/usr/bin/rr"]
CMD ["--help"]
