set -euo pipefail

go get -v github.com/mholt/caddy/caddy
go get -v github.com/caddyserver/builds

sudo setcap "cap_net_bind_service=+ep" ~/bin/caddy
ulimit -n 8192 ${GOPATH}/bin/caddy
