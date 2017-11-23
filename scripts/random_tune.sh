#!/bin/bash
set -euo pipefail
# Play some random tunes.

MUSIC_DIR=${MUSIC_DIR:-/home/$USER/media/music/}
[[ -e $MUSIC_DIR ]] || {
  echo "No such MUSIC_DIR='$MUSIC_DIR'" >&2
  exit 1
}
mplayer "$(find $MUSIC_DIR -name '*.mp3' | shuf -n1)"
