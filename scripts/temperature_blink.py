#!/usr/bin/env python3


import time
import random
import os


def read_avg_temp():
    # xx: could average across more than one sample since temperature readings
    # seem to jump around a lot
    temps = os.popen("cat /sys/class/thermal/thermal_zone?/temp").read().split()
    return sum(int(t) for t in temps) / len(temps)


def read_max_temp():
    temps = os.popen("cat /sys/class/thermal/thermal_zone?/temp").read().split()
    return max(int(t) for t in temps)


def set_color(temp):
    temp = min(temp, 100000.0)
    r = 164*(temp / 100000.0)
    r = int(random.uniform(0.8, 1.2) * r)
    os.system("blink1-tool -m 5 --rgb={},0,25".format(r))


def main(do_demo=False):
    if do_demo:
        for t in range(100):
            set_color(t*1000.0)
            time.sleep(0.1)
    while True:
        time.sleep(0.5 * random.uniform(0.8, 1.2))
        t = read_max_temp()
        print('Read temperature, max is {} C'.format(t/1000.0))
        set_color(t)


if __name__ == '__main__':
    main(do_demo=False)
