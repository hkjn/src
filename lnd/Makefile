include ../make/Makefile
NAME=lnd
VERSION=9479b0

pre-build:
	@echo "Building lnd in container.."
	docker build --build-arg goarch=$(DOCKER_ARCH) --build-arg lnd_version=${VERSION} -t lnd-build -f Dockerfile.build .
	@echo "Running lnd container.."
	docker run --name lnd-build lnd-build
	@echo "Copying out build artifacts.."
	docker cp lnd-build:/home/go/bin .
	docker rm lnd-build

post-build:
	@echo "Not doing anything post-build, no need to squash single layer.."
