#!/bin/bash

# Unmounts the "musashi" logical volume under the /dev/miyato volume group:
# /dev/miyato/musashi: LUKS
#
# This script requires root permissions.

set -e

if mountpoint -q /media/musashi; then
  echo "Unmounting /media/musashi"
  sudo umount /media/musashi
fi
if [[ -e /dev/mapper/musashi_clear ]]; then
  echo "Running fsck on /dev/mapper/musashi_clear.."
  sudo fsck /dev/mapper/musashi_clear
  echo "Removing /dev/mapper/musashi_clear.."
  sudo cryptsetup remove musashi_clear
fi
echo "All done."

