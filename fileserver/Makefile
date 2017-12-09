include ../make/Makefile
NAME=fileserver
VERSION=$(shell cat VERSION)

pre-build:
	@echo "Building fileserver in container.."
	docker build -t fileserver-build -f Dockerfile.build .
	@echo "Running fileserver-build container.."
	docker run --name fsbuild fileserver-build
	@echo "Copying out fsbuild artifacts.."
	docker cp fsbuild:/build/fileserver .
	docker rm fsbuild
