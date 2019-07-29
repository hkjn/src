#!/usr/bin/env bash

set -euo pipefail

sudo mkdir -p /media/wwn /media/wwn_exfat
[[ -e /dev/mapper/wwn-unlocked ]] || {
  sudo cryptsetup luksOpen --key-file=/root/keys/ata-SanDisk_SD9SN8W1T00_183255420707-part3.key /dev/disk/by-id/ata-SanDisk_SD9SN8W1T00_183255420707-part3 wwn-unlocked
}

sudo mount /dev/disk/by-id/ata-SanDisk_SD9SN8W1T00_183255420707-part1 /media/wwn_exfat
sudo mount /dev/mapper/wwn-unlocked /media/wwn/
echo "Mounted /media/wwn."
