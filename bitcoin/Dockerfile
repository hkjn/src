#
# Image for Bitcoin Core.
#
# The bin/ and lib/ directories were produced with Dockerfile.build.
#
FROM hkjn/arch

COPY bin/ /usr/local/bin/
COPY lib/ /usr/lib/

RUN useradd -m bitcoin
USER root
RUN pacman --noconfirm -Syyy && \
    pacman --noconfirm -S python
USER bitcoin

WORKDIR /home/bitcoin
ENTRYPOINT ["bitcoind"]
CMD ["-printtoconsole"]

