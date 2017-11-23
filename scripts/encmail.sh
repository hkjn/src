#!/bin/bash

# Sample script for encrypting a mail to a recipient using gpg.
RECIPIENT=someone@gmail.com
gpg -a -r ${RECIPIENT} -e mail.txt
