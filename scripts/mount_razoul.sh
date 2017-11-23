#!/bin/bash

# Mounts the "razoul" logical volume under the /dev/iago volume group:
# /dev/iago/razoul: LUKS, with keyfile
#
# This script requires root permissions.

set -e

# TODO: Do this in /etc/fstab + udev instead.
cryptsetup luksOpen /dev/iago/razoul razoul_clear --key-file=/root/keys/razoulkey.randomtext
mount /dev/mapper/razoul_clear /media/razoul
echo "Mounted /media/razoul."

