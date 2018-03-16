#
# Image for building Bitcoin Core binaries.
#
FROM hkjn/arch

ARG bitcoin_version

ENV BITCOIN_VERSION=${bitcoin_version}

# Install dependencies.
RUN pacman -Syyu --noconfirm
RUN pacman -S --noconfirm git base-devel boost libevent libunwind jq python

WORKDIR /usr/local/src/bitcoin/
RUN useradd -m bitcoin && \
    chown -R bitcoin:bitcoin /usr/local/src /usr/local/lib && \
    mkdir /usr/local/src/lightning

USER bitcoin

RUN git clone https://github.com/bitcoin/bitcoin .
RUN git checkout ${BITCOIN_VERSION}


RUN echo "Generating Makefile and configure script.." && \
    ./autogen.sh

RUN echo "Building dependencies.." && \
    cd depends && \
    make

# TODO: support x86.
RUN echo "Running configure script.." && \
    ./configure --prefix=$(pwd)/depends/armv7l-unknown-linux-gnueabihf

RUN echo "Building bitcoin at: $(git show)" && \
    make

RUN echo "Running tests.." && \
    make check
# TODO: enable running .py tests if they can be made to work on 2G machines, or we set up larger
# build machines:
#    ./test/functional/test_runner.py

USER root
RUN echo "Installing bitcoin binaries.." && \
    make install
# TODO: Improve this ugly workaround for ./configure --prefix above making 'make install' put the
# bins and libs in that dir.
RUN mv -iv depends/armv7l-unknown-linux-gnueabihf/bin/* /usr/local/bin/ && \
    mv -iv depends/armv7l-unknown-linux-gnueabihf/lib/pkgconfig /usr/local/lib/ && \
    mv -iv depends/armv7l-unknown-linux-gnueabihf/share/man/man1 /usr/local/man/man1/ && \
    mv -iv depends/armv7l-unknown-linux-gnueabihf/include /usr/local/include/

USER bitcoin
RUN mkdir -p ${HOME}/src/ && \
    cp -v share/rpcauth/rpcauth.py ${HOME}/src/

ARG lightning_version
ARG lightning_remote

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
