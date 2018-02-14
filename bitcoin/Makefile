include ../make/Makefile
NAME=bitcoin
VERSION=lightning-2018-02-14
BITCOIN_VERSION=252ae7
LIGHTNING_VERSION=e97ee3

docker-build:
	@echo "Building image $(IMAGE):$(VERSION)-$(DOCKER_ARCH).."
	docker build --build-arg bitcoin_version=${BITCOIN_VERSION} --build-arg lightning_version=${LIGHTNING_VERSION} -t $(IMAGE):$(VERSION)-$(DOCKER_ARCH) .
