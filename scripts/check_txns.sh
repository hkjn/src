#!/usr/bin/env bash
set -euo pipefail
for tx in $(cat waiting-txns.txt); do
	echo "checking $tx.."
	sudo -u btc bitcoin-cli getrawtransaction ${tx} 1 | grep confirmations || echo "not yet"
	echo
done
