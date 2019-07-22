.PHONY: configure-host deploy
.DEFAULT_GOAL=all

all: configure-host deploy

configure-host:
	scp configure-host root@vpn.hkjn.me:
	ssh root@vpn.hkjn.me bash configure-host

deploy:
	scp create-iptables-restore-service iptables.rules vpn:wireguard/
	ssh vpn sudo bash wireguard/create-iptables-restore-service
