include ../make/Makefile
NAME = bitcoin
VERSION = ln-2018-06-07
BITCOIN_VERSION ?= f8bcef3
BITCOIN_ARCH=$(shell bash $(MAKE_DIR)/../bitcoin/get_bitcoin_arch)
LIGHTNING_VERSION ?= 294c76
LIGHTNING_REMOTE ?= https://github.com/ElementsProject/lightning

docker-build:
	@echo "Building image $(IMAGE):$(VERSION)-$(DOCKER_ARCH).."
	docker build --build-arg bitcoin_version=${BITCOIN_VERSION} \
                     --build-arg bitcoin_arch=${BITCOIN_ARCH} \
	             --build-arg lightning_version=${LIGHTNING_VERSION} \
	             --build-arg lightning_remote=${LIGHTNING_REMOTE} \
	             -t $(IMAGE):$(VERSION)-$(DOCKER_ARCH) .
