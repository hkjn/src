#!/usr/bin/env bash
#
# Install Bitcoin Core.
#
echo "Setting flags to exit script if any command fails or if any variable is undefined.."
set -eu

echo "Checking that commands we require are installed.."
command -v cd >/dev/null 2>&1 || { echo >&2 "cd is missing"; exit 1; }
command -v gpg >/dev/null 2>&1 || { echo >&2 "gpg is missing"; exit 1; }
command -v wget >/dev/null 2>&1 || { echo >&2 "mkdir is missing"; exit 1; }
command -v tar >/dev/null 2>&1 || { echo >&2 "cat is missing"; exit 1; }
command -v mkdir >/dev/null 2>&1 || { echo >&2 "chmod is missing"; exit 1; }
cd "${HOME}"
echo
[[ -e bitcoin-0.18.0-x86_64-linux-gnu.tar.gz ]] || {
    echo "*********************************************"
    echo "Fetching Bitcoin Core and signed hashes.."
    wget https://bitcoincore.org/bin/bitcoin-core-0.18.0/bitcoin-0.18.0-x86_64-linux-gnu.tar.gz
    wget https://bitcoincore.org/bin/bitcoin-core-0.18.0/SHA256SUMS.asc
    echo "*********************************************"
}

# if gpg servers are unstable, use local copy of key:
# wget https://j.mp/btc-release-key
# gpg --import < 01EA5486DE18A882D4C2684590C8019E36C2E964.asc
echo "*********************************************"
echo "Fetching release key.."
gpg --recv-key 01EA5486DE18A882D4C2684590C8019E36C2E964
echo "*********************************************"
echo
echo "*********************************************"
echo "Verifying GPG signature.."
gpg --verify SHA256SUMS.asc
echo "*********************************************"
echo
echo "*********************************************"
echo "Verifying SHA256 hashes.."
echo "*********************************************"
sha256sum --ignore-missing --check SHA256SUMS.asc
[[ -d "${HOME}/bitcoin-0.18.0" ]] || {
    echo "*********************************************"
    echo "Extracting archive.."
    tar xzfv bitcoin-0.18.0-x86_64-linux-gnu.tar.gz
    echo "*********************************************"
}
echo "Removing SHA256SUMS.asc file.."
rm -f SHA256SUMS.asc
echo
echo "Adding bitcoin.conf.."
mkdir -p ${HOME}/.bitcoin
cat << EOF > ~/.bitcoin/bitcoin.conf
txindex=1
server=1
disablewallet=1
printtoconsole=1
EOF
echo "*********************************************"

grep -q bitcoin ${HOME}/.bashrc || {
    echo "Adding Bitcoin Core binaries to PATH.."
    echo 'PATH=${PATH}:${HOME}/bitcoin-0.18.0/bin/:.' >> ${HOME}/.bashrc
}

echo "Bitcoin Core has been installed! Try starting it with:"
echo "  source ~/.bashrc"
echo "  bitcoind"
