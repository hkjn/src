include ../make/Makefile
NAME=bitcoin
VERSION=lightning-internal-error
BITCOIN_VERSION=ffc6e4
LIGHTNING_VERSION=03917c

docker-build:
	@echo "Building image $(IMAGE):$(VERSION)-$(DOCKER_ARCH).."
	docker build --build-arg bitcoin_version=${BITCOIN_VERSION} --build-arg lightning_version=${LIGHTNING_VERSION} -t $(IMAGE):$(VERSION)-$(DOCKER_ARCH) .
