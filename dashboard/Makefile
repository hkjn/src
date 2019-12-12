NAME=dashboard
VERSION=$(shell cat VERSION)

build: fetch-deps format
	bash generate_version.sh
	bash generate_bindata.sh
	go build ./cmd/gomon

fetch-deps:
	go get -u github.com/go-bindata/go-bindata/...
	go get -v ./...

format:
	go fmt .

install:
	install gomon /usr/local/bin/
