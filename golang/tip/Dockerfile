FROM hkjn/alpine

RUN set -eux && \
    apk add --no-cache ca-certificates git libcap musl-dev && \
    apk add --no-cache --virtual .build-deps bash gcc openssl go && \
    export GOROOT_BOOTSTRAP="$(go env GOROOT)" && \
    cd /usr/local/ && \
    git clone https://github.com/golang/go.git && \
    cd go/src && \
    echo "Building from commit: $(git show)" && \
    ./make.bash && \
    apk del .build-deps && \
    adduser -D go -u 500 -g 500 -s /bin/bash && \
    mkdir -p /home/go/src && \
    chown -R go:go /usr/local/go/

ENV GOPATH /home/go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

WORKDIR $GOPATH
USER go
