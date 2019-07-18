# ~/.bashrc: executed by bash(1) for non-login shells.

# If not running interactively, don't do anything
[ -z "$PS1" ] && return

# don't put duplicate lines in the history. See bash(1) for more options
# don't overwrite GNU Midnight Commander's setting of `ignorespace'.
export HISTCONTROL=$HISTCONTROL${HISTCONTROL+,}ignoredups
# ... or force ignoredups and ignorespace
export HISTCONTROL=ignoreboth

# Append to the history file, don't overwrite it.
shopt -s histappend

# for setting history length see HISTSIZE and HISTFILESIZE in bash(1)
HISTSIZE=1000
HISTFILESIZE=2000

# check the window size after each command and, if necessary,
# update the values of LINES and COLUMNS.
shopt -s checkwinsize

color_prompt=yes
if ! [ -x /usr/bin/tput ] || ! tput setaf 1 >&/dev/null; then
   # We have no color support; not compliant with Ecma-48
   # (ISO/IEC-6429). (Lack of such support is extremely rare, and such
   # a case would tend to support setf rather than setaf.)
   color_prompt=
fi

# echo the current git branch
gitbranch() {
  git branch 2> /dev/null | sed -e '/^[^*]/d' -e 's/* \(.*\)/\1/'
}

# echo user and host
userhost() {
  local lgreen='\[\033[01;32m\]'
  local normal='\[\033[00m\]'
  local dgray='\[\033[1;30m\]'
  echo "${lgreen}\u@\h${dgray}♾${normal}"
}

# echo current working directory
workdir() {
  local lblue='\[\033[01;34m\]'
  local normal='\[\033[00m\]'
  echo "${lblue}\w${normal}"
}

# echo extra info, if available
#
# TODO(hkjn): Improve the setup for "extrainfo" and make it extensible:

# 1. Read from some inmemory store / unix socket or similar, so no
#    filesystem access is necessary just to draw the prompt
# 2. Have separate timer-based job that checks stuff and writes to socket:
#    - connectivity (can ping 8.8.8.8, DNS resolution works, VPN is up/down)
#    - number of running docker containers (specifically, hkjn/fileserver containers)
extrainfo() {
  local lcyan='\[\033[1;36m\]'
  local normal='\[\033[00m\]'
  local red='\[\e[0;31m\]'
  local awscreds=''
  local dgray='\[\033[1;30m\]'
  if [[ ! -e '.aws/creds.env' ]]; then
    return
  fi
  local expiry=$(grep -Eo '[0-9]{10}$' .aws/creds.env)
  local secsleft=$(($expiry-$(date +%s)))
  if [[ $secsleft -lt 3600 ]]; then
    # Expired credentials. (Apparently.)
    awscreds="💩"
  elif [[ $secsleft -lt 4000 ]]; then
    # Credentails about to expire.
    awscreds="💣"
  else
    # Good credentials.
    awscreds="${lcyan}✓"
  fi

  local stage=${STAGE:-""}
  local stageinfo="?"
  if [[ "$stage" == *"prod"* ]]; then
    stageinfo="${lcyan}☠"
  elif [[ "$stage" == *"stag"* ]]; then
    stageinfo="${lcyan}⚠"
  else
    stageinfo="${lcyan}☯"
  fi
  echo "${dgray}♾${awscreds}${dgray}♾${stageinfo}${normal}"

}

PROMPT_COMMAND=__prompt_command
# Set prompt according to exit status and other info.
__prompt_command() {
  local EXIT="$?"

  local prompt="►"

  local dgray='\[\033[1;30m\]'
  local green='\[\033[00;32m\]'
  local lwhite='\[\033[01;11m\]'
  local normal='\[\033[00m\]'
  local red='\[\e[0;31m\]'

  # TODO: red/green here doesn't seem to exist consistently (alpine),
  # even though other colors do. Switch?
  local pcolor="$green"
  if [ $EXIT != 0 ]; then
    pcolor="$red"
  fi
  # 木 人 ♪
  if [ "$color_prompt" = yes ]; then
    PS1="${lwhite}\$(gitbranch)${dgray}♾${normal}$(userhost)$(workdir)$(extrainfo) \n${pcolor}${prompt}$normal "
  else
    PS1="$"
  fi
}

# If this is an xterm set the title to user@host:dir
case "$TERM" in
xterm*|rxvt*)
    PS1="\[\e]0;${debian_chroot:+($debian_chroot)}\u@\h: \w\a\]$PS1"		
    ;;
*)
    ;;
esac

# Pull in useful functions.
BASH_FUNCS="$HOME/src/hkjn.me/src/scripts/bash_funcs.sh"
if [[ -e "$BASH_FUNCS" ]]; then
	source "$BASH_FUNCS"
else
	echo "No '$BASH_FUNCS' found. Try 'go get hkjn.me/src'?"
fi

# enable color support of ls and also add handy aliases
if [ -x /usr/bin/dircolors ]; then
    eval "`dircolors -b`"
    alias ls='ls --color=auto'
fi

_rmkey() {
	ssh-keygen -R ${1}
	local ip
	ip=$(dig +short ${1})
	ssh-keygen -R ${ip}
}

# enable programmable completion features (you don't need to enable
# this, if it's already enabled in /etc/bash.bashrc and /etc/profile
# sources /etc/bash.bashrc).
[ -f /etc/bash_completion ] && source /etc/bash_completion

alias pp="git pull --rebase && git push"
alias gdc="git diff --cached"
alias gs="git status"

alias docker_rmcontainers='docker rm $(docker ps -a -q -f status=exited)'
alias docker_rmall='docker rm -f $(docker ps -a -q) && docker rmi $(docker images -q --filter "dangling=true")'
alias docker_rmimages='docker rmi $(docker images -q --filter "dangling=true")'

alias e="vim $1"
alias ec="e $HOME/.bash_profile"
alias ecf="e $HOME/.ssh/config"
alias rf="[ -e $HOME/.bash_profile ] && source $HOME/.bash_profile || source $HOME/.bashrc"
alias tc="tmux new -s $1"
alias ta="tmux attach -d -t $1"
alias ll='ls -hsAl'
alias mp="mplayer -af scaletempo $@"
alias mp50="mplayer -af scaletempo -fs -panscanrange -5 $@"
alias xclip="xclip -selection c"
alias shlogs="less ${HOME}/.shell_logs/${HOSTNAME}"
alias rmkey="_rmkey ${1}"
alias elec='electrum --oneserver --server=127.0.0.1:50001:t'
_gs() {
  gpg --decrypt ${1} | tr -d '\n'
}

export LANG="en_US.UTF-8"
export LC_CTYPE="en_US.UTF-8"
export EDITOR=nano
export GOPATH=${HOME}
export PATH=/usr/local/go/bin:${GOPATH}/src/hkjn.me/src/scripts:${HOME}/bin:${HOME}/.local/bin:${HOME}/.cargo/bin:/snap/bin:${HOME}/.npm-global/bin:.:$PATH
export PYTHONPATH=.:..

# GPG always wants to know what TTY it's running on.
export GPG_TTY="$(tty)"

# Set the authentication socket for ssh-agent to gpg-agent.
export SSH_AUTH_SOCK="/run/user/$UID/gnupg/S.gpg-agent.ssh"

# When running nmon, by default show:
# - long-term CPU averages (l)
# - memory & swap          (m)
# - kernel stats & loadavg (k)
# - top processes          (t)
export NMON=lmkt

# Don't scatter __pycache__ directories all over the place.
export PYTHONDONTWRITEBYTECODE=1

export CLOUDSDK_PYTHON=python2

if [ -d ~/google-cloud-sdk ]; then
	# The next line updates PATH for the Google Cloud SDK.
	source "$HOME/google-cloud-sdk/path.bash.inc"

	# The next line enables bash completion for gcloud.
	source "$HOME/google-cloud-sdk/completion.bash.inc"
fi

# Allow current user to connect to X11 socket from any host; required
# to run graphical Docker containers. But not on OS X, since even
# though xhost exists, on at least OS X Mavericks it just stalls
# indefinitely if invoked, preventing new bash sessions.
if which xhost > /dev/null 2>&1 && [ ! -z $DISPLAY ] && [ $(uname) != "Darwin" ]; then
		xhost +si:localuser:$USER >/dev/null
		xhost +si:localuser:root >/dev/null
fi

if [[ -e "${HOME}/.bash_profile_extra" ]]; then
	source "${HOME}/.bash_profile_extra"
fi

kill_other_mosh_sessions() {
	local pids
	pids=$(pgrep mosh-server | grep -v $(ps -o ppid --no-headers $$))
	if [[ "${pids}" ]]; then
		kill ${pids}
	fi
}

declare TMUX=${TMUX:-""}
declare ATTACH_TMUX=${ATTACH_TMUX:-""}
if [[ ! "${TMUX}" ]] && [[ "${ATTACH_TMUX}" ]]; then
	# Kill any other mosh sessions that might be lingering.
	kill_other_mosh_sessions
	# If not already in tmux session, attach to it.
	tmux attach
fi
