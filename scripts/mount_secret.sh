#!/bin/bash

# Mount the ecryptfs filesystem in /home/zero/.secret at /home/zero/secret.
mount -t ecryptfs /home/zero/.secret/ /home/zero/secret
