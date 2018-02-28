set -eu

PROMETHEUS_VERSION=2.1.0
NODE_EXPORTER_VERSION="0.15.2"
# TODO: should look up uname -m or somesuch, but need to pass amd64, not x86_64..
PLATFORM="armv7"

echo "Fetching and installing prometheus.."
curl -vLo prometheus.tar.gz https://github.com/prometheus/prometheus/releases/download/v${PROMETHEUS_VERSION}/prometheus-${PROMETHEUS_VERSION}.linux-${PLATFORM}.tar.gz
tar xzfv prometheus.tar.gz
sudo mkdir -p /etc/prometheus
sudo install prometheus-${PROMETHEUS_VERSION}.linux-${PLATFORM}/prometheus /etc/prometheus/
[ -e /etc/prometheus/prometheus.env ] || sudo touch /etc/prometheus/prometheus.env

echo "Fetching and installing node_exporter.."
curl -vLo node_exporter.tar.gz https://github.com/prometheus/node_exporter/releases/download/v${NODE_EXPORTER_VERSION}/node_exporter-${NODE_EXPORTER_VERSION}.linux-${PLATFORM}.tar.gz
tar xzfv node_exporter.tar.gz
sudo mkdir -p /etc/node_exporter
sudo install node_exporter-${NODE_EXPORTER_VERSION}.linux-${PLATFORM}/node_exporter /etc/node_exporter/
