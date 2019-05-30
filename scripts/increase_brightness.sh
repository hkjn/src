#!/usr/bin/env bash
#
# Increase brightness.
#
set -euo pipefail
mkdir -p /opt/hack/
LIGHT_LEVEL=${LIGHT_LEVEL:-"5"}
[[ -e /opt/hack/light ]] && LIGHT_LEVEL=$(cat /opt/hack/light)
LIGHT_LEVEL=$((LIGHT_LEVEL+1))
[[ "${LIGHT_LEVEL}" -gt 10 ]] && LIGHT_LEVEL=10

echo ${LIGHT_LEVEL} > /opt/hack/light

BRIGHTNESS="0.${LIGHT_LEVEL}"
[[ "${LIGHT_LEVEL}" -gt 9 ]] && BRIGHTNESS="1.${LIGHT_LEVEL}"
xrandr --output eDP-1 --brightness ${BRIGHTNESS}
TEMPERATURE=$(expr 10000 - ${LIGHT_LEVEL} * 1000)
redshift -P -O ${TEMPERATURE}
