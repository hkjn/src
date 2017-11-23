#!/bin/bash
WALLPAPER_PATH=/home/zero/wallpapers/
pcmanfm -w "$(find $WALLPAPER_PATH  -type f | shuf -n1)"
