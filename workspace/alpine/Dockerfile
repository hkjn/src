# The workspace image is a environment for technical work, based on
# Alpine Linux and Golang.
FROM hkjn/golang

ENV LANG=en_US.UTF-8 \
    LANGUAGE=en_US.UTF-8 \
    PROTOBUF_VERSION=v3.3.2 \
    CGO_ENABLED=0 \
    GOOS=linux

USER root
RUN apk add --no-cache bash bind-tools curl coreutils file git python3 mosh openssh tmux vim
RUN echo "Welcome to your workspace." > /etc/motd
RUN chown -R go:go /home/go

USER go

# TODO: Could alternatively take $(id -u) and $(id -g) at run time,
# then create user in entrypoint script which then hands off to bash..
WORKDIR /home/go

COPY [".bash_profile", "./"]
RUN echo "source ~/.bash_profile" > .bashrc

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
CMD [""]
