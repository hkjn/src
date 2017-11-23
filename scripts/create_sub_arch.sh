#!/bin/sh

# Creates a 32bit arch install in /opt, which allows chroot, etc.
#
# This script requires root.
#
# Based on https://wiki.archlinux.org/index.php/Arch64_Install_bundled_32bit_system.
#
# Example usage:
# $0

set -e

TARGET=/opt/arch32

if [ -d "$TARGET" ]; then
		echo "$TARGET already exists; bailing"
		exit -1
fi

# Create target directory and copy over pacman mirrors.
mkdir $TARGET
sed -e 's/\$arch/i686/g' /etc/pacman.d/mirrorlist > $TARGET/mirrorlist
sed -e "s@/etc/pacman.d/mirrorlist@$TARGET/mirrorlist@g" -e '/Architecture/ s,auto,i686,'  /etc/pacman.conf > $TARGET/pacman.conf

# Create core /var directories.
echo "Creating /var directories.."
mkdir -p $TARGET/var/{cache/pacman/pkg,lib/pacman}

# Sync the pacman of the chroot.
pacman --root $TARGET --cachedir $TARGET/var/cache/pacman/pkg --config $TARGET/pacman.conf -Sy

# Install base group.
pacman --root $TARGET --cachedir $TARGET/var/cache/pacman/pkg --config $TARGET/pacman.conf -S base

# Uncomment and modify if more packages are needed.
# pacman --root $TARGET --cachedir $TARGET/var/cache/pacman/pkg --config $TARGET/pacman.conf -S base-devel sudo emacs distcc

mv $TARGET/mirrorlist $TARGET/etc/pacman.d/

# Copy over key config files.
cd $TARGET/etc
cp /etc/passwd* .
cp /etc/shadow* .
cp /etc/group* .
cp /etc/sudoers .
cp /etc/resolv.conf .
cp /etc/localtime .
cp /etc/locale.gen .
cp /etc/profile.d/locale.sh profile.d
cp /etc/mtab .
cp /etc/inputrc .

echo "All done."
exit 0
