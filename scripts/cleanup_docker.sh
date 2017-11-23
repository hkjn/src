#!/bin/bash
#
# Cleans up all untagged Docker images.
#

echo "Removing all stopped containers.."
docker rm $(docker ps -a -q -f status=exited)
echo "Removing all untagged images.."
docker rmi $(docker images -q --filter "dangling=true")
