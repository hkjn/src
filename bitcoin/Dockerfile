#
# Image for Bitcoin Core.
#
# The bin/ and lib/ directories were produced with Dockerfile.build.
#
FROM hkjn/arch

COPY bin/ /usr/local/bin/
COPY lib/ /usr/lib/

RUN useradd -m bitcoin

WORKDIR /usr/local/bin
USER bitcoin
ENTRYPOINT ["bitcoind"]
CMD ["-printtoconsole"]

