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
command -v ssh-agent >/dev/null 2>&1 || { echo >&2 "ssh-agent is missing"; exit 1; }

#
# Note: The PGP message below encodes a SSH private key, using
# a weak passphrase (revealed in the lectures). It's a bad idea
# to share private keys in general! We do it during the lectures
# to allow students to easily access a shared server, so they can
# interact with a VM with a fully synced Bitcoin Core instance
# (requiring ~265 GB+ data including txindex at time of writing)
# without requiring everyone to first set up individual keys.
#
cat << EOF > 21_student_id_rsa.asc
-----BEGIN PGP MESSAGE-----

jA0ECQMCY8nrY7/CrvT20uoBGARvC+o3WRRmFqRXRhd21LmcbKild21Rar7fJYkY
q9b+qHbi2A+WORGD0m+bTFYibwG47GEc8Se6CLjJSN7LB5GIf05SQJH8jbFxYf1d
tH9jcMJJRl9HAyZ3sRBlFJgUYTmPdded6+E3VfEh0DgHr03yA1zL72C5ho0XBcj4
t1BsFEhM9TQbyypuBKrd3tP3w+JJ1xNVjwnHRvkFuZRN25Ikz5hfa9cnrtLSnqTA
kRdNb9kuYsZ4F9ZdHoVIOlw+fc8GR9sm0Q0sLRxWSicLMUnIH9Mvq59l2Q9dumVY
IFKW3ApxyUBLG24C+P51SG6lMLW/W934Ogqm46pOKiyEvtPWqU7F1h8zpEVMoCPA
MDeyY+KpI8GXX+Du+YGqUtknIYQTC6kmNxaW5jSsGoGciW2wtUFTxqo5kaLOCSEI
wJlm9i4iW+4GDkROhBytEyYVk5W5UTTld9JhN/Dq+zfMFHW9GeZx1D3OGj5U+nJk
Xl8LvejwXdR4WurWPzf680+Mim2nj8QgEVuTbJN2aLezki6J55C5aDIl2RmWdvQk
jjZ1UgKYizsc++1gxBjqjEbXkRRhHPCVOasGOwWULh/rRonubXLAe1kD31+l2zUn
jEtwcv8r/bMTkNae/V4h5932OuyS9NOXOo4usiQv/bogaYaskR4YejKibd/orOPp
8gSmXVQ8wZo4CnGTetcJuZLvEVdgenOhUJrDPyAEt09rbkfuhXpCK9PjZ6kDmi+T
iPoo6kiX7oUK8b8VZm+jqg+8E0kCQXtiY1pQuB8t1GtZsjZX3oleKiW+iXNhUUVY
N/Fk4qfWRTtwZUqRio6uljxBUweD0tWvOqJxm4U4XauvA8mGreQQFq0LoGVz/5G2
IS+Z/blqN4k7sjreEYNfIgMmwr/ffAAtxsbSW3svu4lwns4PV/VMXEEc9ESj3jkI
YV4PJ3mvfsTW5lJE6iTH+JVOFOQD/MuQtt5jmlSK12iO9g796Yn7nUcEpQKH6sqC
6vJYsVjpHGCtaqSyK9RvtRgJyCXzeUTj6Ks7gFNlykcfBDQKjDuxVWNJQIIGUnKX
heFKkwAWjUf5POq9ke9xsXKvXJcs77A9N+kpnd5Sb+NYfmn3/OR17xZAitt1M4Zw
y8wLRfd5kRwI0zDYn0BnWAebXact+X/TViF+iVr6zd5eMY3T2RPUintJsgR7XxFK
TBPZDyxU45d/SHsajRQdAYd0fKjbaarRIIf0/fKwi8ZzpPmH3oMPHZ40kB4QwuFW
16URSCdzDpMLmYg+TedKNxTopvZ9JuE8wh7BktFMMQBI3rfZoo/t2JLI2arRvdWa
NEBtrcjzpMOsOHvwqi5fMqI+PsxBG1m1TtwOY8BojZl4wJjHzLbfRUa3M+u8/WZF
B2XpqfK+ajasqSvQ2E2nQUUwvkkaTdb580q4XRII03JE3sUqugRDs6Cue0jtQtHi
+SeCKP64dp3hOsZDa6VlSylTqzzrvOs8K//857Y2aYn6Rjyko2LeMZjzwzTtDOEj
7LwXSEy4tpbu2SACy4tyDiY0nHmVzsOQFzUvjgzhiqSYNjjQNSjFMHAFTgCSTw/W
AAR3OeYGemTbWpFBt9BuLYth1tn0wtkIhbjE4fZVCwjrgD4cVK2ZSxaV5xSGlPWa
ESA/rZ1uihVPfVOCfytQv+47rTIrjS15mkTb/o7qrjYloSX1vxzUfEAfZGstmf5a
3I2jUb5S+lttKjnVe6Zrx23OLjhsxCRH/zaFSrxW/n5IZkHCeF/4rniju0oVADk8
gNzjfwbqAjZwUtWPfdTyjj3/H1OIncU+Kso/IRfHGfuf2bORvw5RhgVA6g==
=hnO4
-----END PGP MESSAGE-----
EOF

# Encrypted file created with:
#   gpg -o msg.asc --armor --symmetric --cipher-algo AES256 21_student_id_rsa
[[ -e 21_student_id_rsa.asc ]] || { echo >&2 "No 21_student_id_rsa.asc key found."; exit 1; }

echo "Creating ~/.ssh directory, if necessary.."
mkdir -p ~/.ssh/

[[ -e ~/.ssh/21_student_id_rsa ]] || {
    echo "Attempting to decrypt SSH key ~/.ssh/21_student_id_rsa. Hint: genesis block."
    gpg -o ~/.ssh/21_student_id_rsa -d 21_student_id_rsa.asc
}

echo "Setting permissions expected by SSH for ~/.ssh/21_student_id_rsa.."
chmod 400 ~/.ssh/21_student_id_rsa

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

echo "Adding SSH key ~/.ssh/21_student_id_rsa (may prompt for passphrase to lock privkey file with)."
eval $(ssh-agent)
ssh-add ~/.ssh/21_student_id_rsa

echo "All done! Try connecting with:"
echo "  ssh 21"
