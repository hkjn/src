#!/bin/bash

# Backs up all of zero-one's files to /media/musashi.
#
# Note that files not readable by current user (notably /root,
# including /root/keys) will not be copied.
BACKUP_DIR=/media/musashi/backup/zero-one-backup-$(date +%Y-%m-%d)
rsync -vaHAX zero-one:/* ${BACKUP_DIR} --exclude={/dev/*,/proc/*,/sys/*,/tmp/*,/run/*,/mnt/*,/media/*,/lost+found}

# Restore with the reverse command:
# rsync -vaHAX ${BACKUP_DIR} zero-one/* --exclude={/dev/*,/proc/*,/sys/*,/tmp/*,/run/*,/mnt/*,/media/*,/lost+found}
