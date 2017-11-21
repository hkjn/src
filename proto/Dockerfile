FROM hkjn/arch

WORKDIR /usr/local/src

RUN pacman --noconfirm -Syyy && \
    pacman --noconfirm -S autoconf automake libtool make gcc git unzip
RUN git clone https://github.com/google/protobuf
RUN cd protobuf/ && \
    ./autogen.sh && \
   ./configure && \
   make && \
   make check && \
   make install

ENTRYPOINT ["bash"]
CMD [""]
