# The workspace image is an environment for technical work, based on
# Arch Linux.
FROM hkjn/arch

ENV UNPRIVILEGED_UID=1001 \
    LANG=en_US.UTF-8 \
    LANGUAGE=en_US.UTF-8 \
    LC_ALL=C \
    GOPATH=/home/user \
    USER=user \
    PROTOBUF_VERSION=v3.3.2 \
    CGO_ENABLED=0 \
    GOOS=linux

RUN pacman --noconfirm -Syy && \
    pacman --noconfirm -S autoconf automake dnsutils emacs icu gcc git go make python openssh tmux unzip && \
    useradd -m user -u $UNPRIVILEGED_UID -s /bin/bash

# note: seems like protobuf expects older libicuuc.so version, so hack around it.
RUN ln -s /usr/lib/libicuuc.so.60 /usr/lib/libicuuc.so.59

RUN mkdir -p /usr/local/src/github.com/google && \
    cd /usr/local/src/github.com/google && \
    git clone https://github.com/google/protobuf && \
    cd protobuf && \
    git checkout ${PROTOBUF_VERSION} && \
    ./autogen.sh && \
    ./configure && \
    make && \
    make check && \
    make install && \
    rm -rf /usr/local/src/github.com/google/protobuf && \
    pacman --noconfirm -R autoconf automake gcc && \
    pacman -Sc

# TODO: Could alternatively take $(id -u) and $(id -g) at run time,
# then create user in entrypoint script which then hands off to bash..
WORKDIR /home/user

COPY [".bash_profile", "./"]
COPY [".emacs", "./"]
COPY [".emacs.d/", "./.emacs.d/"]

RUN echo "source ~/.bash_profile" > .bashrc && \
    chown -R user:user /home/user /usr/lib/go && \
    echo 'en_US.UTF-8 UTF-8' > /etc/locale.gen && \
    locale-gen

USER user

RUN go get github.com/golang/lint/golint \
           golang.org/x/tools/cmd/goimports \
           github.com/golang/dep/cmd/dep \
           github.com/jteeuwen/go-bindata/... \
           google.golang.org/grpc \
           github.com/golang/protobuf/protoc-gen-go \
           github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
           github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger && \
    rm -rf src/golang.org/x/{net,text,tools}/.git src/github.com/golang/dep/.git

USER user

ENTRYPOINT ["bash"]
