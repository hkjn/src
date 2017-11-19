USERNAME=hkjn
NAME=$(shell basename $(PWD))
IMAGE=$(USERNAME)/$(NAME)
DOCKER_ARCH=$(shell bash $(PWD)/get_docker_arch)
VERSION=$(shell cat ../VERSION)
SHELL=/bin/bash

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

push: do-push post-push

do-push:
	@echo "Pushing image.."
	docker push $(IMAGE):$(VERSION)
