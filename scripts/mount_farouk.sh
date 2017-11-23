#!/bin/bash

# Mounts the "farouk" logical volume under the /dev/iago volume group:
# /dev/iago/farouk: LUKS, with passphrase
#
# This script requires root permissions.

set -e

# TODO: Do this in /etc/fstab + udev instead.
cryptsetup luksOpen /dev/iago/farouk farouk_clear
mount /dev/mapper/farouk_clear /media/farouk
echo "Mounted /media/farouk."
