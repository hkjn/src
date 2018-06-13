FROM hkjn/workspace-alpine:1.5.0-arm

ARG version
ENV CGO_ENABLED=0

WORKDIR /home/go/src/hkjn.me/src/exp/kraken/
COPY ["*.go", "./"]
COPY ["vendor/", "./vendor/"]

RUN go test ./... && \
    go vet ./...

RUN go build
