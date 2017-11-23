#!/bin/bash
#
# Git prepush scripts.
#

source "$GOPATH/src/hkjn.me/scripts/go_git_hooks.sh" || exit

has_conflicts || exit
run_go_tests || exit
# TODO(hkjn): There's a bug where 'go vet' doesn't follow symlinks
# correctly ('go test' and others seem to work fine). Repro and report.
# From ~/src/safello/probes (canonical location /home/zero/src/github.com/Safello/probes):
# vet: error walking tree: stat ../../src/github.com/Safello/probes/sslscan/sslscan_test.go: no such file or directory
# vet: ../../src/github.com/Safello/probes/sslscan/sslscan.go: open ../../src/github.com/Safello/probes/sslscan/sslscan.go: no such file or directory
run_go_vet || exit
update_bindata || exit
update_godep || exit
needs_gofmt || exit
prevent_hacks || exit
