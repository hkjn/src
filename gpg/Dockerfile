FROM hkjn/alpine

# GPG needs the GPG_TTY set to not get confused:
# - https://unix.stackexchange.com/a/296496
ENV GPG_TTY=/dev/console
RUN apk add --no-cache bash gnupg vim && \
    adduser -D gpg -s /bin/bash

USER gpg
WORKDIR /home/gpg
ENTRYPOINT ["bash"]
