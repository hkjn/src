FROM hkjn/golang:tip

COPY install.sh /usr/local/bin/
USER root
RUN chown -R go:go /home/go/src
USER go
RUN install.sh

CMD echo "Binary available in ${GOPATH}/bin: $(ls -hsal ${GOPATH}/bin)"
