#!/bin/bash
#
# Logging library for bash scripts.
#
# All logging is done to stderr, since it's common to use stdout for
# the value that a function "returns" (since actual return values in
# bash are just integers). Logging to stderr allows for e.g:
# getFoo() {
#   debug "Getting foo.."
#   echo "foo=42"
# }
# f=$(getFoo)
#
# TODO: Add more colors, taken from:
# for c in {0..255}; do echo -e "\e[38;05;${c}m ${c} Bash Color Chart"; done
#

LOG_PREFIX="$(basename $0)"

# No colo(u)r.
NC='\033[0m'

# ANSI escape codes, via http://stackoverflow.com/a/28938235.
RED='\033[0;31m'
LRED='\033[1;31m'
URED='\033[4;31m'
HIRED='\033[0;101m'

LGREEN='\033[0;32m'

LBLUE='\033[0;34m'
BPURPLE='\033[1;34m'
CYAN='\033[0;36m'
LCYAN='\033[1;36m'

LGRAY='\033[0;37m'
WHITE='\033[1;97m'

BROWN='\033[0;33m'
YELLOW='\033[1;33m'
UBROWN='\033[4;33m'

# Since functions usually can't tell their callers to exit, we'll
# instead register a signal handler for SIGTERM here, then send
# ourselves the signal if fatal() is called.
trap "exit 1" TERM
export LOGGER_PID=$$
export LOGGING_LEVEL=${LOGGING_LEVEL:-2}

# fatal prints given error messages to stderr, then exits the top-level script.
fatal() {
	echo -e "${URED}[$LOG_PREFIX]${NC} ${LRED}FATAL: $@${NC}" >&2
	kill -s TERM $LOGGER_PID
}

# error prints given info messages to stderr.
error() {
	echo -e "${URED}[$LOG_PREFIX]${NC} ${LRED}$@${NC}" >&2
	return 0
}

# info prints given messages to stderr.
info() {
	echo -e "${LCYAN}[$LOG_PREFIX]${NC} ${CYAN}$@${NC}" >&2
	return 0
}

# infon prints given info messages to stderr, without echoing the newline.
infon() {
	echo -ne "${LCYAN}[$LOG_PREFIX]${NC} ${CYAN}$@${NC}" >&2
	return 0
}

# warn prints given warning messages to stderr.
warn() {
	echo -e "${YELLOW}[$LOG_PREFIX] Warning: $@${NC}" >&2
	return 0
}

# debug prints given debug messages to stderr.
debug() {
	echo -e "${CYAN}[$LOG_PREFIX]${NC} ${LGREEN}$@${NC}" >&2
	return 0
}

# debugV prints given verbose debug messages to stdout.
debugV() {
	[[ "$LOGGING_LEVEL" -ge 3 ]] && echo -e "${CYAN}[${LOG_PREFIX}V]${NC} ${LGREEN}$@${NC}" >&2
	return 0
}


# confirm shows a confirmation prompt, and returns success only if the user confirms the action.
confirm() {
	local skipConfirmations=${SKIP_CONFIRMATIONS:-""}
	[[ "$skipConfirmations" ]] && return 0
	local msg="Is this your intent?"
	[[ "$#" -eq 1 ]] && msg="$1"
	infon "$msg [y/N] "
	read -r -p "" response
	[[ $response =~ ^([yY][eE][sS]|[yY])$ ]] || {
		info "Ok."
		return 1
  }
	return 0
}
