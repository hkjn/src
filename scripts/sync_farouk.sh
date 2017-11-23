#!/bin/bash

# Syncs /media/farouk (internal encrypted drive) into /media/clown
# (external encrypted drive).
#
# This script requires root permissions to sync gitrepos (but will
# sync the rest with regular permissions).

set -e

# TODO: Run this in cron (30m?).
rsync --exclude 'lost+found' -avz /media/farouk /media/clown/
echo "Done."


