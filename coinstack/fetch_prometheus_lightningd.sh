#!/usr/bin/env bash

set -euo pipefail

curl -sLo prometheus-lightningd.py https://raw.githubusercontent.com/lightningd/plugins/master/prometheus/prometheus.py
sha256sum prometheus-lightningd.py /opt/lightning/prometheus-lightningd.py
