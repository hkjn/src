include ../make/Makefile
NAME=bitcoin
VERSION=lightning-2018-02-07
BITCOIN_VERSION=1462bd
LIGHTNING_VERSION=b3d4e1

docker-build:
	@echo "Building image $(IMAGE):$(VERSION)-$(DOCKER_ARCH).."
	docker build --build-arg bitcoin_version=${BITCOIN_VERSION} --build-arg lightning_version=${LIGHTNING_VERSION} -t $(IMAGE):$(VERSION)-$(DOCKER_ARCH) .
