FROM hkjn/alpine:0.2.0

ENV GOLANG_VERSION=1.18.1
ENV GOLANG_SRC_URL=https://golang.org/dl/go$GOLANG_VERSION.src.tar.gz \
    GOLANG_SRC_SHA256=efd43e0f1402e083b73a03d444b7b6576bb4c539ac46208b63a916b69aca4088 \
    CGO_ENABLED=0 \
    GOOS=linux

RUN set -eux && \
    apk add --no-cache ca-certificates bash git libcap make musl-dev && \
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

RUN go install -v github.com/jteeuwen/go-bindata/go-bindata@latest

