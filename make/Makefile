USERNAME=hkjn
NAME=undefined-name
VERSION=undefined-version
IMAGE=$(USERNAME)/$(NAME)
MAKE_DIR=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
DOCKER_ARCH=$(shell bash $(MAKE_DIR)/../scripts/get_docker_arch)
SHELL=/bin/bash

.DEFAULT_GOAL=push
.PHONY: pre-build docker-build post-build build	push do-push post-push

build: pre-build docker-build post-build

pre-build:

post-build:
	@echo "Squashing image $(IMAGE).."
	docker run --rm \
		   -v /var/run/docker.sock:/var/run/docker.sock \
		   hkjn/docker-squash \
		     -t $(IMAGE):$(VERSION)-$(DOCKER_ARCH) \
		        $(IMAGE):$(VERSION)-$(DOCKER_ARCH)

post-push:

docker-build:
	@echo "Building image.."
	docker build -t $(IMAGE):$(VERSION)-$(DOCKER_ARCH) .

push: build do-push post-push

do-push:
	@echo "Pushing image.."
	docker push $(IMAGE):$(VERSION)-$(DOCKER_ARCH)
