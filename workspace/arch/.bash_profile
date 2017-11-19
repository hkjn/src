# Shell Options.
shopt -s checkwinsize
# Append to the history file, don't overwrite it.
shopt -s histappend

# don't put duplicate lines in the history. See bash(1) for more options
# don't overwrite GNU Midnight Commander's setting of `ignorespace'.
export HISTCONTROL=$HISTCONTROL${HISTCONTROL+,}ignoredups
# ... or force ignoredups and ignorespace
export HISTCONTROL=ignoreboth

# Color definitions.
COL_YEL="\[\e[1;33m\]"
COL_GRA="\[\e[0;37m\]"
COL_WHI="\[\e[1;37m\]"
COL_GRE="\[\e[1;32m\]"
COL_RED="\[\e[1;31m\]"
COL_NOR='\[\033[00m\]'

PROMPT_COMMAND=__prompt_command
# Set prompt according to exit status and other info.
__prompt_command() {
	exitstatus="$?"

	local usercol
	local psymb
	if [[ "$UID" -eq 0 ]]; then
		usercol=$COL_RED
		psymb="#"
	else
		usercol=$COL_GRE
		psymb="$"
	fi

	local symbcol
	if [[ $exitstatus != 0 ]]; then
		symbcol="$COL_RED"
	else
		symbcol="$COL_GRE"
	fi

	local host
	host=$(hostname -f)
	local prompt
	prompt="${usercol}\u${COL_WHI}@${COL_YEL}${host}${COL_WHI}:\w\n${symbcol}${psymb} $COL_NOR"
	PS1="$prompt"
}

# Starts ssh-agent and stores the SSH_AUTH_SOCK / SSH_AGENT_PID for
# later reuse.
start_ssh_agent() {
  ssh-agent -s > ~/.ssh-agent.conf 2> /dev/null
  source ~/.ssh-agent.conf > /dev/null
}

# Loads SSH identities (starting ssh-agent if necessary), recovering
# from stale sockets.
load_ssh_keys() {
  # SSH-agent setup adapted from
  # http://superuser.com/questions/141044/sharing-the-same-ssh-agent-among-multiple-login-sessions.

  # Time a key should be kept, in seconds.
  key_ttl=$((3600*8))
  if [[ ! -f ~/.ssh-agent.conf ]]; then
    # No existing config, start agent.
    start_ssh_agent
    ssh-add -t $key_ttl > /dev/null 2>&1
    return 0
  fi

  # Found previous config, try loading it. This sources in the path to
  # the authentication socket (SSH_AUTH_SOCK, used below).
  source ~/.ssh-agent.conf > /dev/null
  # List all identities the SSH agent knows about.

  local stat
  ssh-add -l > /dev/null 2>&1
  stat=$?
  if [[ $stat -eq 1 ]]; then
    # $?=1 means the socket is there but contains no key.
    ssh-add -t $key_ttl > /dev/null 2>&1
  elif [[ $stat -eq 2 ]]; then
    # $?=2 means the socket is not there or broken
    rm -f $SSH_AUTH_SOCK
    start_ssh_agent
    ssh-add -t $key_ttl > /dev/null 2>&1
  fi
  # $?=0 means the socket is there and it has a key.
  return 0
}

export PATH=$PATH:/home/user/bin
export EDITOR=nano
alias pp="git pull && git push"
alias gs="git status"
alias gdc="git diff --cached"

# Load SSH keys in new session.
load_ssh_keys
