#!/bin/bash

set -euo pipefail

sudo umount /media/wwn
echo "Unmounted /media/wwn."
echo "Running fsck.."
sudo fsck /dev/mapper/wwn-unlocked
sudo cryptsetup remove wwn-unlocked

echo "All done."
