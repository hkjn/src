#!/bin/bash

# Backs up some /media/musashi files to /mediazero-unity's files to /media/musashi, which is encrypted.
#
# Note that files not readable by current user (notably /root,
# including /root/keys) will not be copied.
BACKUP_DIR=/media/clown/backup/musashi-backup-$(date +%Y-%m-%d)
[[ -d /media/musashi ]] || {
  echo "FATAL: No /media/musashi" >&2
  exit 1
}
[[ -d /media/clown ]] || {
  echo "FATAL: No /media/clown" >&2
  exit 1
}
cd /media/musashi
rsync -vaHAX . ${BACKUP_DIR} --exclude={anime,backup,video,educational,lost+found}
