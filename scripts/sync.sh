#!/bin/bash

set -euo pipefail

HOST=${1:-""}
[[ "$HOST" ]] || {
	echo "FATAL: No HOST specified." >&2
	exit 1
}
DIR="$(pwd)"
echo "Syncing local directory $DIR to $HOST:src/hkjn.me."
PAUSE=${PAUSE:-2s}
while true; do
	rsync -az --exclude=.git/ --exclude=.sync_active $DIR $HOST:src/hkjn.me/
	sleep $PAUSE
done
