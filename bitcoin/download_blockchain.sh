#
# Sample script for downloading checksummsed bitcoin blocks and chainstate from S3-compatible storage.
#
set -eu
BLOCK=508652

[ -e /etc/secrets/digitalocean/digitalocean0_gpg_pass ] || {
	echo "FATAL: No /etc/secrets/digitalocean/digitalocean0_gpg_pass." >&2
	exit 1
}
[ -e /etc/secrets/digitalocean/digitalocean0_spaces_key ] || {
	echo "FATAL: No /etc/secrets/digitalocean/digitalocean0_spaces_key." >&2
	exit 1
}
[ -e /etc/secrets/digitalocean/digitalocean0_spaces_secret ] || {
	echo "FATAL: No /etc/secrets/digitalocean/digitalocean0_spaces_secret." >&2
	exit 1
}

echo "Downloading blocks up to ${BLOCK}.."
docker run --name bitcoin-dl -d -v /crypt:/crypt \
            -e S3_GPG_PASS=$(cat /etc/secrets/digitalocean/digitalocean0_gpg_pass) \
            -e S3_ENDPOINT=nyc3.digitaloceanspaces.com \
            -e S3_KEY=$(cat /etc/secrets/digitalocean/digitalocean0_spaces_key) \
            -e S3_SECRET=$(cat /etc/secrets/digitalocean/digitalocean0_spaces_secret) \
  hkjn/s3cmd:1.0.0 sync s3://zdo/bitcoin-${BLOCK} /crypt -vvv

cd /crypt/bitcoin-${BLOCK}
echo "Verifying checksums.."
sha256sum -c ${BLOCK}_checksums.SHASUMS
echo "Moving to /crypt/bitcoin.."
if [[ -e /crypt/bitcoin ]]; then
	echo "FATAL: Refusing to clobber existing /crypt/bitcoin directory." >&2
	exit 1
fi
mv /crypt/bitcoin-${BLOCK} /crypt/bitcoin
