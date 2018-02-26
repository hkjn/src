set -eu

IMAGE=${IMAGE:-"hkjn/bitcoin:lightning-2018-02-19-amd64"}

fatal() {
	echo "FATAL: $@" >&2
	exit 1
}

[ ! -e /etc/lnmon/lnmon.env ] || fatal "No /etc/lnmon/lnmon.env file."

source /etc/lnmon/lnmon.env

[ "${LNMON_IP_ADDR}" ] || fatal "No LNMON_IP_ADDR specified."
[ "${LNMON_ALIAS}" ] || fatal "No LNMON_ALIAS specified."
[ "${LNMON_RGB}" ] || fatal "No LNMON_RGB specified."
[ "${LNMON_LOG_LEVEL}" ] || fatal "No LNMON_LOG_LEVEL specified."

if ! docker network inspect bitcoin-net 1>/dev/null; then
	echo "Creating bitcoin-net network.."
	docker network create --driver=bridge --subnet=10.4.2.0/24 bitcoin-net
fi

if ! docker container inspect bitcoin 1>/dev/null; then
	echo "Starting bitcoin container.."
	docker run -d \
	           --name bitcoin \
	           -p 8333:8333 \
	           --network bitcoin-net \
	           --entrypoint bash \
	           -v /crypt/bitcoin:/home/bitcoin/.bitcoin \
	           -v /etc/bins:/etc/bins \
	           ${IMAGE} \
	           -c "bitcoind -dbcache=1200 -onlynet=ipv4 -printtoconsole"
fi

if ! docker container inspect ln 1>/dev/null; then
	echo "Starting ln container.."
	docker run -d --name ln \
	           -p 9735:9735 \
	           --network bitcoin-net \
	           --entrypoint lightningd \
	            -v /crypt/bitcoin:/home/bitcoin/.bitcoin:ro \
	            -v /crypt/lightning:/home/bitcoin/.lightning \
	            ${IMAGE} \
	              --network=bitcoin \
	              --ipaddr=${IP_ADDR} \
	              --log-level=${LOG_LEVEL} \
	              --alias=${ALIAS} \
	              --rgb=${RGB}
fi

if ! docker container inspect bcmon 1>/dev/null; then
	echo "Starting bcmon container.."
	docker run -d --name bcmon \
	           -e BCMON_HTTP_PREFIX=/bcmon \
	           -p 9740:9740 \
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
	           --network bitcoin-net \
	           --pid container:ln \
	           --entrypoint /etc/bins/lnmon \
	           -v /etc/bins:/etc/bins:ro \
	           -v /crypt/lightning:/home/bitcoin/.lightning:ro \
	           ${IMAGE}
fi

