#!/bin/bash

# Play some random tunes.
MUSIC_DIR=${MUSIC_DIR:-/home/$USER/media/music/}
[[ -e $MUSIC_DIR ]] || {
  echo "No such MUSIC_DIR='$MUSIC_DIR'" >&2
  exit 1
}
IFS=$'\n' y=($songs)
SONGS="$(find $MUSIC_DIR -name '*.mp3' | shuf -n10)"
for SONG in ${SONGS[@]}; do
  mplayer "$SONG"
done
