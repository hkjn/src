#!/bin/bash
#
# Back up current kernel and upgrade to new one.
#

set -ueo pipefail

LIBS=$(ls -d /lib/modules/?.*-ARCH | tail -n 1)
VERSION=$(echo $LIBS | cut -d '/' -f4)
BACKUP=/root/kernel-backups/${VERSION}-$(date +%Y%m%d)
BOOT=/boot/EFI/arch
KERNEL=${BOOT}/vmlinuz-arch-crypto

# Back up kernel src.
# NOTE: No kernel src seems to be in /usr/src on this system.
# cp -r /usr/src/linux-${VERSION} /usr/src/linux-2.6.28-ARCH-old

echo "Backing up previous kernel modules ${LIBS} to ${BACKUP}/ and ${LIBS}-stable/.."
sudo cp -r ${LIBS} ${BACKUP}/
# We create a copy with a different name, since the upgrade writes
# over the directory otherwise.
sudo cp -r ${LIBS} ${LIBS}-stable/

echo "Backing up previous compiled kernel to ${BACKUP}/${KERNEL}.efi and rotating it to ${KERNEL}-stable.efi.."
sudo cp -v ${KERNEL}.efi ${BACKUP}/
sudo cp -v ${KERNEL}.efi ${KERNEL}-stable.efi

# Back up initial ram image and fallback (initrd).
echo 'Backing up + rotating previous initrd (and fallback) image..'
sudo cp -v ${BOOT}/initramfs-arch-crypto.img ${BACKUP}/
sudo cp -v ${BOOT}/initramfs-arch-crypto.img ${BOOT}/initramfs-arch-crypto-stable.img
sudo cp -v ${BOOT}/initramfs-arch-crypto-fallback.img ${BACKUP}/
sudo cp -v ${BOOT}/initramfs-arch-crypto-fallback.img ${BOOT}/initramfs-arch-crypto-fallback-stable.img

# At this point we can upgrade.  Upgrading the kernel runs `mkinitcpio
# -p linux`, which regenerates /boot/initramfs-linux{-fallback}.img).
echo 'Upgrading system..'
sudo pacman -Syyu

echo "Rotating back old ${LIBS}-stable to ${LIBS}.."
sudo mv ${LIBS}-stable ${LIBS}

# TODO: Make default mkinitcpio preset linux directly put the files
# where we want them, rather than this hack.
echo "Copying in new kernel.."
sudo mv -v /boot/vmlinuz-linux ${KERNEL}.efi

echo "Copying in initial ram images.."
sudo mv -v /boot/initramfs-linux.img ${BOOT}/initramfs-arch-crypto.img
sudo mv -v /boot/initramfs-linux-fallback.img ${BOOT}/initramfs-arch-crypto-fallback.img


