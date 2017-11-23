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

# TODO(hkjn): Should report in that this host was bootstrapped:
# - set up systemd .timer + .service to publish message that this host is available to MQ system on foo.hkjn.me
# - include the ssh-keygen -l -f /etc/ssh/ssh_host_ecdsa_key

