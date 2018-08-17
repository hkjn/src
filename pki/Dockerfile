FROM hkjn/alpine

ENV GOPATH /go
ENV PATH $PATH:$GOPATH/bin
WORKDIR $GOPATH
RUN apk add --no-cache go git gcc musl-dev && \
    go get -tags nopkcs11 github.com/cloudflare/cfssl/cmd/cfssl && \
    go get github.com/cloudflare/cfssl/cmd/cfssljson

COPY configs/ /etc/configs
COPY initca /usr/local/bin/
COPY gencert /usr/local/bin/

WORKDIR /certs

CMD ["sh"]

