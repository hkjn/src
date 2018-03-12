#
# Build image for probes.
#
# As this is a package, there's nothing to release, but this Dockerfile allows
# us to verify that the package continues building. We don't lock in a specific
# version for vendoring here of our own third party deps like github.com/golang/glog,
# as that can be done in the downstream tools importing our package if appropriate.
#
FROM hkjn/golang

WORKDIR /home/go/src/hkjn.me/src/probes

COPY *.go ./
COPY prober ../prober
COPY dnsprobe ./dnsprobe
COPY webprobe ./webprobe
USER root
RUN chown -R go:go /home/go/
USER go
RUN go get -v ./... && \
    go vet && \
    go install
WORKDIR /build
