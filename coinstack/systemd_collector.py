#!/usr/bin/env python3
#
# Provide information about systemd services to Prometheus.
#

import json
import time
import subprocess
import sys
import prometheus_client

# the cpu temp doesn't really belong here as it's not related to
# systemd services, but the /sys/class/hwmon path used by the
# node_collector hwmon collector is not present on arm64..
CPU_TEMP = prometheus_client.Gauge("cpu_temp", "CPU temperature")

## systemd metrics increase every time the unit is found to be active
SYSTEMD_BITCOIND = prometheus_client.Gauge("systemd_bitcoind_active_count", "Systemd unit active count for Bitcoin Core")
SYSTEMD_ELECTRS = prometheus_client.Gauge("systemd_electrs_active_count", "Systemd unit active count for Electrs")
SYSTEMD_LND = prometheus_client.Gauge("systemd_lnd_active_count", "Systemd unit status active count for lnd")
SYSTEMD_PROMETHEUS = prometheus_client.Gauge("systemd_prometheus_active_count", "Systemd unit status active count for Prometheus")
SYSTEMD_GRAFANA = prometheus_client.Gauge("systemd_grafana_active_count", "Systemd unit status active count for Grafana")


def readFile(filepath):
    args = [filepath]
    with open(filepath) as f:
        value = f.readline()
    return value


def systemd_unit_running(unit):
    cmd = ["systemctl", "is-active", unit]
    exit_status = 0
    extra_info = ""
    try:
        subprocess.check_output(cmd)
    except subprocess.CalledProcessError as e:
        exit_status = e.returncode
        extra_info = e.output.strip()
    if exit_status != 0:
        print('command "{}" exited {}: "{}"'.format(' '.join(cmd), exit_status, extra_info))
    return exit_status == 0

def main():
    addr = 8400
    print('Serving systemd status at {}..'.format(addr))
    prometheus_client.start_http_server(addr)
    while True:
        CPU_TEMP.set(readFile("/sys/class/thermal/thermal_zone0/temp"))
        if systemd_unit_running('bitcoind'):
            SYSTEMD_BITCOIND.inc()
        if systemd_unit_running('electrs'):
            SYSTEMD_ELECTRS.inc()
        if systemd_unit_running('lnd'):
            SYSTEMD_LND.inc()
        if systemd_unit_running('prometheus'):
            SYSTEMD_PROMETHEUS.inc()
        if systemd_unit_running('grafana-server'):
            SYSTEMD_GRAFANA.inc()
        time.sleep(10)


if __name__ == "__main__":
    main()
