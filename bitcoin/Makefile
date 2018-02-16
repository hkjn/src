include ../make/Makefile
NAME=bitcoin
VERSION=lightning-2018-02-16
BITCOIN_VERSION=26dc2d
LIGHTNING_VERSION=e48a97

docker-build:
	@echo "Building image $(IMAGE):$(VERSION)-$(DOCKER_ARCH).."
	docker build --build-arg bitcoin_version=${BITCOIN_VERSION} --build-arg lightning_version=${LIGHTNING_VERSION} -t $(IMAGE):$(VERSION)-$(DOCKER_ARCH) .
