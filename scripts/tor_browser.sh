#!/bin/bash
#
# Runs Tor browser in a Docker container.
#
sudo docker run -it --rm \
  --memory 3G \
  --name tor-browser \
  -v /tmp/.X11-unix:/tmp/.X11-unix:ro \
  -e DISPLAY=unix${DISPLAY} \
  hkjn/tor-browser:9.0a2
