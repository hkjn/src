#!/usr/bin/env bash

set -euo pipefail

curl -sLO https://raw.githubusercontent.com/digitalbitbox/bitbox-base/master/armbian/base/scripts/prometheus-bitcoind.py
sha256sum prometheus-bitcoind.py /opt/bitcoin/prometheus-bitcoind.py
