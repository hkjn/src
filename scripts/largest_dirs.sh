#!/bin/sh

# Largest 15 directories in $1. Output is in GBs.
#
# Example usage:
# $0 .

du --block-size=1G $1 | sort -n -r | head -n 15
