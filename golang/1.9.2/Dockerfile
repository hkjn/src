FROM hkjn/alpine

ENV GOLANG_VERSION=1.9.2
ENV GOLANG_SRC_URL=https://golang.org/dl/go$GOLANG_VERSION.src.tar.gz \
    GOLANG_SRC_SHA256=665f184bf8ac89986cfd5a4460736976f60b57df6b320ad71ad4cef53bb143dc \
    CGO_ENABLED=0 \
    GOOS=linux

RUN set -eux && \
    apk add --no-cache ca-certificates git libcap musl-dev && \
    apk add --no-cache --virtual .build-deps bash gcc openssl go && \
    export GOROOT_BOOTSTRAP="$(go env GOROOT)" && \
    wget -q "$GOLANG_SRC_URL" -O golang.tar.gz && \
    echo "$GOLANG_SRC_SHA256  golang.tar.gz" | sha256sum -c - && \
    tar -C /usr/local -xzf golang.tar.gz && \
    rm golang.tar.gz && \
    cd /usr/local/go/src && \
    ./make.bash && \
    apk del .build-deps && \
    adduser -D go -u 500 -g 500 -s /bin/bash && \
    mkdir -p /home/go/src && \
    chown -R go:go /usr/local/go/

ENV GOPATH /home/go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN go get -v github.com/golang/dep/cmd/dep

WORKDIR $GOPATH
USER go
