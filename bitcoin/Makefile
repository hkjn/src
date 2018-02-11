include ../make/Makefile
NAME=bitcoin
VERSION=lightning-2018-02-09
BITCOIN_VERSION=89005
LIGHTNING_VERSION=d5effc

docker-build:
	@echo "Building image $(IMAGE):$(VERSION)-$(DOCKER_ARCH).."
	docker build --build-arg bitcoin_version=${BITCOIN_VERSION} --build-arg lightning_version=${LIGHTNING_VERSION} -t $(IMAGE):$(VERSION)-$(DOCKER_ARCH) .
