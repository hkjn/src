include ../make/Makefile

NAME=bcmon
VERSION=0.1.0
GOOS ?= linux
GOARCH ?= arm
DEPLOYHOST ?= prod.hkjn.me
DEPLOYPORT ?= 6200

.DEFAULT_GOAL=deploy

build-docker: gen-tmpl
	@echo "Building bcmon in container.."
	docker build -t bcmon-build -f Dockerfile.build .
	@echo "Running bcmon-build container.."
	docker run --name bcmonbuild bcmon-build
	@echo "Copying out bcmonbuild artifacts.."
	docker cp bcmonbuild:/home/go/bin/linux_arm/bcmon .
	docker rm bcmonbuild

build-bin: gen-tmpl
	@echo "Building bcmon on host.."
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o bcmon .

gen-tmpl:
	@echo "Generating bindata.go.."
	docker run --rm --user root \
	           -v $(shell pwd):/home/go/src/hkjn.me/src/lnmon \
		   -w /home/go/src/hkjn.me/src/lnmon \
	    hkjn/golang:tip go-bindata *.tmpl

install-local: build-docker
	sudo install bcmon /etc/bins/

deploy: build-bin
	@echo "Deploying bcmon.next binary.."
	scp -P $(DEPLOYPORT) bcmon $(DEPLOYHOST):/etc/bins/bcmon.next
