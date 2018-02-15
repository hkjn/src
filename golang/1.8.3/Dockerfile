FROM hkjn/alpine

ENV GOLANG_VERSION=1.8.3
ENV GOLANG_SRC_URL=https://golang.org/dl/go$GOLANG_VERSION.src.tar.gz \
    GOLANG_SRC_SHA256=5f5dea2447e7dcfdc50fa6b94c512e58bfba5673c039259fd843f68829d99fa6 \
    CGO_ENABLED=0 \
    GOOS=linux

# https://golang.org/issue/14851
COPY no-pic.patch /

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
     patch -p2 -i /no-pic.patch && \
    rm -rf /*.patch && \
    apk del .build-deps && \
    adduser -D go -u 500 -g 500 -s /bin/bash && \
    mkdir -p /home/go/src && \
    chown -R go:go /usr/local/go/ /home/go

ENV GOPATH /home/go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN go get github.com/golang/dep/cmd/dep
RUN go get -v github.com/jteeuwen/go-bindata/...

WORKDIR $GOPATH
USER go
