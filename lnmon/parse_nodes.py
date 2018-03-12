import json
import socket

d = None
with open('nodes.json') as json_file:
    d = json.loads(json_file.read())

for n in d['nodes']:
    if 'addresses' not in n or not n['addresses']:
        continue

    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.settimeout(5.0)
    addr = n['addresses'][0]['address']
    port = n['addresses'][0]['port']
    nd = n['nodeid']
    if 'alias' in n:
        nd = n['alias']
    desc = '{} ({}:{})'.format(nd, addr, port)
    try:
        # print('Connecting to {}'.format(nd))
        # print('s.connect(({}, {}))'.format(addr, port))
        s.connect((addr, port))
        print('lightning-cli connect {} {}:{}'.format(n['nodeid'], addr, port))
    except socket.error:
        pass
        # print('Failed to connect to {}'.format(nd))
