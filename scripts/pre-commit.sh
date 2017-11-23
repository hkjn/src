#!/bin/bash
#
# Git precommit scripts.

source "$GOPATH/src/hkjn.me/scripts/go_git_hooks.sh" || exit

needs_gofmt || exit
