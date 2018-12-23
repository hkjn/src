#!/bin/bash
#
# Terminate all keybase processes.
#
set -euo pipefail
systemctl --user stop keybase kbfs keybase.gui
sudo fusermount -u /keybase
