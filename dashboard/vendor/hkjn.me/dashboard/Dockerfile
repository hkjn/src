FROM hkjn/golang

WORKDIR /home/go/src/hkjn.me/dashboard

COPY ["*.go", "./"]
COPY ["cmd/", "./cmd/"]
COPY ["tmpl/", "./tmpl/"]
COPY ["vendor", "./vendor/"]

RUN go test && \
    go vet && \
    go install hkjn.me/dashboard/cmd/gomon

WORKDIR /home/go/bin
COPY *.yaml ./
COPY tmpl/ ./tmpl/

ENTRYPOINT ["gomon"]
CMD ["-alsologtostderr"]
