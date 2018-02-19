include ../make/Makefile
NAME=bitcoin
VERSION=lightning-2018-02-19
BITCOIN_VERSION=dcfe21
LIGHTNING_VERSION=4fccbb

docker-build:
	@echo "Building image $(IMAGE):$(VERSION)-$(DOCKER_ARCH).."
	docker build --build-arg bitcoin_version=${BITCOIN_VERSION} --build-arg lightning_version=${LIGHTNING_VERSION} -t $(IMAGE):$(VERSION)-$(DOCKER_ARCH) .
