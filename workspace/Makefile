VERSION=$(shell cat ${HOME}/src/hkjn.me/src/workspace/VERSION)

all: alpine aruna arch zero

alpine:
	@echo "Making alpine.."
	$(MAKE) -C alpine

arch:
	@echo "Making arch.."
	$(MAKE) -C arch

aruna:
	@echo "Making aruna.."
	$(MAKE) -C aruna

zero:
	@echo "Making zero.."
	$(MAKE) -C zero

.PHONY : all alpine aruna arch zero
