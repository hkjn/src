#!/usr/bin/env bash
set -euo pipefail
[[ -e /etc/bitcoin/pending-txns.txt ]] || {
	echo "No /etc/bitcoin/pending-txns.txt file to check." >&2
	exit 1
}
ANY_UNCONFIRMED=0
for tx in $(cat /etc/bitcoin/pending-txns.txt); do
	echo "checking $tx.."
	if ! sudo -u btc bitcoin-cli getrawtransaction ${tx} 1 | grep confirmations; then
	       echo "not yet confirmed"
	       ANY_UNCONFIRMED=1
	fi
	echo
done

[[ ${ANY_UNCONFIRMED} ]] && exit 1
exit 0
