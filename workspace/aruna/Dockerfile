FROM hkjn/workspace-arch:1.6.0

COPY [".gitconfig", "./"]
COPY [".bash_profile.extra", "./"]
RUN cat .bash_profile.extra >> ./.bash_profile && \
    rm .bash_profile.extra

RUN go get golang.org/x/crypto/acme/autocert

ENTRYPOINT ["tmux"]
