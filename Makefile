.PHONY : all blocksci bitcoin golang s3cmd

all: blocksci bitcoin golang s3cmd

blocksci:
	@echo "Making blocksci.."
	$(MAKE) -C blocksci

bitcoin:
	@echo "Making bitcoin.."
	$(MAKE) -C bitcoin

golang:
	@echo "Making golang.."
	$(MAKE) -C golang

s3cmd:
	@echo "Making s3cmd.."
	$(MAKE) -C s3cmd

workspace:
	@echo "Making workspace.."
	$(MAKE) -C workspace

