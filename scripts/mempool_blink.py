#!/usr/bin/env python3


import json
import time
import random
import os


def get_mempool_info():
    data = os.popen("bitcoin-cli -datadir=/bitcoin/.bitcoin getmempoolinfo").read()
    # print(data)
    json_data = None
    try:
        json_data = json.loads(data)
    except json.decoder.JSONDecodeError as e:
        print('Failed to fetch mempool info: {}'.format(e))
        raise Exception(e)
    return json_data


def get_raw_mempool():
    data = json.loads(os.popen("bitcoin-cli -datadir=/bitcoin/.bitcoin getrawmempool true").read())
    fee_sum = 0.0
    for key in data:
        fee_base = data[key]['fees']['base']*10**8
        fee_sum += fee_base
    return fee_sum


def set_color(fullness, in_ibd=False):
    """Set the color of the blink1.

    Args:
        fullness: a value between 0.0 and 1.0 representing how
            full the mempool is.
        in_ibd: a bool that if True means we are in initial block
            download (IBD).
    """
    amplification = 1.0
    adjusted_fullness = min(fullness * amplification, 1.0)
    cap = 50.0
    r = int(adjusted_fullness*cap)
    g = 0 if not in_ibd else 50
    b = int(cap-adjusted_fullness*cap)
    os.system("blink1-tool -m 10 --rgb={},{},{}".format(r, g, b))


def block_found_blink():
    c = 127
    num_blinks = 5
    os.system("blink1-tool -m 500 --rgb {},{},{} --blink {}".format(
         c, c, c, num_blinks))


def fee_increase_blink():
    c = 64
    num_blinks = 2
    os.system("blink1-tool -m 500 --rgb {},{},{} --blink {}".format(
         c, c, 0, num_blinks))


def fee_increase(fee_load, increase_fee):
    cap = 50.0
    multiplier = 1.0 + increase_fee
    num_blinks = 1
    if increase_fee > 0.01:
        num_blinks += 1
    if increase_fee > 0.1:
        num_blinks += 1
        multiplier *= 1.5
    adjusted_load = min(fee_load * multiplier, 1.0)
    r = int(adjusted_load*cap)
    g = 0
    b = int(cap-adjusted_load*cap)
    cmd = "blink1-tool -m 10 --rgb {},{},{} --blink {}".format(
         r, g, b, num_blinks)
    os.system(cmd)


def demo(also_blink=False):
    for f in range(100):
        print('Demoing fullness {}..'.format(f/100.0))
        set_color(f/100.0, in_ibd=True if f < 70 else False)
        time.sleep(0.05)
        if f % 10 == 0 and also_blink:
            fee_increase(f/100.0, f/100.0)
        if f % 25 == 0 and also_blink:
            block_found_blink()


def is_in_ibd():
    data = None
    with os.popen('bitcoin-cli -datadir=/bitcoin/.bitcoin getblockchaininfo') as output:
        data=output.read()
    if not data:
        raise Exception('no output from bitcoin-cli')
    json_data = json.loads(data)
    in_ibd = json_data['initialblockdownload']
    return in_ibd


def main(do_demo=False):
    if do_demo:
        demo(also_blink=True)
    
    size = 0
    fee_sum = 0.0
    next_fee_increase_blink_sats = None
    while True:
        time.sleep(2)
        try:
            data = get_mempool_info()
        except Exception as e:
            print('Got exception: {}..'.format(e))
            continue

        in_ibd = is_in_ibd()
        fee_sats = get_raw_mempool()
        fee_load = min(fee_sats, 7.5*10**8) / (7.5*10**8)
        if not next_fee_increase_blink_sats:
            next_fee_increase_blink_sats = fee_sats

        print('')
        print('[IBD={}] Fee sum in mempool: {:2.3} BTC'.format(in_ibd, fee_sats/10.0**8))
        if fee_sats > next_fee_increase_blink_sats:
            fee_increase_blink()
            next_fee_increase_blink_sats += 0.1 * 10.0**8
            print('Blinking next at {:2.3} BTC'.format(next_fee_increase_blink_sats/10.0**8))

        if data['size'] < size:
            print('Size shrank! {} vs {}, a block must have been found?'.format(
                   size, data['size']))
            block_found_blink()
        size = data['size']

        set_color(fee_load, in_ibd=in_ibd)


if __name__ == '__main__':
    main(do_demo=False)
