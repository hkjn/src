#!/bin/bash

# Unmounts the logical volumes under the /dev/dumbo volume group:
# /dev/dumbo/clown: LUKS
# /dev/dumbo/timothy: plain
#
# This script requires root permissions.

set -euo pipefail

# TODO: Also run periodic cron (30m?) in case of unexpected
# disconnect.

# Sync important data to local encrypted storage.
# rsync -avz /media/clown/src /media/farouk/
# rsync -avz /media/clown/notes /media/farouk/

sudo umount /media/clown
echo "Unmounted /media/clown."
echo "Running fsck.."
sudo fsck /dev/mapper/clown_clear
sudo cryptsetup remove clown_clear

sudo umount /media/timothy
echo "Unmounted /media/timothy."
echo "Running fsck.."
sudo fsck /dev/mapper/dumbo-timothy

echo "All done."
