ARG FROM_IMAGE
ARG FROM_TAG
FROM ${FROM_IMAGE}:${FROM_TAG}

# > automatic buildx ARGs
ARG TARGETARCH

# > ARGs before FROM are not accessible
ARG FROM_IMAGE
ARG FROM_TAG

# > Our custom ARGs
ARG NODE_VERSION="12 16"
ARG DISTRO=ubuntu

# > Force apt to not be interactive/not ask
ENV DEBIAN_FRONTEND=noninteractive

SHELL [ "/bin/bash", "--login", "-e", "-o", "pipefail", "-c" ]
WORKDIR /tmp

COPY act.sh /imagegeneration/installers/act.sh
RUN /imagegeneration/installers/act.sh

ARG BUILD_DATE
ARG BUILD_TAG=act
ARG BUILD_REF
ARG BUILD_TAG_VERSION
ARG BUILD_OWNER
ARG BUILD_REPO

USER runner
