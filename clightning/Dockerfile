# TODO: switch to hkjn/alpine, or at least hkjn/arch
# TODO: can make this Dockerfile.build and add a runtime image as well
# TODO: seems that bitcoin-cli and other binaries is also required here.. compile bitcoin as well?
FROM debian

RUN apt-get -y update && \
    apt-get -y install autoconf build-essential git libtool libgmp-dev \
               libsqlite3-dev python python3 net-tools

RUN mkdir -p src/github.com/elementsproject && \
    cd src/github.com/elementsproject && \
    git clone https://github.com/elementsproject/lightning && \
    cd lightning && \
    make && \
    make install

ENTRYPOINT lightningd
CMD [""]
