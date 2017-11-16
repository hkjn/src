.PHONY : all blocksci bitcoin golang

all: blocksci bitcoin golang

blocksci:
	@echo "Making blocksci.."
	$(MAKE) -C blocksci

bitcoin:
	@echo "Making bitcoin.."
	$(MAKE) -C bitcoin

golang:
	@echo "Making golang.."
	$(MAKE) -C golang

