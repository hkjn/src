all: blocksci bitcoin

blocksci:
	@echo "Making blocksci.."
	$(MAKE) -C blocksci

bitcoin:
	@echo "Making bitcoin.."
	$(MAKE) -C bitcoin
