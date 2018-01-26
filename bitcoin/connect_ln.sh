#
# Connect to known LN peers on mainnet.
#
set -euo pipefail
# SLEEPYARK / Blockstream store
lightning-cli connect 02f6725f9c1c40333b67faea92fd211c183050f28df32cac3f9d69685fe9665432 104.198.32.198
# #RECKLESS
lightning-cli connect 035f1498c929d4cefba4701ae36a554691f526ff60b1766badd5a49b3c8b68e1d8 78.63.23.25
# Bitrefill
lightning-cli connect 039514e5d704c59a0eba65d25fc5fe559a1641243ccdf80c980b1fc10ca9c30ca2 52.211.235.81
# ln.keff.org
lightning-cli connect 03ecffae58fab10791a46e89ae00cffb8260875bcdc22549d2dd79d0795e96bf00 194.71.109.91
# Tokensoft.io
lightning-cli connect 0235447c7485ff2b945bac5fbc366d54a87389bab8cacf1b64b26ec01e96bd165a 34.236.113.58:9735


# Unresponsive as of 2018-01-26, disabled:
# COINGAMING
# lightning-cli connect 03b7ca940bc33b882dc1f1bee353a6cf205b1a7472d8ae24d45370a8f510c27d23 18.195.40.124
#
# Coco
# lightning-cli connect 03939ff69d65a13c4bb2585042e7eb7e75a7c77289ab5794d1b973721d86c6839c 213.113.59.152
#
# dx/dy
# lightning-cli connect 03d04b48cb2f277055f765c330de3c3c84f4e7d72129624bdb9b272d1113f13f14 76.10.136.25
