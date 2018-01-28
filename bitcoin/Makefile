include ../make/Makefile
NAME=bitcoin
VERSION=f4c942e3-lightning
BITCOIN_VERSION=126000
LIGHTNING_VERSION=feb7340

pre-build:
	@echo "Building bitcoin in container.."
	docker build --build-arg bitcoin_version=$(BITCOIN_VERSION) --build-arg lightning_version=${LIGHTNING_VERSION} -t bitcoin-build -f Dockerfile.build .
	@echo "Running build container.."
	docker run --name bbuild bitcoin-build
	@echo "Copying out build artifacts.."
	docker cp bbuild:/build/bin .
	docker cp bbuild:/build/lib .
	docker cp bbuild:/build/libexec .
	docker rm bbuild
