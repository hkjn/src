FROM hkjn/docker:18.06.3-ce
RUN apk add --no-cache python3 && \
    pip3 install docker-squash
ENTRYPOINT ["/usr/bin/docker-squash"]
CMD [""]

