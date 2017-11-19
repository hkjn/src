FROM hkjn/workspace-arch:1.5.0-arm
# FIXME: Drop -arm once multiarch is built

COPY [".gitconfig", "./"]
COPY [".bash_profile.extra", "./"]
RUN cat .bash_profile.extra >> ./.bash_profile && \
    rm .bash_profile.extra

RUN go get golang.org/x/crypto/acme/autocert

ENTRYPOINT ["tmux"]
