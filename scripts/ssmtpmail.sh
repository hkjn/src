#!/bin/bash

# Example of how to send mail using SSMTP:
# https://wiki.archlinux.org/index.php/SSMTP
cat mail.txt | mail -v -s "Heading" tousername@somedomain.com
