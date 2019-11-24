NAME=dashboard
VERSION=$(shell cat VERSION)

build:
	bash generate_version.sh
	bash generate_bindata.sh
	go build ./cmd/gomon
