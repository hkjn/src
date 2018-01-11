#
# Image for github.com/jnewbery/bitcointools.
#
FROM hkjn/alpine

ARG tools_version
ENV TOOLS_VERSION=${tools_version}

WORKDIR /usr/local/src/
RUN apk add --no-cache git python3 && \
    git clone https://github.com/jnewbery/bitcointools
WORKDIR /usr/local/src/bitcointools/

RUN git checkout ${TOOLS_VERSION} && \
    echo "Using version: $(git show)"
