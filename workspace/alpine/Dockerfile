# The workspace image is a environment for technical work, based on
# Alpine Linux and Golang.
FROM hkjn/golang

ENV LANG=en_US.UTF-8 \
    LANGUAGE=en_US.UTF-8 \
    PROTOBUF_VERSION=v3.3.2 \
    CGO_ENABLED=0 \
    GOOS=linux

USER root
RUN apk add --no-cache bash bind-tools curl coreutils file git musl-dev python3 mosh openssh tmux vim && \
    apk add --no-cache --virtual .build-deps autoconf automake gcc g++ make libtool libc6-compat
RUN mkdir -p src/github.com/google && \
    cd src/github.com/google && \
    git clone https://github.com/google/protobuf && \
    cd protobuf && \
    git checkout ${PROTOBUF_VERSION} && \
    ./autogen.sh && \
    ./configure
RUN cd src/github.com/google/protobuf && \
    make && \
    make check && \
    make install && \
    ln -s /usr/bin/python3 /usr/bin/python && \
    ln -s /usr/bin/pip3 /usr/bin/pip && \
    rm -rf /home/go/src/github.com/google/protobuf && \
    apk del .build-deps

# TODO: Could alternatively take $(id -u) and $(id -g) at run time,
# then create user in entrypoint script which then hands off to bash..
WORKDIR /home/go

COPY [".bash_profile", "./"]
RUN echo "source ~/.bash_profile" > .bashrc && \
    chown -R go:go /home/go && \
    echo "Welcome to your workspace." > /etc/motd

USER go

RUN go get github.com/golang/lint/golint \
           golang.org/x/tools/cmd/goimports \
           github.com/golang/dep/cmd/dep \
           github.com/jteeuwen/go-bindata/... \
           google.golang.org/grpc \
           github.com/golang/protobuf/protoc-gen-go \
           github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
           github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger && \
     rm -rf src/golang.org/x/{net,text,tools}/.git src/github.com/golang/dep/.git

ENTRYPOINT ["bash"]
