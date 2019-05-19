#!/usr/bin/env python3

import json
import random
import time
import subprocess
import sys


def gentle_blink(mempool_fee):
    print('gentle_blink({})'.format(mempool_fee))
    fee_load = min(mempool_fee, 7.5) / 7.5
    cap = 50.0
    rand = random.uniform(0.75, 1.25)
    r = int(fee_load * cap * rand + 10)
    g = int(fee_load * cap * rand + 35)
    b = int(fee_load * cap * rand + 10)
    cmd = '/opt/monitoring/blink1-tool -m 1500 --rgb={},{},{} --blink 3'.format(r, g, b)
    p = subprocess.Popen(cmd.split())
    output, error = p.communicate()
    if p.returncode != 0 or error:
        raise Exception("command failed: {}".format(error))


def tor_down_blink():
    p = subprocess.Popen('/opt/monitoring/blink1-tool -m 50 --blink 3 --rgb=75,25,25'.split())
    output, error = p.communicate()
    if p.returncode != 0 or error:
        raise Exception("command failed: {}".format(error))


def tor_up():
    p = subprocess.Popen('curl -s localhost:8335'.split(), stdout=subprocess.PIPE, stdin=subprocess.PIPE, stderr=subprocess.PIPE)
    output, error = p.communicate()
    if p.returncode != 0:
        raise Exception("command failed: {}".format(error))
    parts = output.decode().split('\n')
    status = [p for p in parts if 'tor_circuit_up' in p and p[0] != '#']
    if len(status) != 1:
        raise Exception("failed to parse tor_circuit_up: {}".format(parts))
    return status[0].split()[-1] == "1.0"


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


def main():
    tor_is_up = False
    while True:
        try:
            tor_is_up = tor_up()
        except Exception as e:
            print('Error: {}'.format(e))
            tor_is_up = False
        print("now tor_is_up: {}".format(tor_is_up))
        if tor_is_up:
            mempool_fee = bitcoin_mempool_fee()
            gentle_blink(mempool_fee)
        else:
            tor_down_blink()
        time.sleep(30)


if __name__ == '__main__':
    main()
