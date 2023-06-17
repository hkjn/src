USERNAME=hkjn
NAME=undefined-name
VERSION=undefined-version
IMAGE=$(USERNAME)/$(NAME)
MAKE_DIR=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
DOCKER_ARCH=$(shell env bash $(MAKE_DIR)/get_docker_arch)
BUILDER ?= podman
HOST_ARCH = $(shell uname -m)
SHELL=/usr/bin/env bash
.DEFAULT_GOAL=all
.PHONY: pre-build main-build post-build build push do-push

all: build push

build: pre-build main-build post-build

pre-build:

main-build: make-manifest build-amd64 build-arm64

make-manifest:
	@echo "Creating multiarch manifest for $(IMAGE)-$(VERSION)"
	$(BUILDER) manifest create $(IMAGE)-$(VERSION)

build-amd64:
	@echo "Building image $(IMAGE):$(VERSION) for amd64.."
	$(BUILDER) build --platform linux/amd64 -t $(IMAGE):$(VERSION)-amd64 .
	# TODO: 'manifest add' tries to add remote image docker://..
	# to the list (1st arg), so seems we have to push to registry for
	# the specific -amd64 and -arm64 images first
	$(BUILDER) manifest add $(USERNAME)-$(NAME)-$(VERSION) $(IMAGE):$(VERSION)-amd64
	# TODO: above command also gives error if manifest already exists;
	# podman manifest inspect $(USERNAME)-$(NAME)-$(VERSION) would check

build-arm64:
	@echo "Building image $(IMAGE):$(VERSION) for arm64.."
	# TODO: --platform should automatically add ARGS
	# [TARGETARCH TARGETOS TARGETPLATFORM]

	# TODO: try with --squash or --squash-all
	$(BUILDER) build --platform linux/arm64 -t $(IMAGE):$(VERSION)-arm64 .
	$(BUILDER) manifest add $(USERNAME)-$(NAME)-$(VERSION) $(IMAGE):$(VERSION)-arm64

tag:
	$(BUILDER) tag $(IMAGE):$(VERSION)-$(HOST_ARCH) $(IMAGE):$(VERSION)

post-build:
	@echo "Pushing manifest $(USERNAME)-$(IMAGE)-$(VERSION).."
	podman manifest push --all $(USERNAME)-$(NAME)-$(VERSION) $(IMAGE):$(VERSION)

push: build do-push

do-push:
	@echo "Pushing manifest $(IMAGE)-$(VERSION).."
	$(BUILDER) manifest push --all $(USERNAME)-$(NAME)-$(VERSION) $(IMAGE):$(VERSION)
	@echo "Pushing image $(IMAGE):$(VERSION).."
	$(BUILDER) push $(IMAGE):$(VERSION)
	@echo "Cleaning up manifest $(IMAGE)-$(VERSION).."
	$(BUILDER) manifest rm $(USERNAME)-$(NAME)-$(VERSION)
