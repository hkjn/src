#!/bin/sh
#
# Refresh the arch linux keyring for pacman.
#
pacman -Syy
pacman-key --init
pacman-key --populate archlinux
pacman-key --refresh-keys

