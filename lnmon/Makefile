include ../make/Makefile

NAME=lnmon
VERSION=0.1.0
GOOS ?= linux
GOARCH ?= arm
DEPLOYHOST ?= prod.hkjn.me
DEPLOYPORT ?= 6200

.DEFAULT_GOAL=deploy

build-docker: gen-tmpl
	@echo "Building lnmon in container.."
	docker build -t lnmon-build -f Dockerfile.build .
	@echo "Running lnmon-build container.."
	docker run --name lnmonbuild lnmon-build
	@echo "Copying out lnmonbuild artifacts.."
	docker cp lnmonbuild:/home/go/bin/linux_arm/lnmon .
	docker rm lnmonbuild

build-bin: gen-tmpl
	@echo "Building lnmon on host, targeting $(GOOS) / $(GOARCH).."
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
	    -ldflags="-X main.lnmonVersion=$(shell git log -1 --pretty=format:'%h')" \
	    -o lnmon .

gen-tmpl:
	@echo "Generating bindata.go.."
	docker run --rm --user root \
	           -v $(shell pwd):/home/go/src/hkjn.me/src/lnmon \
		   -w /home/go/src/hkjn.me/src/lnmon \
	    hkjn/golang:tip go-bindata *.tmpl

deploy: build-bin
	@echo "Deploying lnmon.next binary to $(DEPLOYHOST):$(DEPLOYPORT).."
	scp -P $(DEPLOYPORT) lnmon $(DEPLOYHOST):/etc/bins/lnmon.next
