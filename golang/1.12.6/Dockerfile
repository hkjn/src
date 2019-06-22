FROM hkjn/alpine:0.2.0

ENV GOLANG_VERSION=1.12.6
ENV GOLANG_SRC_URL=https://golang.org/dl/go$GOLANG_VERSION.src.tar.gz \
    GOLANG_SRC_SHA256=c96c5ccc7455638ae1a8b7498a030fe653731c8391c5f8e79590bce72f92b4ca \
    CGO_ENABLED=0 \
    GOOS=linux

RUN set -eux && \
    apk add --no-cache ca-certificates git libcap musl-dev && \
    apk add --no-cache --virtual .build-deps bash gcc openssl go && \
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

RUN go get -v github.com/golang/dep/cmd/dep
RUN go get -v github.com/jteeuwen/go-bindata/...

WORKDIR ${GOPATH}
USER go
