#!/bin/bash

set -euo pipefail

start() {
	[ -d $2 ] || {
		echo "No such directory: '$2'"
		return
	}
	docker rm -f $1-fileserver && echo "Removed $1-fileserver container"
	docker run --name $1-fileserver -d -v "$2:/var/www" -p $3:8080 hkjn/fileserver
	echo "$1-fileserver is running at $3, serving directory '$2'"
}
start musashi /media/musashi 8080
start staging $HOME/staging 8081
start media $HOME/media 8082
start tvideo /media/timothy/video 8083
start tcomedic /media/timothy/comedic 8084
start teducational /media/timothy/educational 8085
start tmusic /media/timothy/music 8086
start tanime /media/timothy/anime 8087
#start usb1 /run/media/zero/USB20FD 8085
#start usb2 /run/media/zero/f538fa97-80d3-4ef3-9010-99460637a69a 8086

echo "Started media fileserver containers."
