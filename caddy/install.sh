set -eu

GOPATH=${GOPATH:-""}
CADDY_VERSION=v0.10.10

mkdir -p ${GOPATH}/src/github.com/mholt/caddy
cd ${GOPATH}/src/github.com/mholt/caddy
echo "Cloning caddy at $(pwd).."
git clone https://github.com/mholt/caddy .
git checkout ${CADDY_VERSION}

#mkdir -p ${GOPATH}/src/github.com/caddyserver
#cd ${GOPATH}/src/github.com/caddyserver
#git clone https://github.com/caddyserver/builds

go install github.com/mholt/caddy/caddy
# go install github.com/caddyserver/builds
# go get -v -u github.com/mholt/caddy/caddy
# go get -v -u github.com/caddyserver/builds

