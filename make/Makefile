USERNAME=hkjn
NAME=undefined-name
VERSION=undefined-version
IMAGE=$(USERNAME)/$(NAME)
MAKE_DIR=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
DOCKER_ARCH=$(shell env bash $(MAKE_DIR)/get_docker_arch)
DOCKER=podman
SHELL=/usr/bin/env bash
SQUASH_ENABLED=0
PUSH_ENABLED=1
.DEFAULT_GOAL=all
.PHONY: pre-build main-build post-build build push do-push post-push

all: build push

build: pre-build main-build post-build

pre-build:

main-build:
	@echo "Building image $(IMAGE):$(VERSION)-$(DOCKER_ARCH).."
	$(DOCKER) build -t $(IMAGE):$(VERSION)-$(DOCKER_ARCH) .

post-build:
	@if [ "$(SQUASH_ENABLED)" = "1" ]; then \
		$(DOCKER) run --rm \
		           -v /var/run/docker.sock:/var/run/docker.sock \
			   hkjn/docker-squash:1.1.0 -t $(IMAGE):$(VERSION)-$(DOCKER_ARCH) $(IMAGE):$(VERSION)-$(DOCKER_ARCH); \
	else \
		echo "Not squashing image, since SQUASH_ENABLED != 1."; \
	fi

push: build do-push post-push

post-push:
	# xx: how to handle bind mounts on macos?
	@if [ "$(PUSH_ENABLED)" = "1" ]; then \
		$(DOCKER) run -v ${HOME}/.docker:/root/.docker:ro \
		           --rm hkjn/manifest-tool push from-args \
			     --platforms linux/amd64,linux/arm64 \
		             --template $(IMAGE):$(VERSION)-ARCH \
		             --target $(IMAGE):$(VERSION); \
	else \
		echo "Not pushing multiarch manifest, since PUSH_ENABLED=$(PUSH_ENABLED)."; \
	fi

do-push:
	@if [ "$(PUSH_ENABLED)" = "1" ]; then \
		$(DOCKER) push $(IMAGE):$(VERSION)-$(DOCKER_ARCH); \
	else \
		echo "Not pushing image, since PUSH_ENABLED=$(PUSH_ENABLED)"; \
	fi
