#!/usr/bin/env bash
docker run --rm -it -v /containers/eth:/root/.ethereum hkjn/geth:x86_64 --fast --cache 512 console
