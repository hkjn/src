#!/bin/bash
#
# Sample script for setting up LVM on Ubuntu.
#
set -euo pipefail

echo "Installing packages.."
apt-get -y update
apt-get -y install gdisk lvm2

echo "Creating physical volumes from block devices.."
pvcreate /dev/nbd2
pvcreate /dev/nbd3

echo "Creating volume group crypt from physical volumes.."
vgcreate crypt /dev/nbd2
vgextend crypt /dev/nbd3

echo "Creating logical volume crypt0 on volume group crypt.."
lvcreate -l 100%FREE crypt -n crypt0
mkfs.ext4 /dev/mapper/crypt-crypt0

echo "Mounting /dev/mapper/crypt-crypt0 at /crypt.."
mkdir -p /crypt
mount /dev/mapper/crypt-crypt0 /crypt
