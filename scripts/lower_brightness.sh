#!/usr/bin/env bash
#
# Increase brightness.
#
set -euo pipefail
mkdir -p /opt/hack/
LIGHT_LEVEL=${LIGHT_LEVEL:-"5"}
[[ -e /opt/hack/light ]] && LIGHT_LEVEL=$(cat /opt/hack/light)
LIGHT_LEVEL=$((LIGHT_LEVEL-1))
[[ "${LIGHT_LEVEL}" -lt 0 ]] && LIGHT_LEVEL=0

echo ${LIGHT_LEVEL} > /opt/hack/light

xrandr --output eDP-1 --brightness 0.${LIGHT_LEVEL}
TEMPERATURE=$(expr 10000 - ${LIGHT_LEVEL} * 1000)
redshift -P -O ${TEMPERATURE}
