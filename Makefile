.PHONY : all blocksci bitcoin gcloud golang gpg proto s3cmd workspace

all: blocksci bitcoin gcloud golang gpg lnd proto s3cmd workspace

blocksci:
	@echo "Making blocksci.."
	$(MAKE) -C blocksci

bitcoin:
	@echo "Making bitcoin.."
	$(MAKE) -C bitcoin

gcloud:
	@echo "Making gcloud.."
	$(MAKE) -C gcloud

golang:
	@echo "Making golang.."
	$(MAKE) -C golang

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

