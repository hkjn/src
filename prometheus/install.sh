set -eu
declare NODE_EXPORTER_VERSION="0.15.2"
declare PLATFORM="armv7"
curl -vLo node_exporter.tar.gz https://github.com/prometheus/node_exporter/releases/download/v${NODE_EXPORTER_VERSION}/node_exporter-${NODE_EXPORTER_VERSION}.linux-${PLATFORM}.tar.gz
tar xzfv node_exporter.tar.gz
sudo mkdir -p /etc/node_exporter
sudo install node_exporter-${NODE_EXPORTER_VERSION}.linux-${PLATFORM}/node_exporter /etc/node_exporter/
