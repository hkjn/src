FROM hkjn/alpine

# via https://cloud.google.com/sdk/docs/quickstart-linux
ENV GCLOUD_VERSION=183.0.0-linux-x86_64

RUN apk add --no-cache curl python bash && \
    adduser -D gcloud -s /bin/bash

WORKDIR /home/gcloud
USER gcloud

COPY [".bash*", "./"]

RUN curl -fSLO https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-${GCLOUD_VERSION}.tar.gz && \
    tar xzfv google-cloud-sdk-${GCLOUD_VERSION}.tar.gz && \
    rm google-cloud-sdk-${GCLOUD_VERSION}.tar.gz
USER root
RUN chown -R gcloud:gcloud /home/gcloud
USER gcloud
RUN ./google-cloud-sdk/install.sh

ENTRYPOINT ["bash"]
