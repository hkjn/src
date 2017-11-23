#!/bin/bash

# Mounts the logical volumes under the /dev/dumbo volume group:
# /dev/dumbo/clown: LUKS
# /dev/dumbo/timothy: plain
#
# This script requires root permissions.

set -euo pipefail

# TODO: The following does not fail, even when the backing physical
# HDD is disconnected. Why? (We could check if UUID of backing
# physical disk exists first..)
[[ -e /dev/mapper/clown_clear ]] || {
  sudo cryptsetup luksOpen /dev/dumbo/clown clown_clear --key-file=/root/keys/clownkey.randomtext
}
sudo mount /dev/mapper/clown_clear /media/clown
echo "Mounted /media/clown."

sudo mount /dev/mapper/dumbo-timothy /media/timothy
echo "Mounted /media/timothy."

# Set time to power-down disk to 3 min (36 * 5 sec) without activity.
sudo hdparm -S 36 /dev/disk/by-id/wwn-0x5000c500656aa99b
