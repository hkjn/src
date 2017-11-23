#!/bin/bash

# Useful bash functions.

# Starts ssh-agent and stores the SSH_AUTH_SOCK / SSH_AGENT_PID for
# later reuse.
function start-ssh-agent() {
  ssh-agent -s > ~/.ssh-agent.conf 2> /dev/null
  source ~/.ssh-agent.conf > /dev/null
}

# Loads SSH identities (starting ssh-agent if necessary), recovering
# from stale sockets.

# TODO(henrik): Decrypt .ssh/*.pem.gpg files, load into SSH
# agent (ideally with no file touching disk at any time), make sure
# no plaintext files exist.
function load-ssh-key() {
	# SSH-agent setup adapted from
	# http://superuser.com/questions/141044/sharing-the-same-ssh-agent-among-multiple-login-sessions.

	# Time a key should be kept, in seconds.
	key_ttl=$((3600*8))
	if [ ! -f ~/.ssh-agent.conf ] ; then
		# No existing config, start agent.
		start-ssh-agent
		ssh-add -t $key_ttl > /dev/null 2>&1
		return 0
	fi
	# Found previous config, try loading it. This sources in the path to
	# the authentication socket (SSH_AUTH_SOCK, used below).
	source ~/.ssh-agent.conf > /dev/null
	# List all identities the SSH agent knows about.

	# TODO(henrik): Maybe check if output here contains all entries in a
	# known list, and if not, add the missing keys?
	ssh-add -l > /dev/null 2>&1
	stat=$?
	# $?=0 means the socket is there and it has a key.
	if [ $stat -eq 0 ]; then
		return 0
	elif [ $stat -eq 1 ] ; then
		# $?=1 means the socket is there but contains no key.
		ssh-add -t $key_ttl > /dev/null 2>&1
	elif [ $stat -eq 2 ] ; then
		# $?=2 means the socket is not there or broken
		rm -f $SSH_AUTH_SOCK
		start-ssh-agent
		ssh-add -t $key_ttl > /dev/null 2>&1
	fi
}

# Timed GTK dialogs; use like "timer 25m your note here".
function timer() {
	local N=$1; shift

	echo "timer set for $N"
	sleep $N && zenity --info --title="Time's Up" --text="${*:-BING}"
}

# Log all commands typed in host-specific file.
function command_log () {
	# Save the rv
	local -i rv="$?"
	# Get the last line local
	last_line="${BASH_COMMAND}"
	mkdir -p "$HOME/.shell_logs"
	local logfile="${HOME}/.shell_logs/${HOSTNAME}"
	local current_ts="$(date '+%Y%m%d %H:%M:%S')"
	if [ "$last_line" != '' ]; then
		echo "${current_ts} ${LOGNAME} Status[${rv}] SPID[${$}] PWD[${PWD}]" \
		      \'${last_line#        }\' >> "${logfile}"
	fi
}

# mw merges the current branch into specified branch.
function mw() {
	if [ "$#" -ne 1 ]; then
		echo "Usage: mw otherBranch" >&2
		return 1
	fi
	current="$(gitBranch)"
	other="$1"
	if [ ! "$current" ]; then
		echo "mw: Not in a git repo." >&2
		return 2
	fi
	git pull || return 3
	git push || return 4
	git checkout "$other"
	git pull || return 5
	git merge --no-ff "$current"
	git push || return 6
	git checkout "$current"
	return 0
}

if [ $(uname) == "Darwin" ]; then
	export LC_ALL=en_US.UTF-8
	export LANG=en_US.UTF-8
fi


# Trap + log commands.
trap command_log DEBUG

# Load SSH keys in new session.
load-ssh-key
