FROM hkjn/workspace-alpine:1.6.1

COPY [".gitconfig", "./"]
COPY ["known_hosts", "./.ssh/"]
USER root
RUN chown -R $USER_UID:$USER_GID .ssh/
USER user
ENTRYPOINT ["tmux"]
