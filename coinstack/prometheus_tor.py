#!/usr/bin/env python3

import json
import prometheus_client
import time
import subprocess
import sys


# xx: move elsewhere
def get_energy():
    args = ['cat',
            '/sys/class/power_supply/BAT0/energy_now',
            '/sys/class/power_supply/BAT0/energy_full']
    print('running {}'.format(' '.join(args)))
    result = b""
    result = subprocess.check_output(args, stderr=subprocess.STDOUT, timeout=15)
    return [int(x) for x in result.split()]


def get_connection_status(host):
    args = ['curl',
            '-sL',
            '--output', '/dev/null',
            host]
    print('running {}'.format(' '.join(args)))
    try:
        subprocess.check_call(args, stderr=subprocess.STDOUT, timeout=15)
    except subprocess.CalledProcessError:
        return False
    return True


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


def tor_can_reach_frankenbox():
    args = ('curl --socks5 localhost:9050 --socks5-hostname localhost:9050 '
            '-s didqi7crkvxxqz5pnw2khrac7ncchwqgz5yyxcfvbow56wh3fce2peid.onion:22').split()
    print('running {}'.format(' '.join(args)))
    output = subprocess.check_output(args, stderr=subprocess.STDOUT, timeout=15)
    return 'SSH' in str(output)


def main():
    print('prometheus_tor.py starting..')
    network_up_metric = prometheus_client.Gauge('network_up', 'Whether a functional network is up')
    shift_network_up_metric = prometheus_client.Gauge('shift_network_up', 'Whether a functional network is up that can reach rocketchat.shiftcrypto.ch')
    tor_up_metric = prometheus_client.Gauge('tor_circuit_up', 'Whether a functional Tor circuit is up')
    frankenbox_is_alive_metric = prometheus_client.Gauge('frankenbox_is_alive', 'Whether the frankenbox is aliveo')
    wasabi_btc_usd_metric = prometheus_client.Gauge('wasabi_btc_usd', 'BTC/USD price from Wasabi')
    wasabi_sat_per_usd_metric = prometheus_client.Gauge('wasabi_sat_per_usd', 'Number of satoshi for 1 USD, price from Wasabi')
    energy_now_metric = prometheus_client.Gauge('energy_now', 'Current energy level')
    energy_full_metric = prometheus_client.Gauge('energy_full', 'Fully charged energy level')

    prometheus_client.start_http_server(8335)
    while True:
        has_network = False
        has_circuit = False
        frankenbox_is_alive = False
        has_network = get_connection_status("github.com")
        has_network_shift = get_connection_status("rocketchat.shiftcrypto.ch")
        try:
            has_circuit = tor_has_circuit()
        except (subprocess.CalledProcessError, subprocess.TimeoutExpired) as e:
            print('Error: {}'.format(e))
            has_circuit = False
        try:
            has_circuit = tor_has_circuit()
        except (subprocess.CalledProcessError, subprocess.TimeoutExpired) as e:
            print('Error: {}'.format(e))
            has_circuit = False
        try:
            frankenbox_is_alive = tor_can_reach_frankenbox()
        except (subprocess.CalledProcessError, subprocess.TimeoutExpired) as e:
            print('Error: {}'.format(e))
            frankenbox_is_alive = False

        has_network_num = 0 if not has_network else 1
        has_network_shift_num = 0 if not has_network_shift else 1
        print("network_up: {}".format(has_network_num))
        print("shift_network_up: {}".format(has_network_shift_num))

        has_circuit_num = 0 if not has_circuit else 1
        frankenbox_is_alive_num = 0 if not frankenbox_is_alive else 1
        print("tor_circuit_up: {}".format(has_circuit_num))
        print("frankenbox_is_alive_num: {}".format(frankenbox_is_alive_num))

        energy_now, energy_full = get_energy()
        print("energy_now: {}, energy_full: {}".format(energy_now, energy_full))
        energy_now_metric.set(energy_now)
        energy_full_metric.set(energy_full)
        network_up_metric.set(has_network_num)
        shift_network_up_metric.set(has_network_shift_num)
        tor_up_metric.set(has_circuit_num)
        frankenbox_is_alive_metric.set(frankenbox_is_alive_num)
        wasabi_price_btc_usd = None
        try:
            wasabi_price_btc_usd = get_wasabi_price_btc_usd()
        except (subprocess.CalledProcessError, subprocess.TimeoutExpired) as e:
            print('Error getting wasabi price: {}'.format(e))
        if wasabi_price_btc_usd:
            wasabi_btc_usd_metric.set(wasabi_price_btc_usd)
            wasabi_sat_per_usd_metric.set(10.0**8/wasabi_price_btc_usd)

        time.sleep(30)


if __name__ == '__main__':
    main()
