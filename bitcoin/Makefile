include ../make/Makefile
NAME = bitcoin
VERSION = v0.16.0-2018-03-08
BITCOIN_VERSION ?= v0.16.0
LIGHTNING_VERSION ?= 3e9e1
LIGHTNING_REMOTE ?= https://github.com/ElementsProject/lightning

docker-build:
	@echo "Building image $(IMAGE):$(VERSION)-$(DOCKER_ARCH).."
	docker build --build-arg bitcoin_version=${BITCOIN_VERSION} \
	             --build-arg lightning_version=${LIGHTNING_VERSION} \
	             --build-arg lightning_remote=${LIGHTNING_REMOTE} \
	             -t $(IMAGE):$(VERSION)-$(DOCKER_ARCH) .
