include ../make/Makefile
NAME=lnd
VERSION=$(shell cat VERSION)

pre-build:
	@echo "Building lnd in container.."
	docker build --build-arg goarch=$(DOCKER_ARCH) -t lnd-build -f Dockerfile.build .
	@echo "Running lnd container.."
	docker run --name lnd-build lnd-build
	@echo "Copying out build artifacts.."
	docker cp lnd-build:/home/go/bin .
	docker rm lnd-build
