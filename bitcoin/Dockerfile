#
# Image for Bitcoin Core.
#
# The bin/ and lib/ directories were produced with Dockerfile.build.
#
FROM hkjn/arch

COPY bin/ /usr/local/bin/
COPY lib/ /usr/lib/
COPY libexec/ /usr/local/libexec/c-lightning/

RUN useradd -m bitcoin
USER root
RUN pacman --noconfirm -Syyy && \
    pacman --noconfirm -S python
USER bitcoin

COPY ["entrypoint.sh", "/usr/local/bin/"]

WORKDIR /home/bitcoin
ENTRYPOINT ["entrypoint.sh"]
CMD [""]

