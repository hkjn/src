FROM hkjn/alpine

# GPG needs the GPG_TTY set to not get confused:
# - https://unix.stackexchange.com/a/296496
ENV GPG_TTY=/dev/console
RUN apk add --no-cache bash gnupg grep musl-dev nano secure-delete vim && \
    adduser -D gpg -s /bin/bash
COPY keys/ /etc/keys
WORKDIR /usr/local/bin
COPY bin/ ./
USER gpg
WORKDIR /home/gpg
ENTRYPOINT ["entrypoint.sh"]
CMD ["bash"]
