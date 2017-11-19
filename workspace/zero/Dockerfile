FROM hkjn/workspace-alpine:1.5.0-arm
# FIXME: drop -arm once multiarch is built

COPY [".gitconfig", "./"]
COPY ["known_hosts", "./.ssh/"]
USER root
RUN chown -R $USER_UID:$USER_GID .ssh/
USER user
ENTRYPOINT ["tmux"]
