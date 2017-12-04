# This tool looks at bootstrap/ .json configs to verify that the specified URLs
# can be fetched and have matching checksums.
#
# TODO: use "download binaries and verify checksums"
# logic here to fetch" CA cert, client cert + key,
# and then call telemetry/install_client:
#   curl -o /etc/ssl/mon_ca.pem https://admin1.hkjn.me/151f[...]/files/hkjninfra/1.5.5/certs/mon_ca.pem
#
import json
import logging
import os
import subprocess
import sys
import tempfile
import urllib2


def fetch(url, path):
        response = urllib2.urlopen(url)
        html = response.read()
        with open(path, 'w+') as tmpfile:
            tmpfile.write(html)
            print('  Wrote {}.'.format(path))


def verify_digest(path, digest):
        child = subprocess.Popen(["sha512sum", path], stdout=subprocess.PIPE)
        failed = child.wait()
        if failed:
            raise RuntimeError('sha512sum {} call failed.'.format(path))
        shaparts = child.stdout.read().split()
        if shaparts[0] != digest:
            raise RuntimeError('bad checksum for {}: {} vs {}'.format(
                path,
                shaparts[0][:5],
                digest[:5],
            ))

def run():
    if len(sys.argv) != 2:
        raise RuntimeError('Usage: {} node'.format(sys.argv[0]))
    dryrun = os.environ.get('VALIDATE_DRYRUN')

    data=json.loads(open("bootstrap/{}.json".format(sys.argv[1])).read())
    for f in data['storage']['files']:
        if 'source' not in f['contents'] or 'http' not in f['contents']['source']:
            continue
        url = f['contents']['source']
        digest = f['contents']['verification']['hash'].split('sha512-')[1]
        if dryrun:
            print('  URL {} should have checksum {}'.format(url, digest))
        else:
            print('Verifying checksum of {}..'.format(f['path']))
            print('  Fetching {}, which should have checksum {}..'.format(url, digest[:5]))
            path = tempfile.mkdtemp(prefix='hkjninfra_checksums')
            tmppath = os.path.join(path, digest)
            fetch(url, tmppath)
            verify_digest(tmppath, digest)
            print('  Checksum matches!')
    if not dryrun:
        print('All checksums matched.')


if __name__ == '__main__':
    run()
