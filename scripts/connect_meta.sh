#!/bin/bash

# sshfs-mounts /meta locally, whilst allowing everyone (that have
# permissions) to access; following symlinks; attempting repairs across
# network issues.
sshfs -o reconnect,allow_other,follow_symlinks zero-one:/meta /meta

# Unmount with:
# fusermmount -u /meta

# TODO: Set up systemd service that "starts" early on, "stops" on shutdown.

 
