#!/bin/bash
#
# Removes all of /usr/share/locale *except* the specified ones.
#
[ "$UID" == 0 ] || {
	echo "This script requires root privileges." >&2
}
mkdir /tmp/locales/
for k in en en\@boldquot en\@quot en\@shaw en_US; do
		cd /usr/share/locale
		mv -v $k /tmp/locales/
done
sudo rm -rfv /usr/share/locale/*
sudo mv -v /tmp/locales/* /usr/share/locale/
