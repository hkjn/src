#!/bin/bash

# Backs up all of zero-unity's files to /media/musashi, which is encrypted.
#
# Note that files not readable by current user (notably /root,
# including /root/keys) will not be copied.
BACKUP_DIR=/media/musashi/backup/zero-unity-backup-$(date +%Y-%m-%d)
rsync -vaHAX /* ${BACKUP_DIR} --exclude={/dev/*,/proc/*,/sys/*,/tmp/*,/run/*,/mnt/*,/media/*,/meta/*,/home/${USER}/media/*,/home/${USER}/.cache/*,/home/${USER}/.secret/media/*,/home/${USER}/secret/media/*,/lost+found}

# Restore with the reverse command:
# rsync -vaHAX ${BACKUP_DIR} /* --exclude={/dev/*,/proc/*,/sys/*,/tmp/*,/run/*,/mnt/*,/media/*,meta/*,/lost+found}
