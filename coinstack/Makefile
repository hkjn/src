all: refresh-prometheus-lightningd refresh-prometheus-bitcoind

refresh-prometheus-bitcoind:
	bash fetch_prometheus_bitcoind.sh

refresh-prometheus-lightningd:
	bash fetch_prometheus_lightningd.sh

deploy-prometheus-tor:
	sudo cp prometheus_tor.py /opt/monitoring/
	sudo chown monitoring:monitoring /opt/monitoring/prometheus_tor.py
	sudo systemctl daemon-reload
	sudo systemctl restart prometheus-tor
