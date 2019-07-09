all: deploy-prometheus-tor deploy-lightningd-plugins

refresh-prometheus-bitcoind:
	bash fetch_prometheus_bitcoind.sh

refresh-prometheus-lightningd:
	bash fetch_prometheus_lightningd.sh

deploy-prometheus-tor:
	sudo cp prometheus_tor.py /opt/monitoring/
	sudo chown monitoring:monitoring /opt/monitoring/prometheus_tor.py
	sudo systemctl restart prometheus-tor

deploy-blink-network:
	sudo cp blink_network.py /opt/bitcoin/
	sudo chown btc:btc /opt/bitcoin/blink_network.py
	sudo systemctl restart blink-network

deploy-lightningd-plugins:
	sudo cp prometheus-lightningd.py /opt/lightning/
	sudo chown ln:btc /opt/bitcoin/blink_network.py
	sudo systemctl restart lightningd
