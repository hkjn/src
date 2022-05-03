include ../make/Makefile
NAME = terraform
TF_ARCH=$(shell bash $(MAKE_DIR)/../terraform/get_tf_arch)
VERSION = 1.1.9
DOCKER=podman

main-build:
	@echo "Building image $(IMAGE):$(VERSION)-$(DOCKER_ARCH).."
	$(DOCKER) build --build-arg tf_version=${VERSION} \
                     --build-arg tf_arch=${TF_ARCH} \
	             -t $(IMAGE):$(VERSION)-$(DOCKER_ARCH) .
