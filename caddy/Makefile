include ../make/Makefile

build-docker:
	@echo "Building caddy in container.."
	docker build -t caddy-build -f Dockerfile .
	@echo "Running caddy-build container.."
	docker run --name caddy-build caddy-build
	@echo "Copying out caddy artifacts.."
	docker cp caddy-build:/home/go/bin/caddy .
	docker rm caddy-build
