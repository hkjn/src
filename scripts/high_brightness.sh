#!/usr/bin/env bash
#
# Increase brightness.
#
set -euo pipefail
echo "${0}: starting.." >> /tmp/brightness.log
mkdir -p /opt/hack/
LIGHT_LEVEL=${LIGHT_LEVEL:-"5"}
[[ -e /opt/hack/light ]] && LIGHT_LEVEL=$(cat /opt/hack/light)
LIGHT_LEVEL=$((LIGHT_LEVEL+1))
[[ "${LIGHT_LEVEL}" -gt 10 ]] && LIGHT_LEVEL=10
echo "${0}: high bumped to ${LIGHT_LEVEL}" >> /tmp/brightness.log
echo ${LIGHT_LEVEL} > /opt/hack/light

# note: special-case hack for LIGHT_LEVEL=10 -> BRIGHTNESS=0.10
BRIGHTNESS="0.${LIGHT_LEVEL}"
[[ "${LIGHT_LEVEL}" -gt 9 ]] && BRIGHTNESS="1.${LIGHT_LEVEL}"
echo "${0}: setting xrandr to ${BRIGHTNESS}" >> /tmp/brightness.log
xrandr --output eDP-1 --brightness ${BRIGHTNESS}
TEMPERATURE=$((${LIGHT_LEVEL} * 1000+500))
echo "${0} setting temp to ${TEMPERATURE}" >> /tmp/brightness.log
redshift -O ${TEMPERATURE}
