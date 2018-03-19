USERNAME=hkjn
NAME=undefined-name
VERSION=undefined-version
IMAGE=$(USERNAME)/$(NAME)
MAKE_DIR=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
DOCKER_ARCH=$(shell bash $(MAKE_DIR)/get_docker_arch)
SHELL=/bin/bash
PUSH_ENABLED=1
.DEFAULT_GOAL=all
.PHONY: pre-build docker-build post-build build push do-push post-push

all: build push

build: pre-build docker-build post-build

pre-build:

docker-build:
	@echo "Building image $(IMAGE):$(VERSION)-$(DOCKER_ARCH).."
	docker build -t $(IMAGE):$(VERSION)-$(DOCKER_ARCH) .

post-build:
	@if [ "$(PUSH_ENABLED)" = "1" ]; then \
		docker run --rm \
		           -v /var/run/docker.sock:/var/run/docker.sock \
			   hkjn/docker-squash -t $(IMAGE):$(VERSION)-$(DOCKER_ARCH) $(IMAGE):$(VERSION)-$(DOCKER_ARCH); \
	else \
		echo "Not squashing image, since PUSH_ENABLED=$(PUSH_ENABLED)."; \
	fi

push: build do-push post-push

post-push:
	@if [ "$(PUSH_ENABLED)" = "1" ]; then \
		docker run -v ${HOME}/.docker:/root/.docker:ro \
		           --rm hkjn/manifest-tool push from-args \
		             --platforms linux/amd64,linux/arm \
		             --template $(IMAGE):$(VERSION)-ARCH \
		             --target $(IMAGE):$(VERSION); \
	else \
		echo "Not pushing multiarch manifest, since PUSH_ENABLED=$(PUSH_ENABLED)."; \
	fi

do-push:
	@if [ "$(PUSH_ENABLED)" = "1" ]; then \
		docker push $(IMAGE):$(VERSION)-$(DOCKER_ARCH); \
	else \
		echo "Not pushing image, since PUSH_ENABLED=$(PUSH_ENABLED)"; \
	fi
