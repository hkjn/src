#!/bin/bash
#
# Tools for bash scripts
#

LOG_PREFIX="$(basename $0)"

# exists returns 0 if the given command exists on PATH.
exists() {
  if which $1 1>/dev/null; then
    return 0
  fi
  return 1
}
