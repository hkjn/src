#!/bin/bash
#
# Convert all .mkv and .webm files in directory to .mp4
#
set -euo pipefail
if ls *.mkv 2</dev/null; then
  for v in *.mkv; do
    ffmpeg -i "$v" -vcodec copy -acodec copy -scodec mov_text "${v%.mkv}.mp4" ||
      ffmpeg -i "$v" -vcodec copy -c:a aac -scodec mov_text "${v%.mkv}.mp4" ||
      ffmpeg -i "$v" -c:v libx264 -acodec copy -scodec mov_text "${v%.mkv}.mp4"
    [[ -e "${v%.mkv}.mp4" ]] && rm "$v"
  done
fi
if ls *.webm 2>/dev/null; then
  for  v in *.webm; do
    ffmpeg -i "$v" -vcodec libx264 "${v%.webm}.mp4"
    [[ -e "${v%.webm}.mp4" ]] && rm "$v"
  done
fi
