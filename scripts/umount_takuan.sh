#!/usr/bin/env bash

set -euo pipefail

sudo umount /media/takuan
sudo cryptsetup luksClose wwn-unlocked
echo "Unmounted /media/takuan."
