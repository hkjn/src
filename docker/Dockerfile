FROM hkjn/alpine:0.2.0

ARG docker_arch
ARG docker_version
ARG docker_sha256

ENV DOCKER_CHANNEL=stable \
    DOCKER_VERSION=${docker_version} \
    DOCKER_ARCH=${docker_arch} \
    DOCKER_SHA256=${docker_sha256}
RUN apk add --no-cache ca-certificates && \
    apk add --no-cache --virtual .fetch-deps curl tar && \
    curl -fsL -o docker.tgz "https://download.docker.com/linux/static/${DOCKER_CHANNEL}/${DOCKER_ARCH}/docker-${DOCKER_VERSION}.tgz"
RUN echo "installing docker-${DOCKER_VERSION} for ${DOCKER_ARCH} (expecting ${DOCKER_SHA256}):" && \
    sha256sum docker.tgz && \
    echo "${DOCKER_SHA256}  docker.tgz" | sha256sum -c - && \
    tar --extract --file docker.tgz --strip-components 1 --directory /usr/local/bin/ && \
    rm docker.tgz && \
    apk del .fetch-deps && \
    dockerd -v && \
    docker -v

COPY docker-entrypoint.sh /usr/local/bin/

ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["sh"]

