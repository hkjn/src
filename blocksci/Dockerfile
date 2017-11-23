FROM ubuntu

RUN apt-get -y update && \
    apt-get -y install build-essential cmake libfreetype6-dev libssl-dev libboost-all-dev libpng12-dev \
               libsqlite3-dev autogen autoconf git libleveldb-dev libcurl4-openssl-dev libjsoncpp-dev \
               libjsonrpccpp-dev libjsonrpccpp-tools python3-dev python3-pip software-properties-common \
               pkg-config python-software-properties && \
    add-apt-repository ppa:ubuntu-toolchain-r/test -y && \
    apt-get -y update && \
    apt-get -y install gcc-6 g++-6 && \
    update-alternatives --install /usr/bin/gcc gcc /usr/bin/gcc-6 60 --slave /usr/bin/g++ g++ /usr/bin/g++-6

RUN pip3 install --upgrade pip && \
    pip3 install --upgrade multiprocess psutil jupyter pycrypto matplotlib pandas dateparser

WORKDIR /usr/local/src

RUN git clone https://github.com/bitcoin-core/secp256k1 && \
    cd secp256k1 && \
    ./autogen.sh && \
    ./configure --enable-module-recovery && \
    make && \
    make install

RUN git clone https://github.com/citp/BlockSci.git && \
    cd BlockSci && \
    git submodule update --init --recursive --checkout && \
    cd libs/bitcoin-api-cpp && \
    mkdir release && \
    cd release && \
    cmake --DCMAKE_BUILD_TYPE=Release .. && \
    make && \
    make install && \
    cd ../../.. && \
    mkdir release && \
    cd release && \
    cmake --DCMAKE_BUILD_TYPE=Release .. && \
    make && \
    make install

ENTRYPOINT ["bash"]
