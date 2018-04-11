include ../make/Makefile

.DEFAULT_GOAL=install

build-docker:
	@echo "Building caddy in container.."
	docker build -t caddy-build -f Dockerfile .
	@echo "Running caddy-build container.."
	docker run --name caddy-build caddy-build
	@echo "Copying out caddy artifacts.."
	docker cp caddy-build:/home/go/bin/caddy .
	docker rm caddy-build

install: build-docker
	@echo "Setting capabilities to allow binding to privileged ports.."
	sudo setcap "cap_net_bind_service=+ep" caddy
	@echo "Adjusting ulimit.."
	ulimit -n 8192 caddy
	@echo "Installing caddy binary.."
	install caddy /etc/caddy
	# TODO: also install .service here?
