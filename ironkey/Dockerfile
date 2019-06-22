FROM debian:jessie-slim

RUN apt-get -y update && \
    apt-get -y install gcc-multilib

COPY lock unlock /usr/local/bin/

CMD ["lock"]
