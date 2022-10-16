FROM hkjn/alpine:0.2.0

ENV GOLANG_VERSION=1.17.13
ENV GOLANG_SRC_URL=https://golang.org/dl/go$GOLANG_VERSION.src.tar.gz \
    GOLANG_SRC_SHA256=a1a48b23afb206f95e7bbaa9b898d965f90826f6f1d1fc0c1d784ada0cd300fd \
    CGO_ENABLED=0 \
    GOOS=linux

RUN set -eux && \
    apk add --no-cache ca-certificates bash curl git libcap make musl-dev && \
    apk add --no-cache --virtual .build-deps gcc openssl go && \
    export GOROOT_BOOTSTRAP="$(go env GOROOT)" && \
    wget -q "${GOLANG_SRC_URL}" -O golang.tar.gz && \
    echo "${GOLANG_SRC_SHA256}  golang.tar.gz" | sha256sum -c - && \
    tar -C /usr/local -xzf golang.tar.gz && \
    rm golang.tar.gz && \
    cd /usr/local/go/src && \
    ./make.bash && \
    apk del .build-deps && \
    adduser -D go -u 511 -g 511 -s /bin/bash && \
    mkdir -p /home/go/src && \
    chown -R go:go /usr/local/go/ /home/go

ENV GOPATH /home/go
ENV PATH ${GOPATH}/bin:/usr/local/go/bin:${PATH}

USER go
