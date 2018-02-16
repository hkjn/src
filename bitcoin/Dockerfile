#
# Image for building Bitcoin Core binaries.
#
FROM hkjn/arch

ARG bitcoin_version
ARG lightning_version

ENV BITCOIN_VERSION=${bitcoin_version} \
    LIGHTNING_VERSION=${lightning_version}

# Install dependencies.
RUN pacman-key --init
RUN pacman-key --populate archlinux
RUN pacman-key --refresh-keys
RUN pacman -Syyu --noconfirm
RUN pacman -S --noconfirm git base-devel boost libevent jq python


WORKDIR /usr/local/src/bitcoin/
RUN useradd -m bitcoin && \
    chown -R bitcoin:bitcoin /usr/local/src /usr/local/lib && \
    mkdir /usr/local/src/lightning

USER bitcoin

RUN git clone https://github.com/bitcoin/bitcoin .
RUN git checkout ${BITCOIN_VERSION}

RUN ./autogen.sh && \
    ./configure --disable-wallet

RUN echo "Building bitcoin at: $(git show)" && \
    make
USER root
RUN make install
USER bitcoin

WORKDIR /usr/local/src/lightning
USER root
RUN chown -R bitcoin:bitcoin /usr/local/src/lightning
USER bitcoin
RUN git clone https://github.com/elementsproject/lightning . && \
    git checkout ${LIGHTNING_VERSION} && \
    echo "Building lightning at $(git show)"

RUN make

USER root
RUN make install
RUN rm -rf /usr/local/src && \
    pacman -Sc
USER bitcoin

WORKDIR /usr/local/src/bitcoin

COPY sqlversiontest.cpp /usr/local/src/bitcoin
RUN g++ -o sqlversiontest sqlversiontest.cpp -lsqlite3 && \
    ./sqlversiontest

COPY entrypoint.sh /usr/local/bin
WORKDIR /home/bitcoin
ENTRYPOINT ["entrypoint.sh"]
