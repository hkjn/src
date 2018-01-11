include ../make/Makefile
NAME=bitcointools
VERSION=$(shell cat VERSION)

docker-build:
	@echo "Building image $(IMAGE):$(VERSION)-$(DOCKER_ARCH).."
	docker build --build-arg tools_version=${VERSION} -t $(IMAGE):$(VERSION)-$(DOCKER_ARCH) .
