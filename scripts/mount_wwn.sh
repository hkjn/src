#!/usr/bin/env bash

set -euo pipefail

[[ -e /dev/mapper/wwn-unlocked ]] || {
  sudo cryptsetup luksOpen --key-file=/root/keys/ata-SanDisk_SD9SN8W1T00_183255420707-part3.key /dev/disk/by-id/ata-SanDisk_SD9SN8W1T00_183255420707-part3 wwn-unlocked
}

sudo mount /dev/mapper/wwn-unlocked /media/wwn/
echo "Mounted /media/wwn."
