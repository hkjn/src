include ../make/Makefile
NAME=btcd
VERSION=$(shell cat VERSION)

pre-build:
	@echo "Building btcd in container.."
	docker build  -t btcd-build --build-arg goarch=$(DOCKER_ARCH) -f Dockerfile.build .
	@echo "Running btcd container.."
	docker run --name btcdbuild btcd-build
	@echo "Copying out build artifacts.."
	docker cp btcdbuild:/home/go/bin .
	docker rm btcdbuild

post-build:
	@echo "Not doing anything post-build, no need to squash single layer.."
