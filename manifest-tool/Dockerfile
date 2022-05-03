FROM hkjn/golang:1.18.1

USER root
RUN apk add --no-cache gcc make && \
    cd $GOPATH/src && \
    mkdir -p github.com/estesp && \
    cd github.com/estesp && \
    git clone https://github.com/estesp/manifest-tool && \
    cd manifest-tool && \
    make binary && \
    install manifest-tool /usr/local/bin && \
    apk del gcc make

ENTRYPOINT ["manifest-tool"]
