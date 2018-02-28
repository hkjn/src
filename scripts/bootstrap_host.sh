#!/bin/bash
#
# Runs minimal steps to bootstrap a new host to a good state.
#

set -euo pipefail

die() {
  echo "FATAL: $@" >&2
  exit 1
}

[[ $UID -eq 0 ]] || die 'Needs to be root on remote host to bootstrap.'

RUSER=${RUSER:-"zero"}
source /etc/os-release
ID_LIKE=${ID_LIKE:-""}
ID=${ID:-""}
if [[ "$ID_LIKE" = "archlinux" ]] || [[ "$ID" = "coreos" ]]; then
	useradd -G docker --create-home --shell /bin/bash $RUSER
elif [[ "$ID_LIKE" = "debian" ]]; then
	if ! which docker 1>/dev/null; then
		apt-get -y update
		apt-get -y upgrade
		apt-get install -y apt-transport-https ca-certificates curl software-properties-common
		curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
		if ! apt-key fingerprint 0EBFCD88 2>/dev/null | grep "9DC8 5822 9FC7 DD38 854A E2D8 8D81 803C 0EBF CD88" 1>/dev/null; then
			echo "FATAL: Bad fingerprint of Docker gpg key." >&2
			exit 1
		fi
		add-apt-repository "deb [arch=armhf] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
		apt-get -y update
		apt-get -y install docker-ce
	fi
	adduser --ingroup docker --shell /bin/bash --disabled-password $RUSER
else
	adduser -G docker -s /bin/bash $RUSER
fi
mkdir -p /home/$RUSER/.ssh
cp $HOME/.ssh/authorized_keys /home/$RUSER/.ssh/
chown -R $(id -u $RUSER):$(id -g $RUSER) /home/$RUSER/
chmod 700 /home/$RUSER/.ssh
chmod 400 /home/$RUSER/.ssh/authorized_keys

if grep -q 22 /etc/ssh/sshd_config; then
	sed -e s/22/6200/ \
	    -i /etc/ssh/sshd_config
else
	echo 'Port 6200' >> /etc/ssh/sshd_config
fi
if grep -q PermitRootLogin /etc/ssh/sshd_config; then
	sed -e s/'PermitRootLogin without-password'/'PermitRootLogin no'/ \
	    -i /etc/ssh/sshd_config
else
	echo 'PermitRootLogin no' >> /etc/ssh/sshd_config
fi
if which systemctl 1>/dev/null; then
	systemctl restart sshd
else
	service sshd restart
fi
passwd -d $RUSER
if ! which sudo 1>/dev/null; then
	apk add sudo
fi
echo "$RUSER ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/user_sudo
chmod 0440 /etc/sudoers.d/user_sudo

# TODO: Following seems to start interactive shell on ubuntu; should
# run following commands automatically..
su - $RUSER
mkdir -p src/hkjn.me
cd src/hkjn.me
git clone https://github.com/hkjn/scripts.git
git clone https://github.com/hkjn/dotfiles.git
cd dotfiles
cp .bash* ~/

echo 'Done bootstrapping host.'
