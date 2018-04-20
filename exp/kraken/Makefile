include ../../make/Makefile
NAME=kraken
VERSION=$(shell cat VERSION)

run:
	docker run --name kraken --rm -it \
	           -v $(HOME)/src/hkjn.me/src/exp/kraken:/home/go/src/hkjn.me/src/exp/kraken \
	           -w /home/go/src/hkjn.me/src/exp/kraken \
	           --env-file=/etc/secrets/kraken/kraken.env \
	    hkjn/golang:tip go run kraken.go
