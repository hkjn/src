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

jA0ECQMCq1IfLU1mfVz20uoBhsEPT/CEEDqBG7TV2rzuCg7h6KtZxCl9cHQ13fLe
uctY4/MCnhM3LcD4zH1R5vdBi9I3y/UoLTkSQSG+c67MxyKhaCBat09XObCv0QoI
L8h+7W6vbmIRnkL/l5zQpEbH5AuzX+dtItLDMqH8MPCHA9Ham+0NyjKjeWR2ZPCv
KmqwBeGofKJGa3TtVW30J2Ph9LKUqe4BTTZUicsBflLfDRmspy2nCxybVoyXwJUw
xJ5D18pHQRHNN/NL3Gqr6nXmTQEbSYnAR3/BRGY2Z35OXL/7ORF+jeJBUcdP75iZ
h7hl+nliyT6xgQJ5gbFELdhlE7LStgt2Y2jk/bMtX9zWg63haGP6McJRL2u1VoDi
kdQV6j7wg/kTfmeKj7x07kyk0665efplCC5669PQPYL7Eo4hUcRDqV7QHGeemRz1
eHoUqRYJrWSGCsqBT2qEB1QkLjQ8JAZHJ8I+gxG9hETgxR78Q88bcf03KFd5kgNk
ZBSVJS+xmAE6X8GUjnW0ktWb0Y4qAgaN+55C2AuvKO0qcxB0zQBR6voQvtji3wka
+c4otSlPWAxIIWhlqkcauRmBJd3qysB9CS9l8Hw7r3jv53WhBV49zWb0fcF3935T
uqZR7C3JJYKvHW8SdH7zv01J7/ZW/kNtftpIZCWK0gBtPFDJRpMbGF8oF/SRzZi1
5rFX+ADnUCKc9GN1GdP9IBbJk1otaiSKACBLJU9Pw1zc8nZ+j643rnb9HWAe7vx8
0FZ/ho98zIRfTHkrirCg5T7Z34ggh87iTTF+1gbQ+YWHhV2cMb3Fj28b0oHuxYyO
ynZqpYTqFyNw+P73wScErgn1wbClLISurYVP0CHrjcyezdBaMFCO5Ale3YL/oaqH
jBnYQ23VL3PKtYAeL/fp732HvYZD2pmRUuNmW7Y6FnKT4P6IfxBFnwkJC3oS2crG
b3G+ofGaSOfesx0kVP0FI0J5zfnZN6HTofMglAPupfuEWTKhWfMriqkdNgCjuPJA
TVgh9+A1/jDtRooNUGKIBPYgMphpGLwZrspeptCuxci0r/Yhhh83xJAc67x9ehEF
7kFMZovOzIY5zCKkXBcFnhv/soQj9NwJN8UIbSYorFTbM2DLgrZooVR1vtxEeG6Q
JPYOlIEieozX35s+vECfK8bNw08rmABw4hdHGR3RVcZNqaE2NJpcwcGJLaBlO60i
PHRgf6lxyVL2+XQC238G9X6g5BZ5XMH0Rxo8TVHDoYYL3wZNy5Ew70TYZRyhiCbe
OcU4U9mAhUrcMWsidRhE5rgrJEO1POgpgX/PXl70Zuo9wpwD6IInGRU8f8F6FTJZ
rivwg4NDh2cLWoub0KM9knugqmtLQSs51VjBtmPzCzpewJvJDEm/RcmgL3fKmn25
Yk6oNVtWz2Ag0UX7xEsl9ej2VsKZz7i9l3jSoGE8IaPd3B4hd3+56pf1SVMnnMH0
LZtd3l3mR0Wi17MBYnr/NrZNym5hbDocyJCnavh56mSKAKjCy9gfQfGrYbsjtPp1
d9eKcmk/0k2XLGW9ERR9bvK6GEg2ff4nnkbziROo9Lv5ehjsGH94GQDb8Ox1OTkH
Vx5r9H/BB7lU9AkhJfZ4WX2bumaX6hvhcg1b+Nwfg8IkRcPzyYPbK03OdmYICWqK
P11iuooJO1zbEU+g3mxV8tn+9ne8sQgcMPF6UXcA4xXP5rG26o3yncyTWiHT89Db
BmipiN05iMaFA0E4pOC+I34NvHB4OMZSyxxbTxO4KLcS07zdpuVuYBLPewIyzWiR
MC1S3zpKK2/6Ip+5+TQNvEn+h8NHBHspNpboKnE8CGerUWQFpPFPqLD0AUDS9A==
=Hnux
-----END PGP MESSAGE-----
EOF

# Encrypted file created with:
#   gpg --armor --symmetric --cipher-algo AES256 21_student_id_rsa
[[ -e 21_student_id_rsa.asc ]] || { echo >&2 "No 21_student_id_rsa.asc key found."; exit 1; }

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

