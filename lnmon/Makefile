include ../make/Makefile

NAME=lnmon
VERSION=$(shell git log -1 --pretty=format:'%h')
GOOS ?= linux
GOARCH ?= arm
DEPLOYUSER ?= $(USER)
DEPLOYHOST ?= prod.hkjn.me
DEPLOYPORT ?= 6200

.DEFAULT_GOAL=deploy

build-docker: gen-tmpl-docker
	@echo "Building lnmon in container.."
	docker build -t lnmon-build --build-arg goarch=$(GOARCH) -f Dockerfile.build .
	@echo "Running lnmon-build container.."
	docker run --name lnmonbuild lnmon-build
	@echo "Copying out lnmonbuild artifacts.."
	docker cp lnmonbuild:/home/go/bin/lnmon .
	docker rm lnmonbuild

build-bin: gen-tmpl
	@echo "Building lnmon on host, targeting $(GOOS) / $(GOARCH).."
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
	    -ldflags="-X main.lnmonVersion=$(VERSION)" \
	    -o lnmon .

gen-tmpl:
	@echo "Generating bindata.go on host.."
	go-bindata *.tmpl

gen-tmpl-docker:
	@echo "Generating bindata.go in container.."
	docker run --rm --user root \
	           -v $(shell pwd):/home/go/src/hkjn.me/src/lnmon \
		   -w /home/go/src/hkjn.me/src/lnmon \
	    hkjn/golang:tip go-bindata *.tmpl

deploy: build-bin
	@echo "Deploying lnmon.next binary to $(DEPLOYHOST):$(DEPLOYPORT).."
	scp -P $(DEPLOYPORT) lnmon $(DEPLOYUSER)@$(DEPLOYHOST):/etc/bins/lnmon.next
