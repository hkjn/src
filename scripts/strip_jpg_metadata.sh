#!/usr/bin/env bash
#
# Strip metadata from JPG images.
#
set -euo pipefail

if ! which jhead 1>/dev/null; then
  echo "$(basename ${0}): this script needs jhead" >&2
  exit 1
fi
jhead -purejpg ${@}
