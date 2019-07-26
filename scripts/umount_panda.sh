#!/usr/bin/env bash

set -euo pipefail

sudo umount /media/panda
echo "Unmounted /media/panda."
echo "Running fsck.."
sudo fsck /dev/mapper/panda
sudo cryptsetup remove panda

sudo umount /media/panda_son
echo "Unmounted /media/panda_son."
echo "Running fsck.."
sudo fsck /dev/disk/by-id/usb-WD_Elements_SE_25FF_5758333144343835535A304A-0:0-part6

echo "All done."
