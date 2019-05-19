#!/usr/bin/env python3

import json
import prometheus_client
import time
import subprocess
import sys


def get_wasabi_price_btc_usd():
    args = ['curl',
            '--socks5', 'localhost:9050',
            '--socks5-hostname', 'localhost:9050',
            '-s', 'http://wasabiukrxmkdgve5kynjztuovbg43uxcbcxn6y2okcrsg7gb6jdmbad.onion/api/v3/btc/Offchain/exchange-rates',
            '-H', '"accept: application/json"']
    print('running {}'.format(' '.join(args)))
    output = subprocess.check_output(args, stderr=subprocess.STDOUT, timeout=15)
    return json.loads(output)[0]['rate']


def tor_has_circuit():
    args = 'curl --socks5 localhost:9050 --socks5-hostname localhost:9050 -s https://check.torproject.org/'.split()
    print('running {}'.format(' '.join(args)))
    output = subprocess.check_output(args, stderr=subprocess.STDOUT, timeout=15)
    return 'Congratulations' in str(output)


def main():
    tor_up_metric = prometheus_client.Gauge('tor_circuit_up', 'Whether a functional Tor circuit is up')
    wasabi_btc_usd_metric = prometheus_client.Gauge('wasabi_btc_usd', 'BTC/USD price from Wasabi')
    wasabi_sat_per_usd_metric = prometheus_client.Gauge('wasabi_sat_per_usd', 'Number of satoshi for 1 USD, price from Wasabi')
    prometheus_client.start_http_server(8335)
    has_circuit = False
    while True:
        try:
            has_circuit = tor_has_circuit()
        except Exception as e:
            print('Error: {}'.format(e))
            has_circuit = False
        has_circuit_num = 0 if not has_circuit else 1
        print("now tor_circuit_up: {}".format(has_circuit_num))
        tor_up_metric.set(has_circuit_num)
        wasabi_btc_usd_metric.set(get_wasabi_price_btc_usd())
        wasabi_sat_per_usd_metric.set(10.0**8/get_wasabi_price_btc_usd())
        time.sleep(30)


if __name__ == '__main__':
    main()
