set -eu

GOPATH=${GOPATH:-""}
CADDY_VERSION=v0.10.10

mkdir -p ${GOPATH}/src/github.com/mholt/caddy
cd ${GOPATH}/src/github.com/mholt/caddy
echo "Cloning caddy at $(pwd).."
git clone https://github.com/mholt/caddy .
git checkout ${CADDY_VERSION}

go install github.com/mholt/caddy/caddy
