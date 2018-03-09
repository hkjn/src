#!/bin/sh

set -euo pipefail

declare IP_ADDR=${IP_ADDR:-"163.172.162.18"}
declare LOG_LEVEL=${LOG_LEVEL:-"debug"}
declare ALIAS=${ALIAS:-"ln.hkjn.me"}

if [ ! "${IP_ADDR}" ]; then
	echo "FATAL: No IP_ADDR specified." >&2
fi

echo "Starting bitcoind.."
bitcoind -dbcache=800 -onlynet=ipv4 -printtoconsole &

echo "Starting lightningd.."
lightningd --network=bitcoin --ipaddr=${IP_ADDR} --log-level=${LOG_LEVEL} --alias=${ALIAS} --rgb=003366
