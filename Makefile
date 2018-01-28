.PHONY : all blocksci bitcoin btcd clightning config fileserver gcloud golang googleauth gpg lnd proto s3cmd workspace

all: blocksci bitcoin btcd clightning config fileserver gcloud golang googleauth gpg lnd proto s3cmd workspace

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

proto:
	@echo "Making proto.."
	$(MAKE) -C proto

s3cmd:
	@echo "Making s3cmd.."
	$(MAKE) -C s3cmd

workspace:
	@echo "Making workspace.."
	$(MAKE) -C workspace

