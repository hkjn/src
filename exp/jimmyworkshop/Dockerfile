FROM hkjn/alpine

RUN apk add --no-cache python3 py-pip gcc git python3-dev musl-dev libzmq zeromq-dev && \
    pip install virtualenv && \
    git clone https://github.com/jimmysong/pb-exercises && \
    cd pb-exercises/ && \
    virtualenv -p python3 .venv && \
    source .venv/bin/activate && \
    pip install -r requirements.txt

ENTRYPOINT jupyter notebook --allow-root
