#!/usr/bin/env python3


import json
import time
import random
import os



def get_mempool_info():
    data = os.popen("bitcoin-cli getmempoolinfo").read()
    # print(data)
    return json.loads(data)


def get_raw_mempool():
    data = json.loads(os.popen("bitcoin-cli getrawmempool true").read())
    fee_sum = 0.0
    for key in data:
        fee_base = data[key]['fees']['base']*10**8
        fee_sum += fee_base
    return fee_sum


def set_color(fullness):
    """Set the color of the blink1.

    Args:
        fullness: a value between 0.0 and 1.0 representing how
            full the mempool is.
    """
    amplification = 1.0
    adjusted_fullness = min(fullness * amplification, 1.0)
    cap = 64.0
    r = int(adjusted_fullness*cap)
    g = 0
    b = int(cap-adjusted_fullness*cap)
    # print('Setting color for fullness {:2.3%} ({:2.3%}): ({}, {}, {})..'.format(fullness, adjusted_fullness, r, g, b))
    os.system("blink1-tool -m 5 --rgb={},{},{}".format(r, g, b))


def block_found_blink():
    c = 127
    num_blinks = 5
    os.system("blink1-tool -m 5 --rgb {},{},{} --blink {}".format(
         c, c, c, num_blinks))


def fee_increase(fee_load, increase_fee):
    cap = 64.0
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
    cmd = "blink1-tool --rgb {},{},{} --blink {}".format(
         r, g, b, num_blinks)
    os.system(cmd)


def demo(also_blink=False):
    for f in range(100):
        print('Demoing fullness {}..'.format(f/100.0))
        set_color(f/100.0)
        time.sleep(0.05)
        if f % 10 == 0 and also_blink:
            fee_increase(f/100.0, f/100.0)
        if f % 25 == 0 and also_blink:
            block_found_blink()


def main(do_demo=False):
    if do_demo:
        demo(also_blink=True)
    
    size = 0
    fee_sum = 0.0
    while True:
        time.sleep(1.0)
        # print('\n' * 20)
        data = get_mempool_info()
        fullness = 1.0*data['bytes'] / data['maxmempool']
        new_fee_sum = get_raw_mempool() / 10.0**8
        # print('Fee sum in mempool: {:2.3} BTC'.format(new_fee_sum))
        fee_load = min(fee_sum, 5.0) / 5.00
        new_fee_load = min(new_fee_sum, 5.0) / 5.00
        if new_fee_sum > fee_sum:
            print('Sum of fees increased by {:2.3} BTC to {:2.3} BTC, new load is {:2.3%}'.format(new_fee_sum-fee_sum, new_fee_sum, new_fee_load))
            fee_increase(new_fee_load, new_fee_sum-fee_sum)
        fee_sum = new_fee_sum
        set_color(new_fee_load)
        if data['size'] < size:
            print('Size shrank! {} vs {}, a block must have been found?'.format(
                   size, data['size']))
            block_found_blink()
        size = data['size']


if __name__ == '__main__':
    main(do_demo=True)
