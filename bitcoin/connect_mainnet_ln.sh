#!/bin/sh
#
# Connect to some Lightning Network nodes on Bitcoin mainnet.
#
set -euo pipefail

# Coco
lightning-cli connect 03939ff69d65a13c4bb2585042e7eb7e75a7c77289ab5794d1b973721d86c6839c 213.113.59.152
# dx/dy
lightning-cli connect 03d04b48cb2f277055f765c330de3c3c84f4e7d72129624bdb9b272d1113f13f14 76.10.136.25
# Blockstream store / SLEEPYARK
lightning-cli connect 02f6725f9c1c40333b67faea92fd211c183050f28df32cac3f9d69685fe9665432 104.198.32.198
# COINGAMING
lightning-cli connect 03b7ca940bc33b882dc1f1bee353a6cf205b1a7472d8ae24d45370a8f510c27d23 18.195.40.124
# #RECKLESS
lightning-cli connect 035f1498c929d4cefba4701ae36a554691f526ff60b1766badd5a49b3c8b68e1d8 78.63.23.25
# UNITEDBOUNCE
lightning-cli connect 029113620f929df927a4877ae6727b214418c366bf09ebbd4552eda235f48c00f5 68.129.210.251
# HectorJ from reddit
lightning-cli connect 0207481a19a3f51a48f134e95afa67cfeffdb38a99b5ad3494a320c4918aaaf579 163.172.174.151
