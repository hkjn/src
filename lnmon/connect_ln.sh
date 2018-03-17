#
# Connect to known LN peers on mainnet.
#
set -euo pipefail
# ln0.decenter.world
lightning-cli connect 02d2eabbbacc7c25bbd73b39e65d28237705f7bde76f557e94fb41cb18a9ec0084 ln0.decenter.world
# Bitrefill
lightning-cli connect 024a2e265cd66066b78a788ae615acdc84b5b0dec9efac36d7ac87513015eaf6ed lnd.bitrefill.com
# TrueVision.club
lightning-cli connect 0387e3780a4325eb38421fb83000a6f6c0ffa4a69ea0c81db3f00e8e5015c9e8a1 148.251.82.174
# SLEEPYARK / Blockstream store
lightning-cli connect 02f6725f9c1c40333b67faea92fd211c183050f28df32cac3f9d69685fe9665432 104.198.32.198
# #RECKLESS
lightning-cli connect 035f1498c929d4cefba4701ae36a554691f526ff60b1766badd5a49b3c8b68e1d8 78.63.23.25
# mainnet.yalls.org
lightning-cli connect 03457d3e97da4e6a01739f56c0cd168cb14962b91fb44dc3d647816c70e05e9b93 34.200.241.1
# ln.keff.org
lightning-cli connect 03ecffae58fab10791a46e89ae00cffb8260875bcdc22549d2dd79d0795e96bf00 194.71.109.91
# Tokensoft.io
lightning-cli connect 0235447c7485ff2b945bac5fbc366d54a87389bab8cacf1b64b26ec01e96bd165a 34.236.113.58
# UNITEDBOUNCE
lightning-cli connect 029113620f929df927a4877ae6727b214418c366bf09ebbd4552eda235f48c00f5 68.129.210.251
# Cinnober
lightning-cli connect 027f3ee645c11ef3183cb6f929274baed13990e812d157866b4e41c2bf14e352de 34.240.200.0
