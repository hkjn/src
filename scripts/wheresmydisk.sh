#!/bin/bash
#
# Checks where disk usage went using ncdu, excluding some mounted
# partitions.
#
sudo ncdu / --exclude /meta --exclude /run/media --exclude /media --exclude /home
