#!/usr/bin/env bash

set -euo pipefail

sudo umount /media/wwn
echo "Unmounted /media/wwn."
echo "Running fsck.."
sudo fsck /dev/mapper/wwn-unlocked
sudo cryptsetup remove wwn-unlocked

sudo umount /media/wwn_exfat
echo "Unmounted /media/wwn_exfat."
echo "Running fsck.."
sudo fsck /dev/disk/by-id/ata-SanDisk_SD9SN8W1T00_183255420707-part1

echo "All done."
