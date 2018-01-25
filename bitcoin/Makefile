include ../make/Makefile
NAME=bitcoin
VERSION=f4c942e3-lightning
BITCOIN_VERSION=126000
LIGHTNING_VERSION=bab3b1a1

# FIXMEH: first hit this assert on startup:
# https://github.com/ElementsProject/lightning/blob/master/lightningd/peer_control.c#L2750
# then upgraded to newer commit, got 'failing to start with lightningd(79): SQLITE version mismatch: compiled 3021000, now 3019002'
# first version is what's statically coded (in libsqlite?), second is what the .sqlite3 says it's in
# seems to be due to different libsqlite versions in ubuntu and arch linux..

pre-build:
	@echo "Building bitcoin in container.."
	docker build --build-arg bitcoin_version=$(BITCOIN_VERSION) --build-arg lightning_version=${LIGHTNING_VERSION} -t bitcoin-build -f Dockerfile.build .
	@echo "Running build container.."
	docker run --name bbuild bitcoin-build
	@echo "Copying out build artifacts.."
	docker cp bbuild:/build/bin .
	docker cp bbuild:/build/lib .
	docker rm bbuild

docker-build:
	@echo "Building image $(IMAGE):$(VERSION)-$(DOCKER_ARCH).."
	docker build --build-arg bitcoin_version=${BITCOIN_VERSION} --build-arg lightning_version=${LIGHTNING_VERSION} -t $(IMAGE):$(VERSION)-$(DOCKER_ARCH) .
