present:
	go get -v golang.org/x/tools/cmd/present
	present -http=127.0.0.1:8091 -play=false -nacl=false

present-docker:
	docker run --rm -it -p 8091:8091 --name present \
	           -v $(shell pwd):/var/www -w /var/www hkjn/golang:tip \
	           sh -c "go get -v golang.org/x/tools/cmd/present && present -http=0.0.0.0:8091 -play=false -nacl=false"
