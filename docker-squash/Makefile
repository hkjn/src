USERNAME=hkjn
NAME=docker-squash
SHELL=/bin/bash
IMAGE=$(USERNAME)/$(NAME)
DOCKER_ARCH=$(shell bash get_docker_arch)
VERSION=1.1.0

.PHONY: pre-build docker-build post-build build push do-push post-push

build: pre-build docker-build post-build

pre-build:

post-build:

post-push:
	@echo "Pushing multi-arch image manifests.."
	docker run -v $(HOME)/.docker:/root/.docker:ro --rm hkjn/manifest-tool \
	       push from-args --platforms linux/amd64,linux/arm \
	                      --template $(IMAGE):ARCH \
	                      --target $(IMAGE)

docker-build:
	@echo "Building image.."
	docker build -t $(IMAGE):$(VERSION)-$(DOCKER_ARCH) .
	@echo "Tagging image.."
	docker tag $(IMAGE):$(VERSION)-$(DOCKER_ARCH) $(IMAGE):$(DOCKER_ARCH)

push: do-push post-push

do-push:
	@echo "Pushing image.."
	docker push $(IMAGE):$(DOCKER_ARCH)
