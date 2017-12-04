FROM hkjn/workspace-alpine:1.5.0-arm

ENV CGO_ENABLED=0
WORKDIR /home/go/src/hkjn.me/hkjninfra/secretservice/

COPY ["*.go", "./"]
COPY ["cmd", "./cmd/"]
COPY ["vendor", "./vendor/"]

USER go

RUN go test && \
    go build -o /home/go/bin/secretservice ./cmd/
WORKDIR /home/go/bin/

RUN sha512sum secretservice* > SHA512SUMS

CMD echo "Binaries available in $(pwd): $(ls)"
