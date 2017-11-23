#!/bin/bash

# Start Chromium configured to route all traffic through a
# Tor's SOCKS proxy.
chromium --proxy-server="socks://localhost:9050" --user-data-dir=/tmp

