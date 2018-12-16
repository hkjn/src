.PHONY : all alpine blocksci bitcoin btcd clightning config fileserver gcloud golang googleauth gpg lnd openvpn prober probes proto s3cmd terraform workspace

all: alpine blocksci bitcoin btcd clightning config fileserver gcloud golang googleauth gpg lnd openvpn prober probes proto s3cmd terraform workspace

alpine:
	@echo "Making alpine.."
	$(MAKE) -C alpine

blocksci:
	@echo "Making blocksci.."
	$(MAKE) -C blocksci

bitcoin:
	@echo "Making bitcoin.."
	$(MAKE) -C bitcoin

btcd:
	@echo "Making btcd.."
	$(MAKE) -C btcd

clightning:
	@echo "Making clightning.."
	$(MAKE) -C clightning

config:
	@echo "Making config.."
	$(MAKE) -C config

fileserver:
	@echo "Making fileserver.."
	$(MAKE) -C fileserver

gcloud:
	@echo "Making gcloud.."
	$(MAKE) -C gcloud

golang:
	@echo "Making golang.."
	$(MAKE) -C golang

googleauth:
	@echo "Making googleauth.."
	$(MAKE) -C googleauth

gpg:
	@echo "Making gpg.."
	$(MAKE) -C gpg

lnd:
	@echo "Making lnd.."
	$(MAKE) -C lnd

openvpn:
	@echo "Making openvpn.."
	$(MAKE) -C openvpn

prober:
	@echo "Making prober.."
	$(MAKE) -C prober

probes:
	@echo "Making probes.."
	$(MAKE) -C probes

proto:
	@echo "Making proto.."
	$(MAKE) -C proto

s3cmd:
	@echo "Making s3cmd.."
	$(MAKE) -C s3cmd

terraform:
	@echo "Making terraform.."
	$(MAKE) -C terraform

workspace:
	@echo "Making workspace.."
	$(MAKE) -C workspace

