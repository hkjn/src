#
# Image for Bitcoin Core.
#
# The bin/ and lib/ directories were produced with Dockerfile.build.
#
FROM hkjn/arch

ARG bitcoin_version
ARG lightningd_version

ENV BITCOIN_VERSION=bitcoin_version \
    LIGHTNINGD_VERSION=lightningd_version

COPY bin/ /usr/local/bin/
COPY lib/ /usr/lib/

RUN useradd -m bitcoin
USER root
RUN pacman --noconfirm -Syyy && \
    pacman --noconfirm -S python
USER bitcoin

COPY ["entrypoint.sh", "/usr/local/bin/"]

WORKDIR /home/bitcoin
ENTRYPOINT ["entrypoint.sh"]
CMD [""]

