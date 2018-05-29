

all:
	docker run --rm -it -p 3999:3999 --name present \
	           -v $(shell pwd):/var/www -w /var/www hkjn/golang:tip \
	           sh -c "go get -v golang.org/x/tools/cmd/present && present -http=0.0.0.0:3999 -play=false -nacl=false"
