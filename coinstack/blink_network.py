#!/usr/bin/env python3

import json
import time
import subprocess
import sys


def blink(r=50, g=50, b=50, m=15, times=3):
    cmd = '/opt/bitcoin/blink1-tool -m {} --rgb={},{},{} --blink {}'.format(m, r, g, b, times)
    p = subprocess.Popen(cmd.split())
    _, error = p.communicate()
    if p.returncode != 0 or error:
        raise Exception("command {} failed: {}".format(cmd, error))


def gentle_blink(mempool_fee):
    print('gentle_blink({})'.format(mempool_fee))
    fee_load = min(mempool_fee, 7.5) / 7.5
    cap = 50.0
    r = int(fee_load * cap)
    g = int(fee_load * cap + 17)
    b = int(fee_load * cap + 17)
    blink(r, g, b, times=1, m=10)


def get_pending_txns():
    cmd = 'ls /etc/bitcoin/pending-txns/'
    p = subprocess.Popen(cmd.split(), stdout=subprocess.PIPE, stdin=subprocess.PIPE, stderr=subprocess.PIPE)
    output, error = p.communicate()
    if p.returncode != 0 or error:
        raise Exception("command failed: {}".format(error))
    return [line for line in output.decode().split('\n') if line]


def blink_pending_txns():
    pending = get_pending_txns()
    for txid in pending:
        print('xx: checking {}'.format(txid))
        cmd = 'bitcoin-cli getrawtransaction {} 1'.format(txid)
        try:
            output = subprocess.check_output(cmd.split(), stderr=subprocess.STDOUT, timeout=10)
        except Exception as e:
            print("command failed: {}".format(e))
            blink(r=100, times=3)
            continue

        print('xx: getrawtransaction finished')
        rawtx = json.loads(output)
        confirmed = rawtx.get("confirmations", 0) > 0
        if confirmed:
            print('xx: {} confirmed! moving file to confirmed..'.format(txid))
            cmd = 'mv /etc/bitcoin/pending-txns/{} /etc/bitcoin/confirmed-txns/{}'.format(txid, txid)
            p = subprocess.Popen(cmd.split(), stderr=subprocess.PIPE)
            _, error = p.communicate()
            if p.returncode != 0 or error:
                raise Exception("command failed: {}".format(error))
            blink(g=100, m=50, times=10)
        else:
            print('xx: not yet confirmed: {}'.format(txid))


def tor_down_blink():
    p = subprocess.Popen('/opt/bitcoin/blink1-tool -m 50 --blink 3 --rgb=75,25,25'.split())
    _, error = p.communicate()
    if p.returncode != 0 or error:
        raise Exception("command failed: {}".format(error))

def tor_has_circuit():
    args = 'curl --socks5 localhost:9050 --socks5-hostname localhost:9050 -s https://check.torproject.org/'.split()
    print('running {}'.format(' '.join(args)))
    output = subprocess.check_output(args, stderr=subprocess.STDOUT, timeout=15)
    return 'Congratulations' in str(output)


def bitcoin_mempool_fee():
    p = subprocess.Popen('curl -s localhost:8334'.split(), stdout=subprocess.PIPE, stdin=subprocess.PIPE, stderr=subprocess.PIPE)
    output, error = p.communicate()
    if p.returncode != 0:
        raise Exception("command failed: {}".format(error))
    parts = output.decode().split('\n')
    status = [p for p in parts if 'bitcoin_mempool_fee_sum' in p and p[0] != '#']
    if len(status) != 1:
        raise Exception("failed to parse bitcoin_mempool_fee_sum: {}".format(parts))
    return float(status[0].split()[-1])


def get_light():
    with open('/opt/hack/light') as f:
        return int(f.read())


def main():
    print('blink_network.py starting..')
    tor_is_up = False
    while True:
        time.sleep(30)
        # Only blink if the light setting is high enough, so we can turn down
        # the light and the blinking doesn't annoy anyone at night.
        should_blink = get_light() > 3
        if not should_blink:
            continue
        try:
            blink_pending_txns()
        except Exception as e:
            print('Error: {}'.format(e))
        try:
            tor_is_up = tor_has_circuit()
        except Exception as e:
            print('Error: {}'.format(e))
            tor_is_up = False
        print("now tor_is_up: {}".format(tor_is_up))
        if tor_is_up:
            mempool_fee = bitcoin_mempool_fee()
            gentle_blink(mempool_fee)
        else:
            tor_down_blink()


if __name__ == '__main__':
    main()
