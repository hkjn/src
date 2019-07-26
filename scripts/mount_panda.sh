#!/usr/bin/env bash

set -euo pipefail

[[ -e /dev/mapper/panda ]] || {
	sudo cryptsetup luksOpen --key-file=/root/keys/usb-WD_Elements_SE_25FF_5758333144343835535A304A-0\:0-part1.key /dev/disk/by-id/usb-WD_Elements_SE_25FF_5758333144343835535A304A-0\:0-part1 panda
}

sudo mkdir -p /media/panda
sudo mount /dev/mapper/panda /media/panda/
echo "Mounted /media/panda."

sudo mount /dev/disk/by-id/usb-WD_Elements_SE_25FF_5758333144343835535A304A-0:0-part6 /media/panda_son
echo "Mounted /media/panda_son."
