include ../make/Makefile

NAME=bcmon
VERSION=0.1.0

.DEFAULT_GOAL=deploy

build-docker: gen-tmpl bcmon
	@echo "Building bcmon in container.."
	docker build -t bcmon-build -f Dockerfile.build .
	@echo "Running bcmon-build container.."
	docker run --name bcmonbuild lnmon-build
	@echo "Copying out bcmonbuild artifacts.."
	docker cp bcmonbuild:/home/go/bin/linux_arm/bcmon .
	docker rm bcmonbuild

build-bin: gen-tmpl
	@echo "Building bcmon on host.."
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o bcmon .

gen-tmpl:
	@echo "Generating bindata.go.."
	go-bindata *.tmpl

deploy: build-bin
	@echo "Deploying bcmon.next binary.."
	scp -P 6200 bcmon prod.hkjn.me:/etc/bins/bcmon.next
