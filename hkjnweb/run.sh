#!/bin/bash
#
# Runs the web server using production settings.
#
set -euo pipefail

export BIND_ADDR=127.0.0.1:8089
export SERVE_HTTP=1
export PROD=1
echo "[run.sh] Building binary.."
go build ./cmd/server/hkjnserver.go
while pgrep hkjnserver 1>/dev/null; do
	pid=$(pgrep hkjnserver)
	echo "[run.sh] Sending SIGTERM to existing process '$pid'.."
	kill $pid
	sleep 1
done
echo "[run.sh] Starting server.."
./hkjnserver -alsologtostderr
