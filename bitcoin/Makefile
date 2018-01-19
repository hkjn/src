include ../make/Makefile
NAME=bitcoin
VERSION=v0.15.1-lightning
BITCOIN_VERSION=v0.15.1
LIGHTNING_VERSION=ced486

pre-build:
	@echo "Building bitcoin in container.."
	docker build --build-arg bitcoin_version=$(BITCOIN_VERSION) --build-arg lightning-version=${LIGHTNING_VERSION} -t bitcoin-build -f Dockerfile.build .
	@echo "Running build container.."
	docker run --name bbuild bitcoin-build
	@echo "Copying out build artifacts.."
	docker cp bbuild:/build/bin .
	docker cp bbuild:/build/lib .
	docker rm bbuild
