FROM hkjn/alpine

WORKDIR /usr/local/src


RUN apk add --no-cache bash python3 py-pip gcc git python3-dev musl-dev libzmq zeromq-dev && \
    pip install virtualenv && \
    git clone https://github.com/jimmysong/pb-exercises && \
    cd pb-exercises/ && \
    virtualenv -p python3 .venv && \
    source .venv/bin/activate && \
    pip install -r requirements.txt && \
    adduser -D user -s /bin/bash

WORKDIR /usr/local/src/pb-exercises

USER user

ENTRYPOINT source .venv/bin/activate && jupyter notebook --port=8888 --no-browser --ip=0.0.0.0
