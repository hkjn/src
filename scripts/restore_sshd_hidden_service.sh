#!/usr/bin/env bash
set -euo pipefail
if ! grep -q sshd /etc/tor/torrc; then
    echo "Restoring /var/lib/tor/sshd.."
    cat << EOF >> /etc/tor/torrc
HiddenServiceDir /var/lib/tor/sshd/
HiddenServiceVersion 3
HiddenServicePort 22 127.0.0.1:22
EOF
fi