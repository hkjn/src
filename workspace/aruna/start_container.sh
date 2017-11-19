#!/bin/bash -eu
#
# Start a workspace container for aruna.
#
# If container already is running, attach to it.
#
declare VERSION=1.0.5
declare NAME=aruna-cities-$VERSION
declare IMAGE=hkjn/workspace:$(uname -m)-aruna-$VERSION
if docker ps | grep -q "${NAME}$"; then
	echo "Container $NAME is already running, attaching..".
	docker attach $NAME
	exit 1
fi
if docker ps -a | grep -q "${NAME}$"; then
	docker rm $NAME
fi

docker run --name $NAME -it \
       --cpu-shares 1024 \
       -w /home/user/src/github.com/arunaelentari \
       -v /home/aruna/.container_keys/:/home/user/.ssh \
       -v /home/aruna/src/github.com/arunaelentari:/home/user/src/github.com/arunaelentari \
       $IMAGE
