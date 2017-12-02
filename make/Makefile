USERNAME=hkjn
NAME=undefined-name
VERSION=undefined-version
IMAGE=$(USERNAME)/$(NAME)
MAKE_DIR=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
DOCKER_ARCH=$(shell bash $(MAKE_DIR)/../scripts/get_docker_arch)
SHELL=/bin/bash

.DEFAULT_GOAL=all
.PHONY: pre-build docker-build post-build build	push do-push post-push

all: build push

build: pre-build docker-build post-build

pre-build:

docker-build:
	@echo "Building image $(IMAGE):$(VERSION)-$(DOCKER_ARCH).."
	docker build -t $(IMAGE):$(VERSION)-$(DOCKER_ARCH) .

post-build:
	@echo "Squashing image $(IMAGE):$(VERSION)-$(DOCKER_ARCH).."
	docker run --rm \
		   -v /var/run/docker.sock:/var/run/docker.sock \
		   hkjn/docker-squash \
		     -t $(IMAGE):$(VERSION)-$(DOCKER_ARCH) \
		        $(IMAGE):$(VERSION)-$(DOCKER_ARCH)

push: build do-push post-push

post-push:
	@echo "Pushing multiarch manifest for $(IMAGE):$(VERSION)-$(DOCKER_ARCH).."
	docker run -v ${HOME}/.docker:/root/.docker:ro --rm \
	       hkjn/manifest-tool \
	       push from-args --platforms linux/amd64,linux/arm \
	                      --template $(IMAGE):$(VERSION)-ARCH \
	                      --target $(IMAGE):$(VERSION)

do-push:
	@echo "Pushing image.."
	docker push $(IMAGE):$(VERSION)-$(DOCKER_ARCH)
