include ../make/Makefile

NAME=lnmon
VERSION=0.1.0

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
	@echo "Building lnmon on host.."
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o lnmon lnmon.go bindata.go

gen-tmpl:
	@echo "Generating bindata.go.."
	go-bindata *.tmpl

deploy: build-bin
	@echo "Deploying lnmon.next binary.."
	scp -P 6200 lnmon prod.hkjn.me:/crypt/lightning/lnmon.next
	@echo "Setting cap_net_bind_service on remote binary.."
	ssh -p 6200 prod.hkjn.me 'sudo setcap "cap_net_bind_service=+ep" /crypt/lightning/lnmon.next'
