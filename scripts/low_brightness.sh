#!/usr/bin/env bash
#
# Increase brightness.
#
set -euo pipefail
echo "${0}: starting.." >> /tmp/brightness.log
mkdir -p /opt/hack/
LIGHT_LEVEL=${LIGHT_LEVEL:-"5"}
[[ -e /opt/hack/light ]] && LIGHT_LEVEL=$(cat /opt/hack/light)
LIGHT_LEVEL=$((LIGHT_LEVEL-1))
[[ "${LIGHT_LEVEL}" -lt 0 ]] && LIGHT_LEVEL=0
echo "${0}: high lowered to ${LIGHT_LEVEL}" >> /tmp/brightness.log
echo ${LIGHT_LEVEL} > /opt/hack/light

echo "${0}: setting xrandr to 0.${LIGHT_LEVEL}" >> /tmp/brightness.log
xrandr --output eDP-1 --brightness 0.${LIGHT_LEVEL}
TEMPERATURE=$((${LIGHT_LEVEL} * 1000+500))
echo "${0}: setting temp to ${TEMPERATURE}" >> /tmp/brightness.log
redshift -O ${TEMPERATURE}
