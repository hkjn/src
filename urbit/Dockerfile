#
# Usage:
#   make init
#      > |mount %
#
FROM debian:stretch

RUN apt update && \
    apt -y install \
      clang \
      libncurses-dev \
      libcurl4-gnutls-dev \
      libssl-dev \
      libgmp3-dev \
      libsigsegv-dev \
      gdb \
      git \
      ninja-build \
      openssl \
      pkg-config \
      python3-pip \
      unzip \
      zlib1g \
      zlib1g-dev && \
      pip3 install --user -I meson==0.44.1

WORKDIR /urbit

RUN git clone https://github.com/urbit/urbit
RUN cd urbit && \
    bash -c "export PATH=${PATH}:/root/.local/bin; meson . ./build --buildtype=debugoptimized -Dgc=true -Dprof=true" && \
    ./scripts/bootstrap && \
    ./scripts/build && \
    ninja -C ./build/ install

ENTRYPOINT ["/usr/local/bin/urbit"]
CMD ["-F", "zod"]
