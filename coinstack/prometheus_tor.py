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
    result = subprocess.check_output(args, stderr=subprocess.STDOUT, timeout=15)
    return [int(x) for x in result.split()]


def get_pending_txns():
    args = ['ls',
            '/etc/bitcoin/pending-txns/']
    print('running {}'.format(' '.join(args)))
    result = subprocess.check_output(args, stderr=subprocess.STDOUT, timeout=15)
    return [str(x) for x in result.split()]


def get_number_threads():
    args = ['ps',
            '-AL',
            '--no-headers']
    print('running "{}"'.format(' '.join(args)))
    result = subprocess.check_output(args, stderr=subprocess.STDOUT, timeout=15)
    return len(result.split(b'\n'))


def get_connection_status(addr):
    args = ['curl',
            '-sL',
            '--output', '/dev/null',
            addr]
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
    args = ('torsocks nc -w 5 '
            'didqi7crkvxxqz5pnw2khrac7ncchwqgz5yyxcfvbow56wh3fce2peid.onion 8527').split()
    print('running {}'.format(' '.join(args)))
    output = subprocess.check_output(args, stderr=subprocess.STDOUT, timeout=15)
    return 'SSH' in str(output)


def tor_can_reach_spark():
    args = ('curl --socks5 localhost:9050 --socks5-hostname localhost:9050 '
            '-s http://zfzt7mmfqaj263jp33hacwpn55l2nyqj4pgdg5ktdgxd2jy2zd4p3aad.onion').split()
    print('running {}'.format(' '.join(args)))
    output = subprocess.check_output(args, stderr=subprocess.STDOUT, timeout=15)
    return '401 Unauthorized' in str(output)


def main():
    print('prometheus_tor.py starting..')
    network_up_metric = prometheus_client.Gauge('github_reached_count', 'Number of times a HTTPS connection to github.com could be made')
    tor_up_metric = prometheus_client.Gauge('tor_status_reached_count', 'Number of times the Tor status page could be reached via Tor')
    frankenbox_is_alive_metric = prometheus_client.Gauge('frankenbox_sshd_onion_reached_count', 'Number of times the frankenbox sshd onion service was reached')
    wasabi_btc_usd_metric = prometheus_client.Gauge('wasabi_btc_usd', 'BTC/USD price from Wasabi')
    tor_can_reach_wasabi_metric = prometheus_client.Gauge('wasabi_onion_reached_count', 'Number of times the Wasabi .onion service was reached')
    tor_can_reach_spark_metric = prometheus_client.Gauge('spark_onion_reached_count', 'Number of times Spark .onion service was reached')
    wasabi_sat_per_usd_metric = prometheus_client.Gauge('wasabi_sat_per_usd', 'Number of satoshi for 1 USD, price from Wasabi')
    num_pending_txns_metric = prometheus_client.Gauge('num_pending_txns', 'Number of pending unconfirmed txns')
    num_threads_metric = prometheus_client.Gauge('num_threads', 'Number of threads')
    energy_now_metric = prometheus_client.Gauge('energy_now', 'Current energy level')
    energy_full_metric = prometheus_client.Gauge('energy_full', 'Fully charged energy level')

    prometheus_client.start_http_server(8335)
    while True:
        has_network = False
        try:
            has_network = get_connection_status("https://github.com")
        except (subprocess.CalledProcessError, subprocess.TimeoutExpired) as e:
            print('Error: {}'.format(e))
        print("network_up: {}".format(has_network))
        if has_network:
            network_up_metric.inc()

        has_circuit = False
        try:
            has_circuit = tor_has_circuit()
        except (subprocess.CalledProcessError, subprocess.TimeoutExpired) as e:
            print('Error: {}'.format(e))
        print("tor_circuit_up: {}".format(has_circuit))
        if has_circuit:
            tor_up_metric.inc()

        frankenbox_is_alive = False
        try:
            frankenbox_is_alive = tor_can_reach_frankenbox()
        except (subprocess.CalledProcessError, subprocess.TimeoutExpired) as e:
            print('Error: {}'.format(e))
        print("frankenbox_is_alive_num: {}".format(frankenbox_is_alive))
        if frankenbox_is_alive:
            frankenbox_is_alive_metric.inc()

        num_pending_txns_metric.set(len(get_pending_txns()))
        energy_now, energy_full = get_energy()
        print("energy_now: {}, energy_full: {}".format(energy_now, energy_full))
        energy_now_metric.set(energy_now)
        energy_full_metric.set(energy_full)

        num_threads_metric.set(get_number_threads())

        wasabi_price_btc_usd = None
        tor_can_reach_wasabi = 0
        try:
            wasabi_price_btc_usd = get_wasabi_price_btc_usd()
            tor_can_reach_wasabi_metric.inc()
        except (subprocess.CalledProcessError, subprocess.TimeoutExpired) as e:
            print('Error getting wasabi price: {}'.format(e))
        if wasabi_price_btc_usd:
            wasabi_btc_usd_metric.set(wasabi_price_btc_usd)
            wasabi_sat_per_usd_metric.set(10.0**8/wasabi_price_btc_usd)

        can_reach_spark = False
        try:
            can_reach_spark = tor_can_reach_spark()
        except (subprocess.CalledProcessError, subprocess.TimeoutExpired) as e:
            print('Error reaching spark .onion: {}'.format(e))
        if tor_can_reach_spark:
            tor_can_reach_spark_metric.inc()

        time.sleep(30)


if __name__ == '__main__':
    main()
