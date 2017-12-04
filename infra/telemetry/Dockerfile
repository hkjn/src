FROM hkjn/workspace-alpine:1.5.0-arm

ARG version
ENV VERSION=${version} \
    CGO_ENABLED=0

WORKDIR /home/go/src/hkjn.me/hkjninfra/telemetry/
COPY ["server", "./server"]
COPY ["client", "./client"]
COPY ["report", "./report"]
COPY ["vendor", "./vendor"]

USER root
RUN chown -R go:go /home/go/src/hkjn.me
USER go
RUN mkdir -p /home/go/bin/report && \
    protoc \
      -I report/ \
      -I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    report/report.proto --go_out=plugins=grpc:report && \
    cp report/*.pb.go /home/go/bin/report/

WORKDIR /home/go/bin
COPY ["gather_facts", "."]

RUN go test hkjn.me/hkjninfra/telemetry/... && \
    go vet hkjn.me/hkjninfra/telemetry/...

RUN GOARCH=amd64 go build -ldflags "-X main.Version=${VERSION}" -o tclient_x86_64 hkjn.me/hkjninfra/telemetry/client
RUN GOARCH=amd64 go build -ldflags "-X main.Version=${VERSION}" -o tserver_x86_64 hkjn.me/hkjninfra/telemetry/server

RUN GOARCH=arm go build -ldflags "-X main.Version=${VERSION}" -o tclient_armv7l hkjn.me/hkjninfra/telemetry/client
RUN GOARCH=arm go build -ldflags "-X main.Version=${VERSION}"  -o tserver_armv7l hkjn.me/hkjninfra/telemetry/server

RUN sha512sum tclient_* tserver_* gather_facts > SHA512SUMS

ENTRYPOINT ["sh"]


