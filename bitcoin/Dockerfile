#
# Image for building Bitcoin Core binaries.
#
FROM hkjn/arch

ARG bitcoin_version
ARG bitcoin_arch
ARG lightning_version
ARG lightning_remote

ENV BITCOIN_VERSION=${bitcoin_version}

# Install dependencies.
RUN pacman -Syyu --noconfirm && \
    pacman -S --noconfirm base-devel boost doxygen git libevent libunwind jq python

WORKDIR /usr/local/src/bitcoin/
RUN useradd -m bitcoin && \
    chown -R bitcoin:bitcoin /usr/local/src /usr/local/lib && \
    mkdir /usr/local/src/lightning

USER bitcoin

RUN git clone https://github.com/bitcoin/bitcoin .
RUN echo "Checking of bitcoin version '${BITCOIN_VERSION}'.." && \
    git checkout ${BITCOIN_VERSION}

RUN echo "Generating Makefile and configure script.." && \
    ./autogen.sh

RUN echo "Building dependencies.." && \
    cd depends && \
    make

# TODO: support x86:
#     ./configure --prefix=$(pwd)/depends/x86_64-pc-linux-gnu
RUN echo "Running configure script.." && \
    ./configure --prefix=$(pwd)/depends/${bitcoin_arch} \
                --exec-prefix=/usr/local

RUN echo "Building bitcoin at: $(git show)" && \
    make

RUN echo "Running tests.." && \
    make check || \
    cat src/test-suite.log

# TODO: enable running .py tests if they can be made to work on 2G machines, or we set up larger
# build machines:
#    ./test/functional/test_runner.py

USER root
RUN echo "Installing bitcoin binaries.." && \
    make install

USER bitcoin
RUN mkdir -p ${HOME}/src/ && \
    cp -v share/rpcauth/rpcauth.py ${HOME}/src/

ENV LIGHTNING_VERSION=${lightning_version} \
    LIGHTNING_REMOTE=${lightning_remote}

WORKDIR /usr/local/src/lightning
USER root
RUN chown -R bitcoin:bitcoin /usr/local/src/lightning
USER bitcoin
RUN git clone ${LIGHTNING_REMOTE} . && \
    git checkout ${LIGHTNING_VERSION} && \
    echo "Building lightning at $(git show)"

RUN make

USER root
RUN make install
RUN rm -rf /usr/local/src && \
    pacman --noconfirm -Sc && \
    rm -rf /var/cache/pacman/pkg/*
USER bitcoin

WORKDIR /usr/local/src/bitcoin

COPY entrypoint.sh /usr/local/bin
WORKDIR /home/bitcoin
ENTRYPOINT ["entrypoint.sh"]
