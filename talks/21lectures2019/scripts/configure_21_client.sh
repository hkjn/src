#!/usr/bin/env bash
#
# Configure SSH client so we can connect to 21.hkjn.me server.
#
echo "Setting flags to exit script if any command fails or if any variable is undefined.."
set -eu

echo "Checking that commands we require are installed.."
command -v gpg >/dev/null 2>&1 || { echo >&2 "gpg is missing"; exit 1; }
command -v mkdir >/dev/null 2>&1 || { echo >&2 "mkdir is missing"; exit 1; }
command -v cat >/dev/null 2>&1 || { echo >&2 "cat is missing"; exit 1; }
command -v chmod >/dev/null 2>&1 || { echo >&2 "chmod is missing"; exit 1; }
command -v ssh-add >/dev/null 2>&1 || { echo >&2 "ssh-add is missing"; exit 1; }

# Encrypted file created with:
#   gpg --symmetric --cipher-algo AES256 21_student_id_rsa
[[ -e 21_student_id_rsa.gpg ]] || {
    echo >&2 "No 21_student_id_rsa.gpg key found. Make sure to run this from the scripts/ directory."
    exit 1
}

echo "Creating ~/.ssh directory, if necessary.."
mkdir -p ~/.ssh/

[[ -e ~/.ssh/student_id_rsa ]] || {
    echo "Attempting to decrypt SSH key ~/.ssh/21_student_id_rsa. Hint: genesis block."
    gpg -o ~/.ssh/21_student_id_rsa -d 21_student_id_rsa.gpg
}

echo "Setting permissions expected by SSH for ~/.ssh/21_student_id_rsa.."
chmod 400 ~/.ssh/student_id_rsa

echo "Adding SSH key ~/.ssh/21_student_id_rsa (may prompt for passphrase to lock privkey file with)."
ssh-add ~/.ssh/student_id_rsa

echo "Adding aliases to SSH config (if necessary).."
grep -q 21.hkjn.me ~/.ssh/config || cat << EOF >> ~/.ssh/config

# The entries below were added by connect_21.sh.

Host 21
    HostName 21.hkjn.me
    Port 2222
    User student
    IdentityFile ~/.ssh/21_student_id_rsa

Host 22
    HostName 22.hkjn.me
    Port 2222
    User student
    IdentityFile ~/.ssh/21_student_id_rsa
EOF

echo "All done! Try connecting with:"
echo "  ssh 21"

