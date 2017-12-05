#
# Sample script for uploading checksummsed bitcoin blocks and chainstate to S3-compatible storage.
#
set -euo pipefail
declare BLOCK=497220
mkdir /crypt/uploading
mv /crypt/bitcoin/{blocks,chainstate} /crypt/uploading/

echo "Creating checksums file.."
cd /crypt/uploading
find . -type f -print0 | xargs -0 sha1sum > ${BLOCK}_checksums.SHASUMS

echo "Uploading /crypt/uploading.."
docker run -d -it --name bitcoin-upload -v /crypt:/crypt \
           -e S3_GPG_PASS=$(cat /etc/secrets/digitalocean/digitalocean0_gpg_pass) \
           -e S3_ENDPOINT=nyc3.digitaloceanspaces.com \
           -e S3_KEY=$(cat /etc/secrets/digitalocean/digitalocean0_spaces_key) \
           -e S3_SECRET=$(cat /etc/secrets/digitalocean/digitalocean0_spaces_secret) \
    hkjn/s3cmd:1.0.0 sync /crypt/uploading/ s3://zdo/bitcoin-${BLOCK} -vvv
echo "Moving back /crypt/uploading.."
mv -v /crypt/uploading/{blocks,chainstate} /crypt/bitcoin/
