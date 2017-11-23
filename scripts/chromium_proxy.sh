#!/bin/bash

# Start Chromium configured to route all HTTP(S) traffic through a
# SOCKS proxy.

PROXY=193.13.65.254
PORT=5580

chromium --proxy-server="socks5://$PROXY:$PORT" --host-resolver-rules="MAP * 0.0.0.0 , EXCLUDE $PROXY"
