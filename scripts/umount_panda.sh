#!/bin/bash

set -euo pipefail

sudo umount /media/panda
echo "Unmounted /media/panda."
echo "Running fsck.."
sudo fsck /dev/mapper/panda
sudo cryptsetup remove panda

echo "All done."
