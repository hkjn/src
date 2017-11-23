#!/bin/bash
docker run -it \
  --rm \
  --net host \
  --memory 512mb \
  -v /tmp/.X11-unix:/tmp/.X11-unix \
  -e DISPLAY=unix$DISPLAY \
  -v $HOME/Downloads:/root/Downloads \
  -v $HOME/.config/google-chrome/:/data \
  -v /dev/snd:/dev/snd --privileged \
  --name chrome \
  jess/chrome
