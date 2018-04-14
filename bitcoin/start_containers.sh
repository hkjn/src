#
# todo: rewrite to have one bitcoin and one ln container; entrypoint
# which act as tiny init system (handles signals etc.) and keeps bitcoind
# + bcmon processes alive, restarting them as necessary (supervisord)?
#
# todo: set up log forwarding (ideally distributed log aggregation / searching, oklog?) so historical
# logs for ln and other services can be viewed.
#
set -eu

IMAGE=${IMAGE:-"hkjn/bitcoin:ln-2018-04-11"}

fatal() {
	echo "FATAL: $@" >&2
	exit 1
}

[ -e /etc/lnmon/lnmon.env ] || fatal "No /etc/lnmon/lnmon.env file."

. /etc/lnmon/lnmon.env

LNMON_IP_ADDR=${LNMON_IP_ADDR:-""}
LNMON_ALIAS=${LNMON_ALIAS:-""}
LNMON_RGB=${LNMON_RGB:-"112233"}
LNMON_LOG_LEVEL=${LNMON_LOG_LEVEL:-"debug"}

[ "${LNMON_IP_ADDR}" ] || fatal "No LNMON_IP_ADDR specified."
[ "${LNMON_ALIAS}" ] || fatal "No LNMON_ALIAS specified."
[ "${LNMON_RGB}" ] || fatal "No LNMON_RGB specified."
[ "${LNMON_LOG_LEVEL}" ] || fatal "No LNMON_LOG_LEVEL specified."

DOCKER_MIN_VERSION='[1, 13, 1]'
if ! which python 1>/dev/null; then
       echo "FATAL: No python found on PATH." >&2
       exit 1
fi
PARSEVERSIONSCRIPT=$(cat <<EOF
import sys
v=[]
for x in sys.stdin.read().split('.'):
       try:
               v.append(int(x))
       except ValueError:
               v.append(int(x.split('-')[0]))
sys.exit(v >= $DOCKER_MIN_VERSION)
EOF
)
echo "Checking docker version.."
if docker version --format '{{.Server.Version}}' | python -c "${PARSEVERSIONSCRIPT}"; then
	echo "FATAL: Docker version $(docker version --format '{{.Server.Version}}') is too old; need at least ${DOCKER_MIN_VERSION}." >&2
	exit 1
fi

if ! docker network inspect bitcoin-net 1>/dev/null; then
	echo "Creating bitcoin-net network.."
	docker network create --driver=bridge --subnet=10.4.2.0/24 bitcoin-net
fi

if ! docker container inspect bitcoin 1>/dev/null; then
	echo "Starting bitcoin container.."
	docker run -d \
	           --name bitcoin \
	           -p 8333:8333 \
                   --memory 400m \
	           --network bitcoin-net \
	           --entrypoint bitcoind \
	           -v /crypt/bitcoin:/home/bitcoin/.bitcoin \
	           -v /etc/bins:/etc/bins \
	           ${IMAGE}
fi

if ! docker container inspect ln 1>/dev/null; then
	echo "Starting ln container.."
	docker run -d --name ln \
	           -p 9735:9735 \
	           --network bitcoin-net \
	           --memory 100m \
	           --entrypoint lightningd \
	            -v /crypt/bitcoin:/home/bitcoin/.bitcoin:ro \
	            -v /crypt/lightning:/home/bitcoin/.lightning \
	            ${IMAGE} \
	              --network=bitcoin \
	              --ipaddr=${LNMON_IP_ADDR} \
	              --log-level=${LNMON_LOG_LEVEL} \
	              --alias=${LNMON_ALIAS} \
	              --rgb=${LNMON_RGB}
fi

if ! docker container inspect bcmon 1>/dev/null; then
	echo "Starting bcmon container.."
	docker run -d --name bcmon \
	           -e BCMON_HTTP_PREFIX=/bcmon \
	           -p 9740:9740 \
	           --memory 50m \
	           --network bitcoin-net \
	           --pid container:bitcoin \
	           --entrypoint /etc/bins/bcmon \
	           -v /etc/bins:/etc/bins:ro \
	           -v /crypt/bitcoin:/home/bitcoin/.bitcoin:ro \
	           ${IMAGE}
fi

if ! docker container inspect lnmon 1>/dev/null; then
	echo "Starting lnmon container.."
	docker run -d --name lnmon \
	           -e LNMON_HTTP_PREFIX=/lnmon \
	           -p 8380:8380 \
	           --memory 75m \
	           --network bitcoin-net \
	           --pid container:ln \
	           --entrypoint /etc/bins/lnmon \
	           -v /etc/bins:/etc/bins:ro \
	           -v /crypt/lightning:/home/bitcoin/.lightning:ro \
	           ${IMAGE}
fi
echo "All done."

