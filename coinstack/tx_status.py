#!/usr/bin/env python3
#
# Check the confirmation status of transactions of interest.
#
# All files under /etc/bitcoin/pending-txns/<txid> will be checked, and
# if it's confirmed, the marker file will be moved to /etc/bitcoin/confirmed-txns.
#
import json
import subprocess


def get_pending_txns():
    cmd = 'ls /etc/bitcoin/pending-txns/'
    p = subprocess.Popen(cmd.split(), stdout=subprocess.PIPE, stdin=subprocess.PIPE, stderr=subprocess.PIPE)
    output, error = p.communicate()
    if p.returncode != 0 or error:
        raise Exception("command failed: {}".format(error))
    return [line for line in output.decode().split('\n') if line]


def move_confirmed_txns():
    pending = get_pending_txns()
    for txid in pending:
        print('xx: checking {}'.format(txid))
        cmd = 'bitcoin-cli getrawtransaction {} 1'.format(txid)
        try:
            output = subprocess.check_output(cmd.split(), stderr=subprocess.STDOUT, timeout=10)
        except Exception as e:
            print("command failed: {}".format(e))
            continue

        rawtx = json.loads(output)
        confirmed = rawtx.get("confirmations", 0) > 0
        if confirmed:
            print('xx: {} confirmed! moving file to confirmed..'.format(txid))
            cmd = 'mv /etc/bitcoin/pending-txns/{} /etc/bitcoin/confirmed-txns/{}'.format(txid, txid)
            p = subprocess.Popen(cmd.split(), stderr=subprocess.PIPE)
            _, error = p.communicate()
            if p.returncode != 0 or error:
                raise Exception("command failed: {}".format(error))
            # xx: inc prometheus metric here
        else:
            print('not yet confirmed: {}'.format(txid))


if __name__ == '__main__':
    move_confirmed_txns()
