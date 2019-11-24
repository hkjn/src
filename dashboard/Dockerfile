FROM hkjn/golang:1.12.6

WORKDIR /home/go/src/hkjn.me/src/dashboard

COPY ["*.go", "./"]
COPY ["cmd/", "./cmd/"]
COPY ["tmpl/", "./tmpl/"]

RUN go test && \
    go vet && \
    go install hkjn.me/src/dashboard/cmd/gomon

WORKDIR /home/go/bin
COPY *.yaml ./
COPY tmpl/ ./tmpl/

ENTRYPOINT ["gomon"]
CMD ["-alsologtostderr"]
