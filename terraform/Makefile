include ../make/Makefile
NAME = terraform
TF_ARCH=$(shell bash $(MAKE_DIR)/../terraform/get_tf_arch)
VERSION = 0.11.11

docker-build:
	@echo "Building image $(IMAGE):$(VERSION)-$(DOCKER_ARCH).."
	docker build --build-arg tf_version=${VERSION} \
                     --build-arg tf_arch=${TF_ARCH} \
	             -t $(IMAGE):$(VERSION)-$(DOCKER_ARCH) .
