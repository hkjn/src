set -eu

GOPATH=${GOPATH:-""}

go get -vu github.com/mholt/caddy/caddy
go get -vu github.com/caddyserver/builds

sudo setcap "cap_net_bind_service=+ep" ${GOPATH}/bin/caddy
ulimit -n 8192 ${GOPATH}/bin/caddy
